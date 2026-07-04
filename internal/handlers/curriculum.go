package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/afrilearn/curriculum-api/internal/cache"
	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// GetAllSubjects returns all available subjects
// GET /api/v1/subjects
func GetAllSubjects(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, slug, name, description, category, created_at, updated_at
		FROM subjects
		ORDER BY category, name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch subjects",
		})
		return
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var s models.Subject
		if err := rows.Scan(&s.ID, &s.Slug, &s.Name, &s.Description, &s.Category, &s.CreatedAt, &s.UpdatedAt); err != nil {
			continue
		}
		subjects = append(subjects, s)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    subjects,
		Meta:    &models.Meta{Total: len(subjects), Version: "v1"},
	})
}

// GetSubjectBySlug returns a single subject by its slug
// GET /api/v1/subjects/:slug
func GetSubjectBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var s models.Subject
	err := database.DB.QueryRow(`
		SELECT id, slug, name, description, category, created_at, updated_at
		FROM subjects WHERE slug = $1
	`, slug).Scan(&s.ID, &s.Slug, &s.Name, &s.Description, &s.Category, &s.CreatedAt, &s.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Subject not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch subject",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    s,
		Meta:    &models.Meta{Version: "v1"},
	})
}

// GetAllExamBoards returns all exam boards
// GET /api/v1/exam-boards
func GetAllExamBoards(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, slug, name, full_name, country, description, website, created_at, updated_at
		FROM exam_boards ORDER BY name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch exam boards",
		})
		return
	}
	defer rows.Close()

	var boards []models.ExamBoard
	for rows.Next() {
		var b models.ExamBoard
		if err := rows.Scan(&b.ID, &b.Slug, &b.Name, &b.FullName, &b.Country, &b.Description, &b.Website, &b.CreatedAt, &b.UpdatedAt); err != nil {
			continue
		}
		boards = append(boards, b)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    boards,
		Meta:    &models.Meta{Total: len(boards), Version: "v1"},
	})
}

