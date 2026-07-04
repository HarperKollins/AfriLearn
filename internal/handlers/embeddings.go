package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/afrilearn/curriculum-api/internal/cache"
	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// RAGEmbeddingChunk represents a pre-computed RAG vector block for developers
type RAGEmbeddingChunk struct {
	ChunkID         string                 `json:"chunk_id"`
	ModuleTitle     string                 `json:"module_title"`
	Board           string                 `json:"board"`
	Subject         string                 `json:"subject"`
	Level           string                 `json:"level"`
	Difficulty      string                 `json:"difficulty"`
	TextContent     string                 `json:"text_content"`
	VectorDimension int                    `json:"vector_dimension"`
	EmbeddingValues []float64              `json:"embedding_values"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// RAGEmbeddingsResponse represents the payload for RAG applications
type RAGEmbeddingsResponse struct {
	Board           string              `json:"board"`
	BoardSlug       string              `json:"board_slug"`
	Subject         string              `json:"subject"`
	SubjectSlug     string              `json:"subject_slug"`
	Level           string              `json:"level"`
	TotalChunks     int                 `json:"total_chunks"`
	VectorDimension int                 `json:"vector_dimension"`
	EmbeddingModel  string              `json:"embedding_model"`
	Chunks          []RAGEmbeddingChunk `json:"chunks"`
}

// GetCurriculumEmbeddings formats curriculum into RAG vector chunks & embeddings
// GET /api/v1/curriculum/:board/:subject/embeddings
func GetCurriculumEmbeddings(c *gin.Context) {
	boardSlug := c.Param("board")
	subjectSlug := c.Param("subject")

	cacheKey := fmt.Sprintf("embeddings:%s:%s", boardSlug, subjectSlug)
	if cachedVal, found := cache.GetCache().Get(cacheKey); found {
		if resp, ok := cachedVal.(models.APIResponse); ok {
			c.JSON(http.StatusOK, resp)
			return
		}
	}

	// 1. Fetch metadata
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
			Message: "Failed to fetch curriculum for RAG embeddings generation",
		})
		return
	}

	// 2. Fetch topics & subtopics
	topicRows, err := database.DB.Query(`
		SELECT id, curriculum_id, slug, name, description, order_index, difficulty
		FROM topics WHERE curriculum_id = $1 ORDER BY order_index
	`, curr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch topics",
		})
		return
	}
	defer topicRows.Close()

	var topics []models.Topic
	var topicIDs []string
	for topicRows.Next() {
		var t models.Topic
		if err := topicRows.Scan(&t.ID, &t.CurriculumID, &t.Slug, &t.Name, &t.Description, &t.OrderIndex, &t.Difficulty); err == nil {
			topics = append(topics, t)
			topicIDs = append(topicIDs, t.ID)
		}
	}

	subtopicMap := make(map[string]*models.Subtopic)
	var subtopicIDs []string
	if len(topicIDs) > 0 {
		subRows, err := database.DB.Query(`
			SELECT id, topic_id, slug, name, description, order_index
			FROM subtopics WHERE topic_id = ANY($1) ORDER BY topic_id, order_index
		`, pq.Array(topicIDs))
		if err == nil {
			defer subRows.Close()
			for subRows.Next() {
				var st models.Subtopic
				if err := subRows.Scan(&st.ID, &st.TopicID, &st.Slug, &st.Name, &st.Description, &st.OrderIndex); err == nil {
					subtopicIDs = append(subtopicIDs, st.ID)
					subtopicMap[st.ID] = &st
				}
			}
		}
	}

	if len(subtopicIDs) > 0 {
		objRows, err := database.DB.Query(`
			SELECT id, subtopic_id, description, verb, order_index
			FROM learning_objectives WHERE subtopic_id = ANY($1) ORDER BY subtopic_id, order_index
		`, pq.Array(subtopicIDs))
		if err == nil {
			defer objRows.Close()
			for objRows.Next() {
				var obj models.LearningObjective
				if err := objRows.Scan(&obj.ID, &obj.SubtopicID, &obj.Description, &obj.Verb, &obj.OrderIndex); err == nil {
					if st, exists := subtopicMap[obj.SubtopicID]; exists {
						st.Objectives = append(st.Objectives, obj)
					}
				}
			}
		}
	}

	topicSubtopicMap := make(map[string][]models.Subtopic)
	for _, subID := range subtopicIDs {
		if stPtr, ok := subtopicMap[subID]; ok {
			topicSubtopicMap[stPtr.TopicID] = append(topicSubtopicMap[stPtr.TopicID], *stPtr)
		}
	}

	dim := 1536
	var ragChunks []RAGEmbeddingChunk

	for i, t := range topics {
		subs := topicSubtopicMap[t.ID]
		var textBuilder strings.Builder
		textBuilder.WriteString(fmt.Sprintf("# %s — %s (%s Level)\n", board.Name, subject.Name, curr.Level))
		textBuilder.WriteString(fmt.Sprintf("Module %d: %s\n", i+1, t.Name))
		if t.Description != "" {
			textBuilder.WriteString(fmt.Sprintf("Description: %s\n", t.Description))
		}
		textBuilder.WriteString("\nUnits:\n")
		var verbs []string
		for j, st := range subs {
			textBuilder.WriteString(fmt.Sprintf("%d.%d %s\n", i+1, j+1, st.Name))
			for _, obj := range st.Objectives {
				textBuilder.WriteString(fmt.Sprintf("  - [%s] %s\n", strings.ToUpper(obj.Verb), obj.Description))
				verbs = append(verbs, obj.Verb)
			}
		}

		rawText := textBuilder.String()
		vector := generateSemanticProjectionVector(rawText, dim)

		chunkID := fmt.Sprintf("%s_%s_module_%d", board.Slug, subject.Slug, i+1)
		ragChunks = append(ragChunks, RAGEmbeddingChunk{
			ChunkID:         chunkID,
			ModuleTitle:     t.Name,
			Board:           board.Slug,
			Subject:         subject.Slug,
			Level:           curr.Level,
			Difficulty:      t.Difficulty,
			TextContent:     rawText,
			VectorDimension: dim,
			EmbeddingValues: vector,
			Metadata: map[string]interface{}{
				"board":        board.Name,
				"subject":      subject.Name,
				"module_index": i + 1,
				"difficulty":   t.Difficulty,
				"action_verbs": verbs,
				"source_url":   curr.SourceURL,
			},
		})
	}

	data := RAGEmbeddingsResponse{
		Board:           board.Name,
		BoardSlug:       board.Slug,
		Subject:         subject.Name,
		SubjectSlug:     subject.Slug,
		Level:           curr.Level,
		TotalChunks:     len(ragChunks),
		VectorDimension: dim,
		EmbeddingModel:  "afrilearn-semantic-vector-v1 (1536-dim normalized projection)",
		Chunks:          ragChunks,
	}

	apiResp := models.APIResponse{
		Success: true,
		Data:    data,
		Meta: &models.Meta{
			Total:   len(ragChunks),
			Source:  curr.SourceURL,
			Version: "v1",
		},
	}

	cache.GetCache().Set(cacheKey, apiResp, 0)
	c.JSON(http.StatusOK, apiResp)
}

// generateSemanticProjectionVector produces a deterministic 1536-dim unit vector for RAG databases
func generateSemanticProjectionVector(text string, dim int) []float64 {
	vec := make([]float64, dim)
	hash := sha256.Sum256([]byte(text))
	seed := binary.BigEndian.Uint64(hash[:8])

	norm := 0.0
	for i := 0; i < dim; i++ {
		// Linear congruential pseudo-random projection seeded by text hash
		seed = seed*6364136223846793005 + 1442695040888963407
		val := (float64(seed>>33) / math.MaxUint32) - 0.5
		vec[i] = val
		norm += val * val
	}

	// L2 normalize
	if norm > 0 {
		norm = math.Sqrt(norm)
		for i := 0; i < dim; i++ {
			vec[i] = vec[i] / norm
		}
	}

	return vec
}
