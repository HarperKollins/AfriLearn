package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/afrilearn/curriculum-api/internal/cache"
	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// ──────────────────────────────────────────────────────────────────────────────
// Bridge cache TTL
// ──────────────────────────────────────────────────────────────────────────────

const bridgeCacheTTL = 30 * time.Minute

// ──────────────────────────────────────────────────────────────────────────────
// Bridge Request / Response shapes
// ──────────────────────────────────────────────────────────────────────────────

// BridgeRequest is what HK AI (or any EdTech platform) POSTs to this endpoint.
// Provide (board + subject) for direct lookup, or messages[] for auto-detection.
type BridgeRequest struct {
	// Explicit slugs — skip detection if you already know them
	Board   string `json:"board"`   // e.g. "waec", "jamb", "unilag"
	Subject string `json:"subject"` // e.g. "physics", "law", "mathematics"

	// Free-text conversation messages for auto-detection when board/subject omitted.
	// Shape: [{ "role": "user"|"assistant", "content": "..." }]
	Messages []BridgeMessage `json:"messages"`

	// Optional: workspace title gives extra detection signal
	WorkspaceTitle string `json:"workspace_title"`
}

// BridgeMessage is a single chat turn
type BridgeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// BridgeResponse is the self-contained payload HK AI injects into its AI call.
type BridgeResponse struct {
	// Resolved board + subject
	Board         string `json:"board"`
	Subject       string `json:"subject"`
	BoardFullName string `json:"board_full_name"`
	SubjectName   string `json:"subject_name"`
	Level         string `json:"level"`
	// "explicit" | "high" | "medium" | "low" | "none"
	DetectionScore string `json:"detection_score"`

	// LLM enrichment data
	SystemPrompt          string         `json:"system_prompt,omitempty"`
	SubjectSpecificRules  []string       `json:"subject_specific_rules,omitempty"`
	PedagogicalDirectives []string       `json:"pedagogical_directives,omitempty"`
	BloomsTaxonomy        map[string]int `json:"blooms_taxonomy,omitempty"`
	MisconceptionFlags    []string       `json:"misconception_flags,omitempty"`

	// Structured topic tree — lets HK AI organise course lessons around official syllabus
	Topics []BridgeTopic `json:"topics,omitempty"`

	// FTS hits from conversation content
	SearchHits []BridgeSearchHit `json:"search_hits,omitempty"`

	// Drop-in prompt: paste this string directly into your Groq/Gemini system message
	InjectionPrompt string `json:"injection_prompt"`
}

// BridgeTopic is a top-level curriculum topic with its subtopics
type BridgeTopic struct {
	Name       string      `json:"name"`
	Slug       string      `json:"slug"`
	Difficulty string      `json:"difficulty"`
	Subtopics  []BridgeSub `json:"subtopics,omitempty"`
}

