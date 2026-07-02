package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// LLMPromptResponse represents the formatted response tailored for LLM / AI Tutor ingestion
type LLMPromptResponse struct {
	ExamBoard         string            `json:"exam_board"`
	ExamBoardSlug     string            `json:"exam_board_slug"`
	Subject           string            `json:"subject"`
	SubjectSlug       string            `json:"subject_slug"`
	Level             string            `json:"level"`
	SystemPrompt      string            `json:"system_prompt"`
	TopicsSummary     string            `json:"topics_summary"`
	FullContextWindow string            `json:"full_context_window"`
	FormattedModules  []LLMModuleBlock `json:"formatted_modules"`
}

type LLMModuleBlock struct {
	ModuleName     string   `json:"module_name"`
	Difficulty     string   `json:"difficulty"`
	LLMInstruction string   `json:"llm_instruction"`
	Subtopics      []string `json:"subtopics"`
}

// GetLLMPrompt formats full curriculum into LLM System Prompt & Context Window
// GET /api/v1/curriculum/:board/:subject/llm-prompt
func GetLLMPrompt(c *gin.Context) {
	boardSlug := c.Param("board")
	subjectSlug := c.Param("subject")

	// 1. Fetch curriculum metadata
	var curr models.Curriculum
	var board models.ExamBoard
	var subject models.Subject

	err := database.DB.QueryRow(`
		SELECT 
			c.id, c.exam_board_id, c.subject_id, c.year, c.level, c.source_url, c.created_at, c.updated_at,
			eb.slug, eb.name, eb.full_name, eb.country, eb.description, eb.website,
			s.slug, s.name, s.description, s.category
		FROM curricula c
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s ON c.subject_id = s.id
		WHERE eb.slug = $1 AND s.slug = $2
		ORDER BY c.year DESC
		LIMIT 1
	`, boardSlug, subjectSlug).Scan(
		&curr.ID, &curr.ExamBoardID, &curr.SubjectID, &curr.Year, &curr.Level, &curr.SourceURL, &curr.CreatedAt, &curr.UpdatedAt,
		&board.Slug, &board.Name, &board.FullName, &board.Country, &board.Description, &board.Website,
		&subject.Slug, &subject.Name, &subject.Description, &subject.Category,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: fmt.Sprintf("Curriculum not found for %s/%s", boardSlug, subjectSlug),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch curriculum for LLM prompt generation",
		})
		return
	}

	// 2. Fetch topics
	topicRows, err := database.DB.Query(`
		SELECT id, curriculum_id, slug, name, description, order_index, difficulty, created_at, updated_at
		FROM topics WHERE curriculum_id = $1 ORDER BY order_index
	`, curr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch curriculum topics",
		})
		return
	}
	defer topicRows.Close()

	var topics []models.Topic
	var topicIDs []string
	for topicRows.Next() {
		var t models.Topic
		if err := topicRows.Scan(&t.ID, &t.CurriculumID, &t.Slug, &t.Name, &t.Description, &t.OrderIndex, &t.Difficulty, &t.CreatedAt, &t.UpdatedAt); err == nil {
			topics = append(topics, t)
			topicIDs = append(topicIDs, t.ID)
		}
	}

	// 3. Batch query subtopics
	subtopicMap := make(map[string]*models.Subtopic)
	var subtopicIDs []string
	if len(topicIDs) > 0 {
		subRows, err := database.DB.Query(`
			SELECT id, topic_id, slug, name, description, order_index, created_at, updated_at
			FROM subtopics WHERE topic_id = ANY($1) ORDER BY topic_id, order_index
		`, pq.Array(topicIDs))
		if err == nil {
			defer subRows.Close()
			for subRows.Next() {
				var st models.Subtopic
				if err := subRows.Scan(&st.ID, &st.TopicID, &st.Slug, &st.Name, &st.Description, &st.OrderIndex, &st.CreatedAt, &st.UpdatedAt); err == nil {
					subtopicIDs = append(subtopicIDs, st.ID)
					subtopicMap[st.ID] = &st
				}
			}
		}
	}

	// 4. Batch query objectives
	if len(subtopicIDs) > 0 {
		objRows, err := database.DB.Query(`
			SELECT id, subtopic_id, description, verb, order_index, created_at
			FROM learning_objectives WHERE subtopic_id = ANY($1) ORDER BY subtopic_id, order_index
		`, pq.Array(subtopicIDs))
		if err == nil {
			defer objRows.Close()
			for objRows.Next() {
				var obj models.LearningObjective
				if err := objRows.Scan(&obj.ID, &obj.SubtopicID, &obj.Description, &obj.Verb, &obj.OrderIndex, &obj.CreatedAt); err == nil {
					if st, exists := subtopicMap[obj.SubtopicID]; exists {
						st.Objectives = append(st.Objectives, obj)
					}
				}
			}
		}
	}

	// Group subtopics into topics
	topicSubtopicMap := make(map[string][]models.Subtopic)
	for _, stPtr := range subtopicMap {
		topicSubtopicMap[stPtr.TopicID] = append(topicSubtopicMap[stPtr.TopicID], *stPtr)
	}

	for i := range topics {
		if subs, ok := topicSubtopicMap[topics[i].ID]; ok {
			topics[i].Subtopics = subs
		} else {
			topics[i].Subtopics = []models.Subtopic{}
		}
	}

	// Level-specific instruction rules
	var levelRule string
	switch strings.ToLower(board.Slug) {
	case "bece", "nerdc":
		levelRule = "Stricly use simple, beginner-friendly explanations appropriate for Junior Secondary (JSS1-JSS3) students. Use simple word equations (e.g. Carbon Dioxide + Water -> Glucose + Oxygen) instead of balanced chemical formulas unless explicitly requested. Focus on foundational concepts without overloading with senior secondary or university details."
	case "waec", "neco", "jamb":
		levelRule = "Provide comprehensive Senior Secondary (SS1-SS3 / UTME) depth aligned with WAEC/NECO marking schemes. Use standard scientific notation, balanced chemical equations, and exam-style practice questions."
	default:
		levelRule = "Provide advanced university-level / polytechnic-level depth matching the official NUC/NBTE degree benchmarks."
	}

	systemPrompt := fmt.Sprintf(
		"You are an expert AI Tutor specialized in the official %s (%s) %s curriculum (%s level). "+
			"Your primary instruction is to explain concepts, solve practice problems, and guide students strictly aligned with "+
			"the official %s syllabus standards. %s Always provide clear, step-by-step explanations with relevant African examples.",
		board.Name, board.FullName, subject.Name, curr.Level, board.Name, levelRule,
	)

	var topicsSummaryBuilder strings.Builder
	var fullContextBuilder strings.Builder
	var moduleBlocks []LLMModuleBlock

	fullContextBuilder.WriteString(fmt.Sprintf("# %s — %s Official Curriculum Context\n\n", board.Name, subject.Name))
	fullContextBuilder.WriteString(fmt.Sprintf("**Level**: %s | **Category**: %s | **Source**: %s\n\n", curr.Level, subject.Category, curr.SourceURL))
	fullContextBuilder.WriteString("## System Directive for AI Tutor\n")
	fullContextBuilder.WriteString(systemPrompt + "\n\n")
	fullContextBuilder.WriteString("## Complete Syllabus Breakdown & Learning Objectives\n\n")

	for i, t := range topics {
		topicsSummaryBuilder.WriteString(fmt.Sprintf("%d. %s (%d subtopics)\n", i+1, t.Name, len(t.Subtopics)))

		fullContextBuilder.WriteString(fmt.Sprintf("### Module %d: %s\n", i+1, t.Name))
		if t.Description != "" {
			fullContextBuilder.WriteString(fmt.Sprintf("*Description*: %s\n\n", t.Description))
		}

		var subtopicNames []string
		for j, st := range t.Subtopics {
			subtopicNames = append(subtopicNames, st.Name)
			fullContextBuilder.WriteString(fmt.Sprintf("#### Unit %d.%d: %s\n", i+1, j+1, st.Name))
			if len(st.Objectives) > 0 {
				fullContextBuilder.WriteString("Learning Objectives:\n")
				for _, obj := range st.Objectives {
					fullContextBuilder.WriteString(fmt.Sprintf("- [%s] %s\n", strings.ToUpper(obj.Verb), obj.Description))
				}
			}
			fullContextBuilder.WriteString("\n")
		}

		moduleBlocks = append(moduleBlocks, LLMModuleBlock{
			ModuleName:     t.Name,
			Difficulty:     t.Difficulty,
			LLMInstruction: fmt.Sprintf("Teach '%s' with focus on: %s", t.Name, strings.Join(subtopicNames, ", ")),
			Subtopics:      subtopicNames,
		})
	}

	response := LLMPromptResponse{
		ExamBoard:         board.Name,
		ExamBoardSlug:     board.Slug,
		Subject:           subject.Name,
		SubjectSlug:       subject.Slug,
		Level:             curr.Level,
		SystemPrompt:      systemPrompt,
		TopicsSummary:     topicsSummaryBuilder.String(),
		FullContextWindow: fullContextBuilder.String(),
		FormattedModules:  moduleBlocks,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
		Meta: &models.Meta{
			Source:  curr.SourceURL,
			Version: "v1",
		},
	})
}
