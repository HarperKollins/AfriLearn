package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/models"
	"github.com/gin-gonic/gin"
)

// boardLevelOrder defines the curriculum hierarchy for progression sorting
var boardLevelOrder = map[string]int{
	"bece":       1,
	"nerdc":      2,
	"neco":       3,
	"waec":       4,
	"jamb":       5,
	"nbte":       6,
	"yabatech":   7,
	"imt":        8,
	"auchi":      9,
	"fedpoly-nek": 10,
	"nuc":        11,
	"ebsu":       12,
	"funai":      13,
	"unec":       14,
	"unn":        15,
	"unilag":     16,
	"ui":         17,
	"oau":        18,
	"abu":        19,
	"futo":       20,
	"futa":       21,
	"covenant":   22,
}

// boardDepthLabel returns a human-readable depth label for a board
func boardDepthLabel(slug string) string {
	switch slug {
	case "bece", "nerdc":
		return "introductory"
	case "waec", "neco":
		return "intermediate"
	case "jamb":
		return "exam-focused"
	case "nbte", "yabatech", "imt", "auchi", "fedpoly-nek":
		return "polytechnic"
	default:
		return "university"
	}
}

// CrossBoardMatch represents a single board's coverage of a topic
type CrossBoardMatch struct {
	Board          string       `json:"board"`
	BoardFullName  string       `json:"board_full_name"`
	Level          string       `json:"level"`
	Depth          string       `json:"depth"`
	Subject        string       `json:"subject"`
	SubjectSlug    string       `json:"subject_slug"`
	MatchedTopics  []TopicMatch `json:"matched_topics"`
	TopicCount     int          `json:"topic_count"`
}

// TopicMatch is a simplified topic for cross-board results
type TopicMatch struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Difficulty string `json:"difficulty"`
}

// GetCurriculumMatch queries ALL boards simultaneously for a topic keyword
// GET /api/v1/curriculum/match/:topic
func GetCurriculumMatch(c *gin.Context) {
	topicSlug := c.Param("topic")
	topicSlug = strings.ReplaceAll(topicSlug, "-", " ")

	if topicSlug == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Topic parameter is required",
		})
		return
	}

	// Phase 2 FTS Fix: use plainto_tsquery so GIN indexes are actually used.
	// This replaces the old LIKE which did a full sequential scan even with GIN indexes.
	rows, err := database.DB.Query(`
		SELECT
			t.id, t.name, t.slug, t.difficulty,
			eb.slug      AS board_slug,
			eb.name      AS board_name,
			eb.full_name AS board_full_name,
			c.level,
			s.name       AS subject_name,
			s.slug       AS subject_slug
		FROM topics t
		JOIN curricula c    ON t.curriculum_id = c.id
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s     ON c.subject_id = s.id
		WHERE to_tsvector('english', t.name) @@ plainto_tsquery('english', $1)
		ORDER BY
			ts_rank(
				to_tsvector('english', t.name),
				plainto_tsquery('english', $1)
			) DESC,
			eb.slug, t.order_index
	`, topicSlug)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Cross-board query failed: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	// Aggregate results by board slug
	boardMap := make(map[string]*CrossBoardMatch)
	boardOrder := []string{}

	for rows.Next() {
		var (
			topicID, topicName, topicSlugDB, difficulty string
			boardSlug, boardName, boardFullName, level   string
			subjectName, subjectSlug                     string
		)
		if err := rows.Scan(
			&topicID, &topicName, &topicSlugDB, &difficulty,
			&boardSlug, &boardName, &boardFullName, &level,
			&subjectName, &subjectSlug,
		); err != nil {
			continue
		}

		if _, exists := boardMap[boardSlug]; !exists {
			boardMap[boardSlug] = &CrossBoardMatch{
				Board:         boardSlug,
				BoardFullName: boardFullName,
				Level:         level,
				Depth:         boardDepthLabel(boardSlug),
				Subject:       subjectName,
				SubjectSlug:   subjectSlug,
				MatchedTopics: []TopicMatch{},
			}
			boardOrder = append(boardOrder, boardSlug)
		}

		boardMap[boardSlug].MatchedTopics = append(boardMap[boardSlug].MatchedTopics, TopicMatch{
			ID:         topicID,
			Name:       topicName,
			Slug:       topicSlugDB,
			Difficulty: difficulty,
		})
		boardMap[boardSlug].TopicCount++
	}

	if len(boardMap) == 0 {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: fmt.Sprintf("No curriculum found matching '%s' across any board", topicSlug),
		})
		return
	}

	// Sort matches by curriculum hierarchy (BECE → WAEC → JAMB → NUC → universities)
	coverage := make([]CrossBoardMatch, 0, len(boardMap))
	// Build sorted by boardLevelOrder
	type boardWithOrder struct {
		slug  string
		order int
	}
	sorted := []boardWithOrder{}
	for _, slug := range boardOrder {
		o := boardLevelOrder[slug]
		if o == 0 {
			o = 99
		}
		sorted = append(sorted, boardWithOrder{slug, o})
	}
	// Simple insertion sort (small slice)
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && sorted[j].order < sorted[j-1].order; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}
	for _, b := range sorted {
		coverage = append(coverage, *boardMap[b.slug])
	}

	// Build progression path string
	progressionParts := []string{}
	for _, m := range coverage {
		progressionParts = append(progressionParts, fmt.Sprintf("%s (%s)", m.BoardFullName, m.Level))
	}
	progressionPath := strings.Join(progressionParts, " → ")

	// Generate unified LLM prompt spanning all boards
	var promptBuilder strings.Builder
	promptBuilder.WriteString(fmt.Sprintf(
		"You are an expert AI Tutor covering '%s' across ALL Nigerian curriculum levels. ",
		strings.Title(topicSlug),
	))
	promptBuilder.WriteString("You have deep knowledge spanning:\n")
	for _, m := range coverage {
		promptBuilder.WriteString(fmt.Sprintf(
			"- %s (%s level, %s depth): %d matched topic(s)\n",
			m.BoardFullName, m.Level, m.Depth, m.TopicCount,
		))
	}
	promptBuilder.WriteString("\nAlways teach progressively: start from the most foundational level and build up.")

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: fmt.Sprintf("Found '%s' across %d board(s)", strings.Title(topicSlug), len(coverage)),
		Data: gin.H{
			"topic":                topicSlug,
			"total_boards_matched": len(coverage),
			"cross_board_coverage": coverage,
			"progression_path":     progressionPath,
			"llm_unified_prompt":   promptBuilder.String(),
		},
		Meta: &models.Meta{
			Total:   len(coverage),
			Version: "v1",
		},
	})
}
