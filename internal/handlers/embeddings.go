package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/afrilearn/curriculum-api/internal/cache"
	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// RAGEmbeddingChunk represents a semantic text chunk ready for embedding into a vector database.
// The embedding_values field contains a DETERMINISTIC PLACEHOLDER vector derived from a hash of
// the text content. It is NOT a real ML/semantic embedding. To use real semantic search, replace
// embedding_values with vectors from OpenAI text-embedding-3-small, Google text-embedding-004,
// or any compatible model. See ROADMAP.md for integration details.
type RAGEmbeddingChunk struct {
	ChunkID         string                 `json:"chunk_id"`
	ModuleTitle     string                 `json:"module_title"`
	Board           string                 `json:"board"`
	Subject         string                 `json:"subject"`
	Level           string                 `json:"level"`
	Difficulty      string                 `json:"difficulty"`
	TextContent     string                 `json:"text_content"`
	TokenEstimate   int                    `json:"token_estimate"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// RAGEmbeddingsResponse represents curriculum content pre-chunked for RAG/vector database ingestion.
type RAGEmbeddingsResponse struct {
	Board              string              `json:"board"`
	BoardSlug          string              `json:"board_slug"`
	Subject            string              `json:"subject"`
	SubjectSlug        string              `json:"subject_slug"`
	Level              string              `json:"level"`
	TotalChunks        int                 `json:"total_chunks"`
	TotalTokenEstimate int                 `json:"total_token_estimate"`
	ChunkingStrategy   string              `json:"chunking_strategy"`
	EmbeddingNote      string              `json:"embedding_note"`
	IntegrationGuide   map[string]string   `json:"integration_guide"`
	Chunks             []RAGEmbeddingChunk `json:"chunks"`
}

// GetCurriculumEmbeddings formats curriculum into RAG-ready text chunks for vector database ingestion.
// The chunks contain pre-formatted text content (topics + subtopics + objectives) ready to be
// embedded with any embedding model. The endpoint does NOT generate real ML embeddings — it provides
// the structured text payload that you embed with your chosen model (OpenAI, Gemini, Ollama, etc.)
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
			Message: "Failed to fetch curriculum for RAG chunk generation",
		})
		return
	}

	// 2. Fetch topics
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

	var ragChunks []RAGEmbeddingChunk
	totalTokens := 0

	for i, t := range topics {
		subs := topicSubtopicMap[t.ID]
		var textBuilder strings.Builder

		// Header context — critical for RAG retrieval relevance
		textBuilder.WriteString(fmt.Sprintf("CURRICULUM: %s — %s\n", board.FullName, subject.Name))
		textBuilder.WriteString(fmt.Sprintf("LEVEL: %s | BOARD: %s | CATEGORY: %s\n\n", curr.Level, board.Name, subject.Category))
		textBuilder.WriteString(fmt.Sprintf("MODULE %d: %s\n", i+1, t.Name))
		textBuilder.WriteString(fmt.Sprintf("DIFFICULTY: %s\n", strings.ToUpper(t.Difficulty)))
		if t.Description != "" {
			textBuilder.WriteString(fmt.Sprintf("OVERVIEW: %s\n", t.Description))
		}
		textBuilder.WriteString("\nLEARNING UNITS:\n")

		var verbs []string
		var objectiveCount int
		for j, st := range subs {
			textBuilder.WriteString(fmt.Sprintf("\n%d.%d %s\n", i+1, j+1, st.Name))
			if st.Description != "" {
				textBuilder.WriteString(fmt.Sprintf("  Overview: %s\n", st.Description))
			}
			if len(st.Objectives) > 0 {
				textBuilder.WriteString("  Learning Objectives:\n")
				for _, obj := range st.Objectives {
					textBuilder.WriteString(fmt.Sprintf("  - [%s] %s\n", strings.ToUpper(obj.Verb), obj.Description))
					verbs = append(verbs, obj.Verb)
					objectiveCount++
				}
			}
		}

		rawText := textBuilder.String()
		// Rough token estimate: ~0.75 tokens per character for English educational text
		tokenEst := int(float64(len(rawText)) * 0.75 / 4)
		totalTokens += tokenEst

		chunkID := fmt.Sprintf("%s_%s_module_%02d", board.Slug, subject.Slug, i+1)
		ragChunks = append(ragChunks, RAGEmbeddingChunk{
			ChunkID:       chunkID,
			ModuleTitle:   t.Name,
			Board:         board.Slug,
			Subject:       subject.Slug,
			Level:         curr.Level,
			Difficulty:    t.Difficulty,
			TextContent:   rawText,
			TokenEstimate: tokenEst,
			Metadata: map[string]interface{}{
				"board_full_name": board.FullName,
				"subject_full":    subject.Name,
				"module_index":    i + 1,
				"difficulty":      t.Difficulty,
				"subtopic_count":  len(subs),
				"objective_count": objectiveCount,
				"action_verbs":    uniqueStrings(verbs),
				"source_url":      curr.SourceURL,
				"curriculum_year": curr.Year,
			},
		})
	}

	data := RAGEmbeddingsResponse{
		Board:              board.Name,
		BoardSlug:          board.Slug,
		Subject:            subject.Name,
		SubjectSlug:        subject.Slug,
		Level:              curr.Level,
		TotalChunks:        len(ragChunks),
		TotalTokenEstimate: totalTokens,
		ChunkingStrategy:   "one-chunk-per-topic-module (200–800 tokens per chunk, optimized for RAG retrieval)",
		EmbeddingNote: "Text chunks are ready for embedding. Pass text_content to your embedding model " +
			"(OpenAI text-embedding-3-small, Google text-embedding-004, Ollama nomic-embed-text, etc.) " +
			"to generate real semantic vectors. Store with chunk_id as the vector ID.",
		IntegrationGuide: map[string]string{
			"openai":       "client.embeddings.create(model='text-embedding-3-small', input=chunk['text_content'])",
			"google_gemini": "genai.embed_content(model='models/text-embedding-004', content=chunk['text_content'])",
			"ollama":        "ollama.embeddings(model='nomic-embed-text', prompt=chunk['text_content'])",
			"pgvector_upsert": "INSERT INTO topics (id, embedding) VALUES ($1, $2::vector) ON CONFLICT (id) DO UPDATE SET embedding = EXCLUDED.embedding",
		},
		Chunks: ragChunks,
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

// uniqueStrings returns deduplicated string slice preserving order
func uniqueStrings(in []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(in))
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// VectorSearchRequest is the request body for the semantic search endpoint
type VectorSearchRequest struct {
	Query  string `json:"query" binding:"required"`
	Limit  int    `json:"limit"`
	Board  string `json:"board"`   // optional: filter by board slug
	Subject string `json:"subject"` // optional: filter by subject slug
}

// VectorSearchResult represents a single search result
type VectorSearchResult struct {
	TopicID         string  `json:"topic_id"`
	TopicName       string  `json:"topic_name"`
	TopicSlug       string  `json:"topic_slug"`
	BoardSlug       string  `json:"board_slug"`
	BoardName       string  `json:"board_name"`
	SubjectSlug     string  `json:"subject_slug"`
	SubjectName     string  `json:"subject_name"`
	Level           string  `json:"level"`
	Difficulty      string  `json:"difficulty"`
	Snippet         string  `json:"snippet"`
	MatchedIn       string  `json:"matched_in"` // "topic", "subtopic", or "objective"
	RelevanceScore  float64 `json:"relevance_score"`
}

// HandleVectorSearch performs full-text semantic search across topics, subtopics, and objectives.
// NOTE: This endpoint uses PostgreSQL GIN full-text search (FTS), not ML vector embeddings.
// The pgvector HNSW index is provisioned and schema-ready for future real embedding integration.
// To enable true semantic search: (1) call /embeddings to get text chunks, (2) embed with your
// chosen model, (3) upsert into the topics.embedding column via pgvector, then queries will
// automatically use HNSW cosine similarity. Until then, FTS provides excellent keyword ranking.
// POST /api/v1/search/vector
func HandleVectorSearch(c *gin.Context) {
	var req VectorSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body. Required: {\"query\": \"your search text\"}",
		})
		return
	}

	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 10
	}

	// Build optional filters
	var filterClauses []string
	var args []interface{}
	args = append(args, req.Query) // $1
	args = append(args, req.Query) // $2 (used twice in ts_rank)
	args = append(args, req.Limit) // $3

	argIdx := 4
	if req.Board != "" {
		filterClauses = append(filterClauses, fmt.Sprintf("AND eb.slug = $%d", argIdx))
		args = append(args, strings.ToLower(req.Board))
		argIdx++
	}
	if req.Subject != "" {
		filterClauses = append(filterClauses, fmt.Sprintf("AND s.slug = $%d", argIdx))
		args = append(args, strings.ToLower(req.Subject))
		argIdx++
	}
	filterSQL := strings.Join(filterClauses, " ")

	// Phase 1: Search topics by name + description (highest weight)
	topicQuery := fmt.Sprintf(`
		SELECT 
			t.id, t.name, t.slug, t.difficulty,
			eb.slug, eb.name,
			s.slug, s.name,
			c.level,
			COALESCE(t.description, t.name) AS snippet,
			'topic' AS matched_in,
			ts_rank_cd(
				to_tsvector('english', t.name || ' ' || COALESCE(t.description, '')),
				plainto_tsquery('english', $1)
			) AS score
		FROM topics t
		JOIN curricula c ON t.curriculum_id = c.id
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s ON c.subject_id = s.id
		WHERE to_tsvector('english', t.name || ' ' || COALESCE(t.description, '')) @@ plainto_tsquery('english', $2)
		%s
		ORDER BY score DESC
		LIMIT $3
	`, filterSQL)

	var results []VectorSearchResult
	resultIDs := make(map[string]bool)

	rows, err := database.DB.Query(topicQuery, args...)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var r VectorSearchResult
			if err := rows.Scan(
				&r.TopicID, &r.TopicName, &r.TopicSlug, &r.Difficulty,
				&r.BoardSlug, &r.BoardName,
				&r.SubjectSlug, &r.SubjectName,
				&r.Level, &r.Snippet, &r.MatchedIn, &r.RelevanceScore,
			); err == nil {
				results = append(results, r)
				resultIDs[r.TopicID] = true
			}
		}
	}

	// Phase 2: Search subtopics — join back to parent topic
	subArgs := []interface{}{req.Query, req.Query}
	subFilterClauses := []string{}
	subArgIdx := 3
	if req.Board != "" {
		subFilterClauses = append(subFilterClauses, fmt.Sprintf("AND eb.slug = $%d", subArgIdx))
		subArgs = append(subArgs, strings.ToLower(req.Board))
		subArgIdx++
	}
	if req.Subject != "" {
		subFilterClauses = append(subFilterClauses, fmt.Sprintf("AND s.slug = $%d", subArgIdx))
		subArgs = append(subArgs, strings.ToLower(req.Subject))
		subArgIdx++
	}
	remaining := req.Limit - len(results)
	if remaining <= 0 {
		remaining = req.Limit
	}
	subArgs = append(subArgs, remaining)
	subFilterSQL := strings.Join(subFilterClauses, " ")

	subQuery := fmt.Sprintf(`
		SELECT DISTINCT ON (t.id)
			t.id, t.name, t.slug, t.difficulty,
			eb.slug, eb.name,
			s.slug, s.name,
			c.level,
			st.name AS snippet,
			'subtopic' AS matched_in,
			ts_rank_cd(
				to_tsvector('english', st.name || ' ' || COALESCE(st.description, '')),
				plainto_tsquery('english', $1)
			) AS score
		FROM subtopics st
		JOIN topics t ON st.topic_id = t.id
		JOIN curricula c ON t.curriculum_id = c.id
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s ON c.subject_id = s.id
		WHERE to_tsvector('english', st.name || ' ' || COALESCE(st.description, '')) @@ plainto_tsquery('english', $2)
		%s
		ORDER BY t.id, score DESC
		LIMIT $%d
	`, subFilterSQL, subArgIdx)

	subRows, err := database.DB.Query(subQuery, subArgs...)
	if err == nil {
		defer subRows.Close()
		for subRows.Next() {
			var r VectorSearchResult
			if err := subRows.Scan(
				&r.TopicID, &r.TopicName, &r.TopicSlug, &r.Difficulty,
				&r.BoardSlug, &r.BoardName,
				&r.SubjectSlug, &r.SubjectName,
				&r.Level, &r.Snippet, &r.MatchedIn, &r.RelevanceScore,
			); err == nil && !resultIDs[r.TopicID] {
				results = append(results, r)
				resultIDs[r.TopicID] = true
			}
		}
	}

	// Phase 3: Search learning objectives — join back to parent topic
	objArgs := []interface{}{req.Query, req.Query}
	objFilterClauses := []string{}
	objArgIdx := 3
	if req.Board != "" {
		objFilterClauses = append(objFilterClauses, fmt.Sprintf("AND eb.slug = $%d", objArgIdx))
		objArgs = append(objArgs, strings.ToLower(req.Board))
		objArgIdx++
	}
	if req.Subject != "" {
		objFilterClauses = append(objFilterClauses, fmt.Sprintf("AND s.slug = $%d", objArgIdx))
		objArgs = append(objArgs, strings.ToLower(req.Subject))
		objArgIdx++
	}
	objArgs = append(objArgs, req.Limit)
	objFilterSQL := strings.Join(objFilterClauses, " ")

	objQuery := fmt.Sprintf(`
		SELECT DISTINCT ON (t.id)
			t.id, t.name, t.slug, t.difficulty,
			eb.slug, eb.name,
			s.slug, s.name,
			c.level,
			lo.description AS snippet,
			'objective' AS matched_in,
			ts_rank_cd(
				to_tsvector('english', lo.description),
				plainto_tsquery('english', $1)
			) AS score
		FROM learning_objectives lo
		JOIN subtopics st ON lo.subtopic_id = st.id
		JOIN topics t ON st.topic_id = t.id
		JOIN curricula c ON t.curriculum_id = c.id
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s ON c.subject_id = s.id
		WHERE to_tsvector('english', lo.description) @@ plainto_tsquery('english', $2)
		%s
		ORDER BY t.id, score DESC
		LIMIT $%d
	`, objFilterSQL, objArgIdx)

	objRows, err := database.DB.Query(objQuery, objArgs...)
	if err == nil {
		defer objRows.Close()
		for objRows.Next() {
			var r VectorSearchResult
			if err := objRows.Scan(
				&r.TopicID, &r.TopicName, &r.TopicSlug, &r.Difficulty,
				&r.BoardSlug, &r.BoardName,
				&r.SubjectSlug, &r.SubjectName,
				&r.Level, &r.Snippet, &r.MatchedIn, &r.RelevanceScore,
			); err == nil && !resultIDs[r.TopicID] {
				results = append(results, r)
				resultIDs[r.TopicID] = true
			}
		}
	}

	// Cap total results
	if len(results) > req.Limit {
		results = results[:req.Limit]
	}

	if results == nil {
		results = []VectorSearchResult{}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"query":         req.Query,
			"search_engine": "PostgreSQL GIN full-text search (pgvector HNSW slot reserved — see /embeddings for chunk extraction)",
			"search_scope":  "topics + subtopics + learning objectives",
			"filters_applied": gin.H{
				"board":   req.Board,
				"subject": req.Subject,
			},
			"total_matches": len(results),
			"results":       results,
			"upgrade_path": "To enable real semantic vector search: (1) GET /api/v1/curriculum/:board/:subject/embeddings " +
				"(2) embed text_content with OpenAI/Gemini/Ollama (3) upsert into topics.embedding column via pgvector",
		},
		Meta: &models.Meta{
			Total:   len(results),
			Version: "v1",
		},
	})
}
