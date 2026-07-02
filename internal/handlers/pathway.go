package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/models"
	"github.com/gin-gonic/gin"
)

// PathwayStep is one step in a curriculum learning journey
type PathwayStep struct {
	Stage      int    `json:"stage"`
	Board      string `json:"board"`
	BoardName  string `json:"board_name"`
	Level      string `json:"level"`
	TopicName  string `json:"topic_name"`
	TopicSlug  string `json:"topic_slug"`
	Difficulty string `json:"difficulty"`
	Depth      string `json:"depth"`
}

// GetLearningPathway returns an ordered step-by-step learning journey between two boards for a subject
// GET /api/v1/curriculum/pathway?subject=mathematics&from=bece&to=jamb
func GetLearningPathway(c *gin.Context) {
	subjectSlug := strings.TrimSpace(c.Query("subject"))
	fromBoard := strings.TrimSpace(strings.ToLower(c.Query("from")))
	toBoard := strings.TrimSpace(strings.ToLower(c.Query("to")))

	if subjectSlug == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Query param 'subject' is required. Example: ?subject=mathematics&from=bece&to=jamb",
		})
		return
	}

	// If no from/to specified, return the full journey across all boards
	var boardFilter string
	var queryArgs []interface{}
	queryArgs = append(queryArgs, subjectSlug) // $1

	if fromBoard != "" && toBoard != "" {
		fromOrder := boardLevelOrder[fromBoard]
		toOrder := boardLevelOrder[toBoard]
		if fromOrder == 0 || toOrder == 0 {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: fmt.Sprintf("Unknown board slug in from/to. Valid boards: bece, waec, jamb, nuc, yabatech, imt, unilag, futo, ebsu, funai, unn, unec, etc."),
			})
			return
		}
		if fromOrder > toOrder {
			fromOrder, toOrder = toOrder, fromOrder
		}
		boardFilter = fmt.Sprintf(" AND eb.slug = ANY($2) ")
		// Build list of board slugs in range
		inRange := []string{}
		for slug, order := range boardLevelOrder {
			if order >= fromOrder && order <= toOrder {
				inRange = append(inRange, slug)
			}
		}
		queryArgs = append(queryArgs, postgresTextArray(inRange)) // $2
	}

	query := fmt.Sprintf(`
		SELECT
			t.name, t.slug, t.difficulty, t.order_index,
			eb.slug   AS board_slug,
			eb.name   AS board_name,
			eb.full_name,
			c.level
		FROM topics t
		JOIN curricula c     ON t.curriculum_id = c.id
		JOIN exam_boards eb  ON c.exam_board_id = eb.id
		JOIN subjects s      ON c.subject_id = s.id
		WHERE s.slug = $1 %s
		ORDER BY
			COALESCE((
				SELECT order_val FROM (VALUES
					('bece',1),('nerdc',2),('neco',3),('waec',4),('jamb',5),
					('nbte',6),('yabatech',7),('imt',8),('auchi',9),('fedpoly-nek',10),
					('nuc',11),('ebsu',12),('funai',13),('unec',14),('unn',15),
					('unilag',16),('ui',17),('oau',18),('abu',19),('futo',20),('futa',21),('covenant',22)
				) AS board_order(board_slug, order_val)
				WHERE board_slug = eb.slug
			), 99),
			t.order_index
	`, boardFilter)

	rows, err := database.DB.Query(query, queryArgs...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Pathway query failed: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	var steps []PathwayStep
	stage := 0
	for rows.Next() {
		var (
			topicName, topicSlug, difficulty string
			orderIndex                       int
			boardSlug, boardName, boardFull  string
			level                            string
		)
		if err := rows.Scan(&topicName, &topicSlug, &difficulty, &orderIndex,
			&boardSlug, &boardName, &boardFull, &level); err != nil {
			continue
		}
		stage++
		steps = append(steps, PathwayStep{
			Stage:      stage,
			Board:      boardSlug,
			BoardName:  boardFull,
			Level:      level,
			TopicName:  topicName,
			TopicSlug:  topicSlug,
			Difficulty: difficulty,
			Depth:      boardDepthLabel(boardSlug),
		})
	}

	if len(steps) == 0 {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: fmt.Sprintf("No pathway found for subject '%s'", subjectSlug),
		})
		return
	}

	// Build summary
	boardsSeen := []string{}
	boardSeenMap := make(map[string]bool)
	for _, s := range steps {
		if !boardSeenMap[s.Board] {
			boardsSeen = append(boardsSeen, s.BoardName)
			boardSeenMap[s.Board] = true
		}
	}
	progressionSummary := strings.Join(boardsSeen, " → ")

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: fmt.Sprintf("Learning pathway for '%s': %d steps across %d board(s)", subjectSlug, len(steps), len(boardsSeen)),
		Data: gin.H{
			"subject":             subjectSlug,
			"from":                fromBoard,
			"to":                  toBoard,
			"total_steps":         len(steps),
			"boards_covered":      len(boardsSeen),
			"progression_summary": progressionSummary,
			"pathway":             steps,
		},
		Meta: &models.Meta{
			Total:   len(steps),
			Version: "v1",
		},
	})
}

