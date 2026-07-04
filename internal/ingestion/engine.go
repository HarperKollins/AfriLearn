// Package ingestion reads curriculum JSON files from disk and upserts them into PostgreSQL.
// This replaces the old internal/scraper package which embedded data directly in Go code.
package ingestion

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/google/uuid"
)

// CurriculumFile mirrors the JSON schema in data/curricula/**/*.json
type CurriculumFile struct {
	Board     string      `json:"board"`
	Subject   string      `json:"subject"`
	Level     string      `json:"level"`
	SourceURL string      `json:"source_url"`
	Topics    []TopicData `json:"topics"`
}

// TopicData holds a top-level topic and all its children.
type TopicData struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Difficulty  string         `json:"difficulty"`
	Subtopics   []SubtopicData `json:"subtopics"`
}

// SubtopicData holds a subtopic, optional course attributes, and its learning objectives.
type SubtopicData struct {
	Name        string   `json:"name"`
	CourseCode  string   `json:"course_code,omitempty"`
	CreditUnits string   `json:"credit_units,omitempty"`
	Semester    string   `json:"semester,omitempty"`
	Objectives  []string `json:"objectives"`
}

// Engine walks data/curricula/ and ingests every JSON file into Postgres.
type Engine struct {
	dataDir string
}

// NewEngine creates a new Engine. dataDir should point to data/curricula/.
func NewEngine(dataDir string) *Engine {
	return &Engine{dataDir: dataDir}
}

// Run walks all JSON files in dataDir and ingests them.
func (e *Engine) Run() error {
	var files []string
	if err := filepath.WalkDir(e.dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".json") {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to walk data directory: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no JSON curriculum files found in %s", e.dataDir)
	}

	log.Printf("📂 Found %d curriculum files to ingest", len(files))

	var successCount, failCount int
	for i, f := range files {
		log.Printf("── [%d/%d] %s", i+1, len(files), filepath.Base(f))
		if err := e.ingestFile(f); err != nil {
			log.Printf("  ❌ Error: %v", err)
			failCount++
		} else {
			successCount++
		}
	}

	log.Printf("\n✅ Ingestion complete: %d succeeded, %d failed", successCount, failCount)
	return nil
}

// IngestFile reads a single JSON file and upserts it into the database.
func (e *Engine) IngestFile(path string) error {
	return e.ingestFile(path)
}

// ValidationError contains details of a curriculum JSON validation failure
type ValidationError struct {
	FilePath string `json:"file_path"`
	Field    string `json:"field"`
	Message  string `json:"message"`
}

// ValidateCurriculumFile performs schema, difficulty, and quality validation on a CurriculumFile
func ValidateCurriculumFile(cf CurriculumFile, filePath string) []ValidationError {
	var errs []ValidationError

	if strings.TrimSpace(cf.Board) == "" {
		errs = append(errs, ValidationError{FilePath: filePath, Field: "board", Message: "board slug is required"})
	}
	if strings.TrimSpace(cf.Subject) == "" {
		errs = append(errs, ValidationError{FilePath: filePath, Field: "subject", Message: "subject slug is required"})
	}
	if strings.TrimSpace(cf.Level) == "" {
		errs = append(errs, ValidationError{FilePath: filePath, Field: "level", Message: "curriculum level is required"})
	}
	if strings.TrimSpace(cf.SourceURL) == "" || (!strings.HasPrefix(cf.SourceURL, "http://") && !strings.HasPrefix(cf.SourceURL, "https://")) {
		errs = append(errs, ValidationError{FilePath: filePath, Field: "source_url", Message: "valid http/https source_url is required"})
	}

	if len(cf.Topics) == 0 {
		errs = append(errs, ValidationError{FilePath: filePath, Field: "topics", Message: "curriculum must contain at least 1 topic"})
	}

	for i, t := range cf.Topics {
		topicPath := fmt.Sprintf("topics[%d]", i)
		if strings.TrimSpace(t.Name) == "" {
			errs = append(errs, ValidationError{FilePath: filePath, Field: topicPath + ".name", Message: "topic name cannot be empty"})
		}
		if t.Difficulty != "" {
			diff := strings.ToLower(t.Difficulty)
			if diff != "easy" && diff != "medium" && diff != "hard" {
				errs = append(errs, ValidationError{FilePath: filePath, Field: topicPath + ".difficulty", Message: fmt.Sprintf("invalid difficulty '%s' (must be easy, medium, or hard)", t.Difficulty)})
			}
		}
		if len(t.Subtopics) == 0 {
			errs = append(errs, ValidationError{FilePath: filePath, Field: topicPath + ".subtopics", Message: fmt.Sprintf("topic '%s' has 0 subtopics", t.Name)})
		}
		for j, st := range t.Subtopics {
			subPath := fmt.Sprintf("%s.subtopics[%d]", topicPath, j)
			if strings.TrimSpace(st.Name) == "" {
				errs = append(errs, ValidationError{FilePath: filePath, Field: subPath + ".name", Message: "subtopic name cannot be empty"})
			}
			if len(st.Objectives) == 0 {
				errs = append(errs, ValidationError{FilePath: filePath, Field: subPath + ".objectives", Message: fmt.Sprintf("subtopic '%s' has 0 learning objectives", st.Name)})
			}
			for k, obj := range st.Objectives {
				if len(strings.TrimSpace(obj)) < 5 {
					errs = append(errs, ValidationError{FilePath: filePath, Field: fmt.Sprintf("%s.objectives[%d]", subPath, k), Message: "objective description too short (<5 chars)"})
				}
			}
		}
	}

	return errs
}