// GetCurriculum returns full curriculum with topics, subtopics, and objectives for a board+subject
// GET /api/v1/curriculum/:board/:subject
func GetCurriculum(c *gin.Context) {
	boardSlug := c.Param("board")
	subjectSlug := c.Param("subject")

	cacheKey := fmt.Sprintf("curr:%s:%s", boardSlug, subjectSlug)
	if cachedVal, found := cache.GetCache().Get(cacheKey); found {
		if resp, ok := cachedVal.(models.APIResponse); ok {
			c.JSON(http.StatusOK, resp)
			return
		}
	}

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
			Message: "Curriculum not found for " + boardSlug + "/" + subjectSlug,
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch curriculum",
		})
		return
	}

	curr.ExamBoard = &board
	curr.Subject = &subject

	// 2. Query all topics for this curriculum
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
	topicMap := make(map[string]*models.Topic)

	for topicRows.Next() {
		var t models.Topic
		if err := topicRows.Scan(&t.ID, &t.CurriculumID, &t.Slug, &t.Name, &t.Description, &t.OrderIndex, &t.Difficulty, &t.CreatedAt, &t.UpdatedAt); err != nil {
			continue
		}
		topics = append(topics, t)
		topicIDs = append(topicIDs, t.ID)
	}

	if len(topicIDs) == 0 {
		curr.Topics = []models.Topic{}
		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Data:    curr,
			Meta:    &models.Meta{Source: curr.SourceURL, Version: "v1"},
		})
		return
	}

	// 3. Query all subtopics for all topics in a single batch query
	subRows, err := database.DB.Query(`
		SELECT id, topic_id, slug, name, description, order_index, created_at, updated_at
		FROM subtopics WHERE topic_id = ANY($1) ORDER BY topic_id, order_index
	`, pq.Array(topicIDs))

	var subtopicIDs []string
	subtopicMap := make(map[string]*models.Subtopic)

	if err == nil {
		defer subRows.Close()
		for subRows.Next() {
			var st models.Subtopic
			if err := subRows.Scan(&st.ID, &st.TopicID, &st.Slug, &st.Name, &st.Description, &st.OrderIndex, &st.CreatedAt, &st.UpdatedAt); err != nil {
				continue
			}
			subtopicIDs = append(subtopicIDs, st.ID)
			subtopicMap[st.ID] = &st
		}
	}

	// 4. Query all learning_objectives for all subtopics in a single batch query
	if len(subtopicIDs) > 0 {
		objRows, err := database.DB.Query(`
			SELECT id, subtopic_id, description, verb, order_index, created_at
			FROM learning_objectives WHERE subtopic_id = ANY($1) ORDER BY subtopic_id, order_index
		`, pq.Array(subtopicIDs))
		if err == nil {
			defer objRows.Close()
			for objRows.Next() {
				var obj models.LearningObjective
				if err := objRows.Scan(&obj.ID, &obj.SubtopicID, &obj.Description, &obj.Verb, &obj.OrderIndex, &obj.CreatedAt); err != nil {
					continue
				}
				if st, exists := subtopicMap[obj.SubtopicID]; exists {
					st.Objectives = append(st.Objectives, obj)
				}
			}
		}
	}

	// Reassemble subtopics into topics in strict order_index order
	topicSubtopicMap := make(map[string][]models.Subtopic)
	for _, subID := range subtopicIDs {
		if stPtr, ok := subtopicMap[subID]; ok {
			topicSubtopicMap[stPtr.TopicID] = append(topicSubtopicMap[stPtr.TopicID], *stPtr)
		}
	}

	for i := range topics {
		if subs, ok := topicSubtopicMap[topics[i].ID]; ok {
			topics[i].Subtopics = subs
		} else {
			topics[i].Subtopics = []models.Subtopic{}
		}
		topicMap[topics[i].ID] = &topics[i]
	}

	curr.Topics = topics

	apiResp := models.APIResponse{
		Success: true,
		Data:    curr,
		Meta: &models.Meta{
			Source:  curr.SourceURL,
			Version: "v1",
		},
	}

	cache.GetCache().Set(cacheKey, apiResp, 0)
	c.JSON(http.StatusOK, apiResp)
}

// SearchTopics searches across topics and subtopics
// GET /api/v1/search?q=quadratic
func SearchTopics(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Query parameter 'q' is required",
		})
		return
	}


	rows, err := database.DB.Query(`
		SELECT 
			t.id, t.slug, t.name, t.difficulty,
			eb.slug as board_slug, eb.name as board_name,
			s.slug as subject_slug, s.name as subject_name
		FROM topics t
		JOIN curricula c ON t.curriculum_id = c.id
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s ON c.subject_id = s.id
		WHERE to_tsvector('english', t.name) @@ plainto_tsquery('english', $1)
		ORDER BY ts_rank(to_tsvector('english', t.name), plainto_tsquery('english', $1)) DESC
		LIMIT 50
	`, query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Search failed",
		})
		return
	}
	defer rows.Close()

	type SearchResult struct {
		ID          string `json:"id"`
		Slug        string `json:"slug"`
		Name        string `json:"name"`
		Difficulty  string `json:"difficulty"`
		BoardSlug   string `json:"board_slug"`
		BoardName   string `json:"board_name"`
		SubjectSlug string `json:"subject_slug"`
		SubjectName string `json:"subject_name"`
	}

	var results []SearchResult
	for rows.Next() {
		var r SearchResult
		if err := rows.Scan(&r.ID, &r.Slug, &r.Name, &r.Difficulty, &r.BoardSlug, &r.BoardName, &r.SubjectSlug, &r.SubjectName); err != nil {
			continue
		}
		results = append(results, r)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    results,
		Meta:    &models.Meta{Total: len(results), Version: "v1"},
	})
}

// HealthCheck returns the API health status
// GET /health
func HealthCheck(c *gin.Context) {
	dbStatus := "connected"
	if err := database.DB.Ping(); err != nil {
		dbStatus = "disconnected"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"version":  "v1",
		"database": dbStatus,
		"service":  "AfriLearn Curriculum API",
	})
}
