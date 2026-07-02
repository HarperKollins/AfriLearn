package scraper

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/google/uuid"
)

// Engine manages execution of curriculum scrapers and data ingestion into PostgreSQL.
type Engine struct{}

// NewEngine creates a new Engine instance.
func NewEngine() *Engine {
	return &Engine{}
}

// Execute runs a given Scraper and ingests its full curriculum hierarchy into the database.
func (e *Engine) Execute(s Scraper) error {
	log.Printf("🚀 Starting ingestion for [%s / %s]...", s.BoardSlug(), s.SubjectSlug())

	// 1. Ensure Curriculum Record exists
	curriculumID, err := e.ensureCurriculum(s)
	if err != nil {
		return fmt.Errorf("failed to resolve curriculum for [%s/%s]: %w", s.BoardSlug(), s.SubjectSlug(), err)
	}
	log.Printf("  ✅ Curriculum ID resolved: %s", curriculumID)

	// 2. Extract Golden Record Topics
	topics := s.Topics()
	log.Printf("  📦 Ingesting %d topics for [%s/%s]...", len(topics), s.BoardSlug(), s.SubjectSlug())

	// 3. Insert/Upsert Topics, Subtopics & Objectives in DB
	for i, topic := range topics {
		topicID, err := e.upsertTopic(curriculumID, topic, i+1)
		if err != nil {
			log.Printf("  ⚠️ Failed to upsert topic '%s': %v", topic.Name, err)
			continue
		}

		for j, subtopic := range topic.Subtopics {
			subtopicID, err := e.upsertSubtopic(topicID, subtopic, j+1)
			if err != nil {
				log.Printf("    ⚠️ Failed to upsert subtopic '%s': %v", subtopic.Name, err)
				continue
			}

			for k, objective := range subtopic.Objectives {
				if err := e.upsertObjective(subtopicID, objective, k+1); err != nil {
					log.Printf("      ⚠️ Failed to insert objective [%d]: %v", k+1, err)
				}
			}
		}
		log.Printf("  ✅ [%d/%d] Ingested topic: '%s' (%d subtopics)", i+1, len(topics), topic.Name, len(topic.Subtopics))
	}

	log.Printf("🎉 Ingestion complete for [%s / %s]!\n", s.BoardSlug(), s.SubjectSlug())
	return nil
}

// ensureCurriculum guarantees exam board, subject, and curriculum entries exist in DB.
func (e *Engine) ensureCurriculum(s Scraper) (string, error) {
	var boardID string
	err := database.DB.QueryRow(`SELECT id FROM exam_boards WHERE slug = $1`, s.BoardSlug()).Scan(&boardID)
	if err != nil {
		return "", fmt.Errorf("exam board '%s' not found: %w", s.BoardSlug(), err)
	}

	var subjectID string
	err = database.DB.QueryRow(`SELECT id FROM subjects WHERE slug = $1`, s.SubjectSlug()).Scan(&subjectID)
	if err != nil {
		return "", fmt.Errorf("subject '%s' not found: %w", s.SubjectSlug(), err)
	}

	currentYear := time.Now().Year()

	var currID string
	err = database.DB.QueryRow(`
		SELECT id FROM curricula 
		WHERE exam_board_id = $1 AND subject_id = $2 AND year = $3
		LIMIT 1
	`, boardID, subjectID, currentYear).Scan(&currID)

	if err == nil {
		return currID, nil
	}

	if err != sql.ErrNoRows {
		return "", err
	}

	// Create new curriculum record
	currID = uuid.New().String()
	_, err = database.DB.Exec(`
		INSERT INTO curricula (id, exam_board_id, subject_id, year, level, source_url)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, currID, boardID, subjectID, currentYear, s.Level(), s.SourceURL())

	if err != nil {
		return "", fmt.Errorf("failed to insert curriculum: %w", err)
	}

	return currID, nil
}

func (e *Engine) upsertTopic(curriculumID string, topic TopicData, order int) (string, error) {
	slug := Slugify(topic.Name)
	diff := topic.Difficulty
	if diff == "" {
		diff = ClassifyDifficulty(topic.Name, "")
	}

	var topicID string
	err := database.DB.QueryRow(`
		SELECT id FROM topics WHERE curriculum_id = $1 AND slug = $2
	`, curriculumID, slug).Scan(&topicID)

	if err == nil {
		// Update existing
		_, err = database.DB.Exec(`
			UPDATE topics 
			SET name = $1, description = $2, order_index = $3, difficulty = $4, updated_at = NOW()
			WHERE id = $5
		`, topic.Name, topic.Description, order, diff, topicID)
		return topicID, err
	}

	// Insert new
	topicID = uuid.New().String()
	_, err = database.DB.Exec(`
		INSERT INTO topics (id, curriculum_id, slug, name, description, order_index, difficulty)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, topicID, curriculumID, slug, topic.Name, topic.Description, order, diff)

	return topicID, err
}

func (e *Engine) upsertSubtopic(topicID string, subtopic SubtopicData, order int) (string, error) {
	slug := Slugify(subtopic.Name)

	var subtopicID string
	err := database.DB.QueryRow(`
		SELECT id FROM subtopics WHERE topic_id = $1 AND slug = $2
	`, topicID, slug).Scan(&subtopicID)

	if err == nil {
		_, err = database.DB.Exec(`
			UPDATE subtopics 
			SET name = $1, order_index = $2, updated_at = NOW()
			WHERE id = $3
		`, subtopic.Name, order, subtopicID)
		return subtopicID, err
	}

	subtopicID = uuid.New().String()
	_, err = database.DB.Exec(`
		INSERT INTO subtopics (id, topic_id, slug, name, order_index)
		VALUES ($1, $2, $3, $4, $5)
	`, subtopicID, topicID, slug, subtopic.Name, order)

	return subtopicID, err
}

func (e *Engine) upsertObjective(subtopicID, description string, order int) error {
	verb := ExtractBloomsVerb(description)

	// Check if already exists for this subtopic
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