// ValidateAll scans and validates all JSON files without making database changes
func (e *Engine) ValidateAll() (int, []ValidationError, error) {
	var files []string
	if err := filepath.WalkDir(e.dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".json") {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return 0, nil, fmt.Errorf("failed to walk data directory: %w", err)
	}

	var allErrors []ValidationError
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			allErrors = append(allErrors, ValidationError{FilePath: f, Field: "file", Message: "cannot read file: " + err.Error()})
			continue
		}
		var cf CurriculumFile
		if err := json.Unmarshal(data, &cf); err != nil {
			allErrors = append(allErrors, ValidationError{FilePath: f, Field: "json", Message: "invalid JSON syntax: " + err.Error()})
			continue
		}
		errs := ValidateCurriculumFile(cf, filepath.Base(f))
		allErrors = append(allErrors, errs...)
	}

	return len(files), allErrors, nil
}

func (e *Engine) ingestFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", path, err)
	}

	var cf CurriculumFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return fmt.Errorf("invalid JSON in %s: %w", path, err)
	}

	// Run pre-ingestion validation
	valErrors := ValidateCurriculumFile(cf, filepath.Base(path))
	if len(valErrors) > 0 {
		var errMsg []string
		for _, ve := range valErrors {
			errMsg = append(errMsg, fmt.Sprintf("[%s]: %s", ve.Field, ve.Message))
		}
		return fmt.Errorf("validation failed for %s:\n  - %s", filepath.Base(path), strings.Join(errMsg, "\n  - "))
	}

	log.Printf("  🚀 Ingesting [%s / %s] (%d topics)...", cf.Board, cf.Subject, len(cf.Topics))

	curriculumID, err := e.ensureCurriculum(cf)
	if err != nil {
		return fmt.Errorf("cannot ensure curriculum [%s/%s]: %w", cf.Board, cf.Subject, err)
	}

	for i, topic := range cf.Topics {
		topicID, err := e.upsertTopic(curriculumID, topic, i+1)
		if err != nil {
			log.Printf("    ⚠️  Failed to upsert topic '%s': %v", topic.Name, err)
			continue
		}
		for j, subtopic := range topic.Subtopics {
			subtopicID, err := e.upsertSubtopic(topicID, subtopic, j+1)
			if err != nil {
				log.Printf("      ⚠️  Failed to upsert subtopic '%s': %v", subtopic.Name, err)
				continue
			}
			// Clear existing objectives for this subtopic to ensure clean, idempotent upsert
			_, _ = database.DB.Exec(`DELETE FROM learning_objectives WHERE subtopic_id = $1`, subtopicID)

			for k, objective := range subtopic.Objectives {
				if err := e.upsertObjective(subtopicID, objective, k+1); err != nil {
					log.Printf("        ⚠️  Objective [%d] error: %v", k+1, err)
				}
			}
		}
	}
	log.Printf("  ✅ Done [%s / %s]", cf.Board, cf.Subject)
	return nil
}

// ensureCurriculum resolves board + subject + curriculum rows, creating if missing.
func (e *Engine) ensureCurriculum(cf CurriculumFile) (string, error) {
	var boardID string
	err := database.DB.QueryRow(`SELECT id FROM exam_boards WHERE slug = $1`, cf.Board).Scan(&boardID)
	if err != nil {
		return "", fmt.Errorf("exam board '%s' not found: %w", cf.Board, err)
	}

	var subjectID string
	err = database.DB.QueryRow(`SELECT id FROM subjects WHERE slug = $1`, cf.Subject).Scan(&subjectID)
	if err != nil {
		return "", fmt.Errorf("subject '%s' not found: %w", cf.Subject, err)
	}

	currentYear := time.Now().Year()

	var currID string
	err = database.DB.QueryRow(`
		INSERT INTO curricula (id, exam_board_id, subject_id, year, level, source_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (exam_board_id, subject_id, year)
		DO UPDATE SET source_url = EXCLUDED.source_url, updated_at = NOW()
		RETURNING id
	`, uuid.New().String(), boardID, subjectID, currentYear, cf.Level, cf.SourceURL).Scan(&currID)

	if err != nil {
		return "", fmt.Errorf("failed to upsert curriculum: %w", err)
	}
	return currID, nil
}