// BridgeSub is a subtopic within a topic
type BridgeSub struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// BridgeSearchHit is a single FTS result from the conversation
type BridgeSearchHit struct {
	Board    string  `json:"board"`
	Subject  string  `json:"subject"`
	Topic    string  `json:"topic"`
	Subtopic string  `json:"subtopic,omitempty"`
	Rank     float64 `json:"rank"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Board / Subject detector — reuses the maps in query.go (same package)
// ──────────────────────────────────────────────────────────────────────────────

// bridgeDetect infers board + subject from concatenated message text + title.
// boardKeywords and subjectKeywords are declared in query.go.
func bridgeDetect(messages []BridgeMessage, title string) (board, subject, score string) {
	parts := []string{strings.ToLower(title)}
	for _, m := range messages {
		parts = append(parts, strings.ToLower(m.Content))
	}
	text := strings.Join(parts, " ")

	// board — check explicit board names first (highest priority), then inferred keywords.
	// Ordered slice ensures deterministic priority: direct board slugs win over grade-level inferences.
	explicitBoards := []string{"waec", "jamb", "bece", "neco", "unilag", "unn", "unec", "ebsu", "funai", "futo", "yabatech", "imt", "nuc", "nbte"}
	for _, slug := range explicitBoards {
		if strings.Contains(text, slug) {
			board = slug
			break
		}
	}
	// Fallback: scan boardKeywords map for inferred phrases (e.g. "ss1", "university level")
	if board == "" {
		for kw, slug := range boardKeywords {
			if strings.Contains(text, strings.ToLower(kw)) {
				board = slug
				break
			}
		}
	}

	// subject — try topicToSubject first (more precise), then subjectKeywords
	for kw, slug := range topicToSubject {
		if strings.Contains(text, strings.ToLower(kw)) {
			subject = slug
			break
		}
	}
	if subject == "" {
		for kw, slug := range subjectKeywords {
			if strings.Contains(text, strings.ToLower(kw)) {
				subject = slug
				break
			}
		}
	}

	switch {
	case board != "" && subject != "":
		score = "high"
	case subject != "":
		score = "medium"
	case board != "":
		score = "low"
	default:
		score = "none"
	}
	return
}

// ──────────────────────────────────────────────────────────────────────────────
// HandleBridgeEnrich — POST /api/v1/bridge/enrich
// ──────────────────────────────────────────────────────────────────────────────

// HandleBridgeEnrich is the single AfriLearn→HK AI integration point.
//
// HK AI sends its workspace conversation + optional board/subject slugs.
// AfriLearn returns:
//   - system_prompt  — full LLM system prompt for the detected curriculum
//   - topics         — official topic tree for lesson structuring
//   - search_hits    — FTS matches from the conversation
//   - injection_prompt — ready-to-paste string for Groq/Gemini system message
func HandleBridgeEnrich(c *gin.Context) {
	var req BridgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	if req.Board == "" && req.Subject == "" && len(req.Messages) == 0 && req.WorkspaceTitle == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Provide at least one of: board+subject, messages[], or workspace_title",
		})
		return
	}

	resp := BridgeResponse{}

	// ── 1. Resolve board + subject ─────────────────────────────────────────────
	if req.Board != "" && req.Subject != "" {
		resp.Board = strings.ToLower(strings.TrimSpace(req.Board))
		resp.Subject = strings.ToLower(strings.TrimSpace(req.Subject))
		resp.DetectionScore = "explicit"
	} else {
		board, subject, score := bridgeDetect(req.Messages, req.WorkspaceTitle)
		resp.Board = board
		resp.Subject = subject
		resp.DetectionScore = score
	}

	// ── 2. Fetch static curriculum data (cached per board+subject) ─────────────
	if resp.Board != "" && resp.Subject != "" {
		cacheKey := fmt.Sprintf("bridge:%s:%s", resp.Board, resp.Subject)

		if cached, found := cache.GetCache().Get(cacheKey); found {
			if cr, ok := cached.(BridgeResponse); ok {
				cr.DetectionScore = resp.DetectionScore
				cr.SearchHits = bridgeFetchSearchHits(req.Messages)
				cr.InjectionPrompt = bridgeBuildInjectionPrompt(cr)
				c.JSON(http.StatusOK, models.APIResponse{
					Success: true,
					Data:    cr,
					Meta:    &models.Meta{Source: "cache", Version: "v1"},
				})
				return
			}
		}

		// Fetch board + subject metadata + curriculum ID
		var board models.ExamBoard
		var subject models.Subject
		var curr models.Curriculum

		err := database.DB.QueryRow(`
			SELECT c.id, eb.slug, eb.name, eb.full_name, s.slug, s.name, c.level
			FROM curricula c
			JOIN exam_boards eb ON c.exam_board_id = eb.id
			JOIN subjects s ON c.subject_id = s.id
			WHERE eb.slug = $1 AND s.slug = $2
			LIMIT 1
		`, resp.Board, resp.Subject).Scan(
			&curr.ID,
			&board.Slug, &board.Name, &board.FullName,
			&subject.Slug, &subject.Name,
			&curr.Level,
		)

		if err == nil {
			resp.BoardFullName = board.FullName
			resp.SubjectName = subject.Name
			resp.Level = curr.Level
			resp.Topics = bridgeFetchTopics(curr.ID)
			bridgeFetchLLMData(resp.Board, resp.Subject, &resp)
		} else if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Database error fetching curriculum",
			})
			return
		}

		// Cache the static parts (search hits are per-request)
		staticCopy := resp
		staticCopy.SearchHits = nil
		staticCopy.InjectionPrompt = ""
		cache.GetCache().Set(cacheKey, staticCopy, bridgeCacheTTL)
	}

	// ── 3. Conversation-specific search hits ───────────────────────────────────
	resp.SearchHits = bridgeFetchSearchHits(req.Messages)

	// ── 4. Build the ready-to-inject prompt string ────────────────────────────
	resp.InjectionPrompt = bridgeBuildInjectionPrompt(resp)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    resp,
		Meta:    &models.Meta{Source: "db", Version: "v1"},
	})
}

// ── bridgeFetchTopics pulls topic + subtopic tree for a curriculum ─────────────
func bridgeFetchTopics(curriculumID string) []BridgeTopic {
	rows, err := database.DB.Query(`
		SELECT id, slug, name, COALESCE(difficulty, 'intermediate')
		FROM topics
		WHERE curriculum_id = $1
		ORDER BY order_index
	`, curriculumID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var topics []BridgeTopic
	for rows.Next() {
		var t BridgeTopic
		var topicID string
		if err := rows.Scan(&topicID, &t.Slug, &t.Name, &t.Difficulty); err != nil {
			continue
		}

		srows, serr := database.DB.Query(`
			SELECT slug, name FROM subtopics
			WHERE topic_id = $1
			ORDER BY order_index
		`, topicID)
		if serr == nil {
			defer srows.Close()
			for srows.Next() {
				var sub BridgeSub
				if serr := srows.Scan(&sub.Slug, &sub.Name); serr == nil {
					t.Subtopics = append(t.Subtopics, sub)
				}
			}
		}
		topics = append(topics, t)
	}
	return topics
}

// ── bridgeFetchLLMData pulls system prompt + pedagogy from DB ─────────────────
func bridgeFetchLLMData(boardSlug, subjectSlug string, resp *BridgeResponse) {
	// Try the pre-built llm_prompts table first
	var systemPrompt, topicsSummary string
	var subjectRules, pedagogicalDirectives, misconceptionFlags pq.StringArray

	err := database.DB.QueryRow(`
		SELECT
			COALESCE(system_prompt, '') as system_prompt,
			COALESCE(topics_summary, '') as topics_summary,
			COALESCE(subject_specific_rules, ARRAY[]::text[]) as subject_specific_rules,
			COALESCE(pedagogical_directives, ARRAY[]::text[]) as pedagogical_directives,
			COALESCE(misconception_flags, ARRAY[]::text[]) as misconception_flags
		FROM llm_prompts
		WHERE board_slug = $1 AND subject_slug = $2
		LIMIT 1
	`, boardSlug, subjectSlug).Scan(
		&systemPrompt,
		&topicsSummary,
		&subjectRules,
		&pedagogicalDirectives,
		&misconceptionFlags,
	)

	if err == nil && systemPrompt != "" {
		resp.SystemPrompt = systemPrompt
		resp.SubjectSpecificRules = []string(subjectRules)
		resp.PedagogicalDirectives = []string(pedagogicalDirectives)
		resp.MisconceptionFlags = []string(misconceptionFlags)
	} else {
		// Fallback: build from topic tree
		resp.SystemPrompt = bridgeBuildFallbackPrompt(resp)
	}

	// Bloom's taxonomy breakdown from learning objectives
	bloomsRows, err := database.DB.Query(`
		SELECT lo.verb, COUNT(*) as count
		FROM learning_objectives lo
		JOIN subtopics st ON lo.subtopic_id = st.id
		JOIN topics t ON st.topic_id = t.id
		JOIN curricula c ON t.curriculum_id = c.id
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s ON c.subject_id = s.id
		WHERE eb.slug = $1 AND s.slug = $2
		  AND lo.verb IS NOT NULL AND lo.verb != ''
		GROUP BY lo.verb
	`, boardSlug, subjectSlug)
	if err == nil {
		defer bloomsRows.Close()
		resp.BloomsTaxonomy = make(map[string]int)
		for bloomsRows.Next() {
			var verb string
			var count int
			if err := bloomsRows.Scan(&verb, &count); err == nil {
				resp.BloomsTaxonomy[verb] = count
			}
		}
	}
}

// ── bridgeFetchSearchHits does FTS from the last 3 user messages ──────────────
func bridgeFetchSearchHits(messages []BridgeMessage) []BridgeSearchHit {
	if len(messages) == 0 {
		return nil
	}

	var userTexts []string
	count := 0
	for i := len(messages) - 1; i >= 0 && count < 3; i-- {
		if messages[i].Role == "user" && len(strings.TrimSpace(messages[i].Content)) > 3 {
			userTexts = append(userTexts, messages[i].Content)
			count++
		}
	}
	if len(userTexts) == 0 {
		return nil
	}

	q := strings.Join(userTexts, " ")
	if len(q) > 400 {
		q = q[:400]
	}

	rows, err := database.DB.Query(`
		SELECT
			eb.slug,
			s.slug,
			t.name,
			st.name,
			ts_rank(
				to_tsvector('english', t.name || ' ' || COALESCE(t.description,'') || ' ' || st.name),
				plainto_tsquery('english', $1)
			) AS rank
		FROM topics t
		JOIN subtopics st ON st.topic_id = t.id
		JOIN curricula c ON t.curriculum_id = c.id
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s ON c.subject_id = s.id
		WHERE to_tsvector('english', t.name || ' ' || COALESCE(t.description,'') || ' ' || st.name)
		      @@ plainto_tsquery('english', $1)
		ORDER BY rank DESC
		LIMIT 5
	`, q)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var hits []BridgeSearchHit
	for rows.Next() {
		var h BridgeSearchHit
		if err := rows.Scan(&h.Board, &h.Subject, &h.Topic, &h.Subtopic, &h.Rank); err == nil {
			hits = append(hits, h)
		}
	}
	return hits
}

// ── bridgeBuildFallbackPrompt builds a system prompt from the topic tree ───────
func bridgeBuildFallbackPrompt(resp *BridgeResponse) string {
	if resp.BoardFullName == "" {
		return ""
	}
	var sb strings.Builder
	fmt.Fprintf(&sb,
		"You are an expert AI Tutor specialized in the official %s %s curriculum for Nigerian students.\n\n",
		resp.BoardFullName, resp.SubjectName,
	)
	if len(resp.Topics) > 0 {
		sb.WriteString("OFFICIAL CURRICULUM TOPICS:\n")
		for i, t := range resp.Topics {
			fmt.Fprintf(&sb, "%d. %s", i+1, t.Name)
			if len(t.Subtopics) > 0 {
				subs := make([]string, len(t.Subtopics))
				for j, s := range t.Subtopics {
					subs[j] = s.Name
				}
				sb.WriteString(" → " + strings.Join(subs, ", "))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("Ground every response in this official syllabus. " +
		"Use Nigerian analogies, WAEC/JAMB exam conventions, and local context. " +
		"Always reference the correct exam board format.\n")
	return sb.String()
}

// ── bridgeBuildInjectionPrompt assembles the final paste-ready prompt ──────────
func bridgeBuildInjectionPrompt(resp BridgeResponse) string {
	if resp.SystemPrompt == "" && len(resp.Topics) == 0 {
		return ""
	}
	var sb strings.Builder

	sb.WriteString("=== AFRILEARN CURRICULUM INTELLIGENCE CONTEXT ===\n")

	if resp.SystemPrompt != "" {
		sb.WriteString(resp.SystemPrompt)
		sb.WriteString("\n\n")
	}

	if len(resp.SubjectSpecificRules) > 0 {
		sb.WriteString("SUBJECT-SPECIFIC RULES:\n")
		for _, r := range resp.SubjectSpecificRules {
			fmt.Fprintf(&sb, "• %s\n", r)
		}
		sb.WriteString("\n")
	}

	if len(resp.MisconceptionFlags) > 0 {
		sb.WriteString("COMMON MISCONCEPTIONS TO ADDRESS PROACTIVELY:\n")
		for _, m := range resp.MisconceptionFlags {
			fmt.Fprintf(&sb, "⚠️  %s\n", m)
		}
		sb.WriteString("\n")
	}

	if len(resp.Topics) > 0 {
		sb.WriteString("OFFICIAL TOPIC STRUCTURE (use this to organise course lessons):\n")
		for i, t := range resp.Topics {
			if i >= 12 {
				fmt.Fprintf(&sb, "  ... and %d more topics\n", len(resp.Topics)-12)
				break
			}
			fmt.Fprintf(&sb, "%d. %s", i+1, t.Name)
			if len(t.Subtopics) > 0 && len(t.Subtopics) <= 5 {
				subs := make([]string, len(t.Subtopics))
				for j, s := range t.Subtopics {
					subs[j] = s.Name
				}
				fmt.Fprintf(&sb, " (%s)", strings.Join(subs, ", "))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	if len(resp.SearchHits) > 0 {
		sb.WriteString("RELEVANT CURRICULUM CONTEXT FROM CONVERSATION:\n")
		for _, h := range resp.SearchHits {
			fmt.Fprintf(&sb, "• [%s %s] %s", strings.ToUpper(h.Board), h.Subject, h.Topic)
			if h.Subtopic != "" {
				fmt.Fprintf(&sb, " → %s", h.Subtopic)
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("=== END AFRILEARN CONTEXT ===\n")
	return sb.String()
}

// ──────────────────────────────────────────────────────────────────────────────
// HandleBridgeBoards — GET /api/v1/bridge/boards
// Returns all boards + their available subjects for HK AI's workspace picker UI
// ──────────────────────────────────────────────────────────────────────────────

// BridgeBoardEntry represents one board with its available subjects
type BridgeBoardEntry struct {
	Slug     string   `json:"slug"`
	Name     string   `json:"name"`
	FullName string   `json:"full_name"`
	Level    string   `json:"level"`
	Subjects []string `json:"subjects"`
}

// HandleBridgeBoards returns all available boards + subjects.
// HK AI uses this to build the workspace creation picker
// (e.g. a dropdown: "WAEC → Physics, Chemistry, Mathematics...").
func HandleBridgeBoards(c *gin.Context) {
	cacheKey := "bridge:boards"
	if cached, found := cache.GetCache().Get(cacheKey); found {
		if boards, ok := cached.([]BridgeBoardEntry); ok {
			c.JSON(http.StatusOK, models.APIResponse{
				Success: true,
				Data:    boards,
				Meta:    &models.Meta{Source: "cache", Version: "v1"},
			})
			return
		}
	}

	rows, err := database.DB.Query(`
		SELECT
			eb.slug,
			eb.name,
			eb.full_name,
			ARRAY_AGG(DISTINCT s.slug ORDER BY s.slug) AS subject_slugs
		FROM curricula c
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s ON c.subject_id = s.id
		GROUP BY eb.slug, eb.name, eb.full_name
		ORDER BY eb.name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch boards",
		})
		return
	}
	defer rows.Close()

	var boards []BridgeBoardEntry
	for rows.Next() {
		var b BridgeBoardEntry
		var subjects pq.StringArray
		if err := rows.Scan(&b.Slug, &b.Name, &b.FullName, &subjects); err != nil {
			continue
		}
		b.Subjects = []string(subjects)
		boards = append(boards, b)
	}

	cache.GetCache().Set(cacheKey, boards, bridgeCacheTTL)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    boards,
		Meta:    &models.Meta{Total: len(boards), Version: "v1"},
	})
}