// GetTopicPrerequisites returns prerequisite topics for a specific topic
// GET /api/v1/curriculum/prerequisites/:board/:subject/:topic
func GetTopicPrerequisites(c *gin.Context) {
	boardSlug := c.Param("board")
	subjectSlug := c.Param("subject")
	topicSlug := c.Param("topic")

	// Find the topic first
	var topicID, topicName string
	err := database.DB.QueryRow(`
		SELECT t.id, t.name
		FROM topics t
		JOIN curricula c    ON t.curriculum_id = c.id
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s     ON c.subject_id = s.id
		WHERE eb.slug = $1 AND s.slug = $2 AND t.slug = $3
		LIMIT 1
	`, boardSlug, subjectSlug, topicSlug).Scan(&topicID, &topicName)

	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: fmt.Sprintf("Topic '%s' not found for %s/%s", topicSlug, boardSlug, subjectSlug),
		})
		return
	}

	// Fetch prerequisites from topic_prerequisites table
	rows, err := database.DB.Query(`
		SELECT
			t.name, t.slug, t.difficulty,
			eb.slug AS board_slug, eb.full_name AS board_name,
			c.level,
			tp.order_index
		FROM topic_prerequisites tp
		JOIN topics t           ON tp.prerequisite_topic_id = t.id
		JOIN curricula c        ON t.curriculum_id = c.id
		JOIN exam_boards eb     ON c.exam_board_id = eb.id
		WHERE tp.topic_id = $1
		ORDER BY tp.order_index
	`, topicID)

	var prereqs []gin.H
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name, slug, diff, bSlug, bName, level string
			var orderIdx int
			if err := rows.Scan(&name, &slug, &diff, &bSlug, &bName, &level, &orderIdx); err != nil {
				continue
			}
			prereqs = append(prereqs, gin.H{
				"topic_name":  name,
				"topic_slug":  slug,
				"difficulty":  diff,
				"board":       bSlug,
				"board_name":  bName,
				"level":       level,
				"order_index": orderIdx,
			})
		}
	}

	// If no explicit prerequisites exist yet, derive them from the pathway (topics that come before this one in the same subject)
	if len(prereqs) == 0 {
		prereqs = []gin.H{}
		derivedRows, err2 := database.DB.Query(`
			SELECT t2.name, t2.slug, t2.difficulty, eb2.slug, eb2.full_name, c2.level
			FROM topics t2
			JOIN curricula c2    ON t2.curriculum_id = c2.id
			JOIN exam_boards eb2 ON c2.exam_board_id = eb2.id
			JOIN subjects s2     ON c2.subject_id = s2.id
			JOIN curricula c1    ON c1.subject_id = s2.id
			JOIN topics t1       ON t1.curriculum_id = c1.id AND t1.id = $1
			WHERE t2.order_index < (SELECT order_index FROM topics WHERE id = $1)
			  AND c2.subject_id = c1.subject_id
			ORDER BY t2.order_index DESC
			LIMIT 5
		`, topicID)
		if err2 == nil {
			defer derivedRows.Close()
			for derivedRows.Next() {
				var name, slug, diff, bSlug, bName, level string
				if err := derivedRows.Scan(&name, &slug, &diff, &bSlug, &bName, &level); err != nil {
					continue
				}
				prereqs = append(prereqs, gin.H{
					"topic_name": name,
					"topic_slug": slug,
					"difficulty": diff,
					"board":      bSlug,
					"board_name": bName,
					"level":      level,
					"derived":    true,
				})
			}
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"topic":         topicName,
			"topic_slug":    topicSlug,
			"board":         boardSlug,
			"subject":       subjectSlug,
			"prerequisites": prereqs,
			"total":         len(prereqs),
		},
		Meta: &models.Meta{Version: "v1"},
	})
}

// postgresTextArray converts a Go string slice to a pq-compatible array literal
func postgresTextArray(items []string) string {
	if len(items) == 0 {
		return "{}"
	}
	quoted := make([]string, len(items))
	for i, s := range items {
		quoted[i] = `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
	}
	return "{" + strings.Join(quoted, ",") + "}"
}