func (e *Engine) upsertTopic(curriculumID string, topic TopicData, order int) (string, error) {
	slug := Slugify(topic.Name)
	diff := topic.Difficulty
	if diff == "" {
		diff = ClassifyDifficulty(topic.Name)
	}

	var topicID string
	err := database.DB.QueryRow(`
		INSERT INTO topics (id, curriculum_id, slug, name, description, order_index, difficulty, embedding)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULL)
		ON CONFLICT (curriculum_id, slug)
		DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, order_index = EXCLUDED.order_index, difficulty = EXCLUDED.difficulty, updated_at = NOW()
		RETURNING id
	`, uuid.New().String(), curriculumID, slug, topic.Name, topic.Description, order, diff).Scan(&topicID)

	return topicID, err
}

func (e *Engine) upsertSubtopic(topicID string, subtopic SubtopicData, order int) (string, error) {
	slug := Slugify(subtopic.Name)

	var subtopicID string
	err := database.DB.QueryRow(`
		INSERT INTO subtopics (id, topic_id, slug, name, course_code, credit_units, semester, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (topic_id, slug)
		DO UPDATE SET name = EXCLUDED.name, course_code = EXCLUDED.course_code, credit_units = EXCLUDED.credit_units, semester = EXCLUDED.semester, order_index = EXCLUDED.order_index, updated_at = NOW()
		RETURNING id
	`, uuid.New().String(), topicID, slug, subtopic.Name, subtopic.CourseCode, subtopic.CreditUnits, subtopic.Semester, order).Scan(&subtopicID)

	return subtopicID, err
}

func (e *Engine) upsertObjective(subtopicID, description string, order int) error {
	description = strings.TrimSpace(description)
	verb := ExtractBloomsVerb(description)

	var exists bool
	err := database.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM learning_objectives WHERE subtopic_id = $1 AND description = $2)
	`, subtopicID, description).Scan(&exists)
	if err == nil && exists {
		return nil
	}

	_, err = database.DB.Exec(`
		INSERT INTO learning_objectives (id, subtopic_id, description, verb, order_index)
		VALUES ($1, $2, $3, $4, $5)
	`, uuid.New().String(), subtopicID, description, verb, order)
	return err
}

// Slugify converts any string to a URL-safe slug.
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, "–", "-")
	s = strings.ReplaceAll(s, "—", "-")
	s = strings.ReplaceAll(s, " & ", "-and-")
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	for _, r := range []string{",", "(", ")", "'", "\"", ".", ":", "?"} {
		s = strings.ReplaceAll(s, r, "")
	}
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}

// ClassifyDifficulty assigns a difficulty based on topic title keywords.
func ClassifyDifficulty(topicName string) string {
	hardKeywords := []string{"calculus", "matrices", "vectors", "quantum", "nuclear", "electronics", "electromagnetism", "genetics", "organic"}
	mediumKeywords := []string{"trigonometry", "coordinate", "statistics", "probability", "algebra", "waves", "optics", "thermodynamics", "ecology"}
	lower := strings.ToLower(topicName)
	for _, kw := range hardKeywords {
		if strings.Contains(lower, kw) {
			return "hard"
		}
	}
	for _, kw := range mediumKeywords {
		if strings.Contains(lower, kw) {
			return "medium"
		}
	}
	return "easy"
}

// ExtractBloomsVerb extracts the leading Bloom's Taxonomy action verb.
func ExtractBloomsVerb(objective string) string {
	bloomsVerbs := []string{
		"define", "identify", "calculate", "solve", "apply", "analyze", "analyse",
		"evaluate", "create", "construct", "design", "build", "interpret", "explain", "distinguish",
		"state", "use", "find", "draw", "describe", "determine", "perform", "express",
		"compare", "differentiate", "illustrate", "demonstrate", "relate", "formulate",
		"list", "name", "recall",
	}
	lower := strings.ToLower(strings.TrimSpace(objective))
	for _, verb := range bloomsVerbs {
		if strings.HasPrefix(lower, verb) {
			return verb
		}
	}
	return "understand"
}
