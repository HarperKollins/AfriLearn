package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// ──────────────────────────────────────────────────────────────────────────────
// Session store — holds partial query state for clarification loops
// ──────────────────────────────────────────────────────────────────────────────

type QuerySession struct {
	Intent    ParsedIntent
	ExpiresAt time.Time
}

var (
	sessionStore = make(map[string]QuerySession)
	sessionMu    sync.Mutex
)

func setSession(id string, s QuerySession) {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	sessionStore[id] = s
}

func getSession(id string) (QuerySession, bool) {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	s, ok := sessionStore[id]
	if ok && time.Now().After(s.ExpiresAt) {
		delete(sessionStore, id)
		return QuerySession{}, false
	}
	return s, ok
}

func deleteSession(id string) {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	delete(sessionStore, id)
}

// ──────────────────────────────────────────────────────────────────────────────
// Intent parser — keyword/regex, no LLM required
// ──────────────────────────────────────────────────────────────────────────────

// ParsedIntent holds extracted fields from a user query
type ParsedIntent struct {
	Board   string // e.g. "waec", "jamb", "nuc"
	Subject string // e.g. "mathematics", "physics"
	Topic   string // e.g. "quadratic-equations"
	Action  string // "curriculum", "topics", "llm-prompt", "match", "pathway", "prerequisites"
	Raw     string
}

var boardKeywords = map[string]string{
	"waec": "waec", "west african": "waec",
	"jamb": "jamb", "utme": "jamb", "joint admission": "jamb",
	"bece": "bece", "jss": "bece", "junior": "bece",
	"neco": "neco",
	"nuc":  "nuc", "university degree": "nuc", "degree": "nuc",
	"nbte": "nbte", "polytechnic": "nbte",
	"yabatech": "yabatech", "yaba": "yabatech",
	"imt": "imt", "institute of management": "imt",
	"unilag": "unilag", "university of lagos": "unilag",
	"unn": "unn", "university of nigeria nsukka": "unn",
	"unec": "unec", "university of nigeria enugu": "unec",
	"ebsu": "ebsu", "ebonyi": "ebsu",
	"funai": "funai", "funai ae-funai": "funai",
	"futo": "futo", "federal university of technology owerri": "futo",
	"oau": "oau", "obafemi awolowo": "oau",
	"ui": "ui", "university of ibadan": "ui",
	"abu": "abu", "ahmadu bello": "abu",
	"covenant": "covenant",
}

var subjectKeywords = map[string]string{
	"math":        "mathematics", "maths": "mathematics", "mathematics": "mathematics",
	"physics":     "physics",
	"chemistry":   "chemistry", "chem": "chemistry",
	"biology":     "biology", "bio": "biology",
	"economics":   "economics", "econ": "economics",
	"government":  "government", "govt": "government",
	"english":     "english-studies", "english studies": "english-studies", "literature": "literature-in-english",
	"computer science": "computer-science", "cs": "computer-science", "computing": "computer-science",
	"law":         "law",
	"accounting":  "accounting", "accounts": "accounting",
	"business":    "business-administration", "business admin": "business-administration",
	"nursing":     "nursing-science",
	"medicine":    "medicine-surgery", "mbbs": "medicine-surgery", "medical": "medicine-surgery",
	"mechanical":  "mechanical-engineering", "mechanical engineering": "mechanical-engineering",
	"electrical":  "electrical-engineering", "electrical engineering": "electrical-engineering",
	"petroleum":   "petroleum-engineering",
	"mass comm":   "mass-communication", "mass communication": "mass-communication", "journalism": "mass-communication",
	"social studies": "social-studies",
	"basic science": "basic-science",
	"basic technology": "basic-technology",
	"business studies": "business-studies",
	"science lab":  "science-laboratory-technology", "slt": "science-laboratory-technology",
	"computer engineering": "computer-engineering-technology",
}

var actionKeywords = map[string]string{
	"topics":        "topics",
	"curriculum":    "curriculum",
	"study":         "curriculum",
	"syllabus":      "curriculum",
	"ai tutor":      "llm-prompt",
	"llm":           "llm-prompt",
	"system prompt": "llm-prompt",
	"teach me":      "llm-prompt",
	"match":         "match",
	"across":        "match",
	"all boards":    "match",
	"pathway":       "pathway",
	"path":          "pathway",
	"journey":       "pathway",
	"progression":   "pathway",
	"order":         "pathway",
	"prerequisite":  "prerequisites",
	"before":        "prerequisites",
	"need to know":  "prerequisites",
}

func parseIntent(text string, existing *ParsedIntent) ParsedIntent {
	lower := strings.ToLower(text)
	intent := ParsedIntent{Raw: text}
	if existing != nil {
		intent = *existing
		intent.Raw = text
	}

	// Extract board
	if intent.Board == "" {
		for keyword, board := range boardKeywords {
			if strings.Contains(lower, keyword) {
				intent.Board = board
				break
			}
		}
	}

	// Extract subject
	if intent.Subject == "" {
		bestLen := 0
		for keyword, subject := range subjectKeywords {
			if strings.Contains(lower, keyword) && len(keyword) > bestLen {
				intent.Subject = subject
				bestLen = len(keyword)
			}
		}
	}

	// Extract action
	if intent.Action == "" {
		for keyword, action := range actionKeywords {
			if strings.Contains(lower, keyword) {
				intent.Action = action
				break
			}
		}
		if intent.Action == "" {
			intent.Action = "curriculum" // default
		}
	}

	return intent
}

func missingFields(intent ParsedIntent) []string {
	var missing []string
	if intent.Board == "" {
		missing = append(missing, "Which exam board or institution? (e.g. WAEC, JAMB, BECE, NECO, NUC University, YABATECH, UNILAG, FUTO)")
	}
	if intent.Subject == "" {
		missing = append(missing, "Which subject or degree programme? (e.g. Mathematics, Physics, Computer Science, Law, Medicine)")
	}
	return missing
}

// ──────────────────────────────────────────────────────────────────────────────
// Smart Query Cache
// ──────────────────────────────────────────────────────────────────────────────

func normalizeQuery(q string) string {
	q = strings.ToLower(strings.TrimSpace(q))
	// Strip filler words
	fillers := []string{"please", "can you", "tell me", "what", "give me", "show me", "i want", "do you have", "about", "the", "a ", "an "}
	for _, f := range fillers {
		q = strings.ReplaceAll(q, f, "")
	}
	// Collapse whitespace
	fields := strings.Fields(q)
	return strings.Join(fields, " ")
}

func hashQuery(normalized string) string {
	h := sha256.Sum256([]byte(normalized))
	return fmt.Sprintf("%x", h)[:32]
}

func cacheGet(hash string) (map[string]interface{}, bool) {
	var responseJSON []byte
	err := database.DB.QueryRow(
		`SELECT response_json FROM query_cache WHERE query_hash = $1`, hash,
	).Scan(&responseJSON)
	if err != nil {
		return nil, false
	}
	// Bump hit count asynchronously
	go func() {
		_, _ = database.DB.Exec(
			`UPDATE query_cache SET hit_count = hit_count + 1, last_hit_at = NOW() WHERE query_hash = $1`, hash,
		)
	}()
	var result map[string]interface{}
	if err := json.Unmarshal(responseJSON, &result); err != nil {
		return nil, false
	}
	return result, true
}

func cachePut(hash, rawQuery, normalized string, intentTags []string, response map[string]interface{}) {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return
	}
	_, _ = database.DB.Exec(`
		INSERT INTO query_cache (query_hash, raw_query, normalized, intent_tags, response_json)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (query_hash) DO UPDATE SET
			hit_count  = query_cache.hit_count + 1,
			last_hit_at = NOW()
	`, hash, rawQuery, normalized, pq.Array(intentTags), responseJSON)
}

// ──────────────────────────────────────────────────────────────────────────────
// Query Orchestrator — routes to correct handler based on intent
// ──────────────────────────────────────────────────────────────────────────────

func orchestrate(intent ParsedIntent) (map[string]interface{}, error) {
	switch intent.Action {
	case "match":
		return orchestrateMatch(intent)
	case "pathway":
		return orchestratePathway(intent)
	case "llm-prompt":
		return orchestrateLLMPrompt(intent)
	default:
		return orchestrateCurriculum(intent)
	}
}

func orchestrateCurriculum(intent ParsedIntent) (map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT t.name, t.slug, t.difficulty, eb.slug, eb.full_name, s.name, c.level
		FROM topics t
		JOIN curricula c    ON t.curriculum_id = c.id
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s     ON c.subject_id = s.id
		WHERE eb.slug = $1 AND s.slug = $2
		ORDER BY t.order_index
	`, intent.Board, intent.Subject)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type topicRow struct {
		Name       string `json:"name"`
		Slug       string `json:"slug"`
		Difficulty string `json:"difficulty"`
	}
	var topics []topicRow
	var boardFull, subjectName, level string
	for rows.Next() {
		var t topicRow
		var bSlug string
		if err := rows.Scan(&t.Name, &t.Slug, &t.Difficulty, &bSlug, &boardFull, &subjectName, &level); err != nil {
			continue
		}
		topics = append(topics, t)
	}
	if len(topics) == 0 {
		return nil, fmt.Errorf("no curriculum found for %s/%s", intent.Board, intent.Subject)
	}
	return map[string]interface{}{
		"board":       intent.Board,
		"board_name":  boardFull,
		"subject":     subjectName,
		"level":       level,
		"total_topics": len(topics),
		"topics":      topics,
	}, nil
}

func orchestrateMatch(intent ParsedIntent) (map[string]interface{}, error) {
	searchTerm := "%" + strings.ReplaceAll(intent.Subject, "-", " ") + "%"
	rows, err := database.DB.Query(`
		SELECT t.name, t.slug, t.difficulty, eb.slug, eb.full_name, c.level, s.slug
		FROM topics t
		JOIN curricula c    ON t.curriculum_id = c.id
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s     ON c.subject_id = s.id
		WHERE LOWER(s.slug) LIKE $1 OR LOWER(s.name) LIKE $1
		ORDER BY eb.slug, t.order_index
	`, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boardMap := make(map[string][]map[string]string)
	for rows.Next() {
		var name, slug, diff, bSlug, bFull, level, sSlug string
		if err := rows.Scan(&name, &slug, &diff, &bSlug, &bFull, &level, &sSlug); err != nil {
			continue
		}
		boardMap[bSlug] = append(boardMap[bSlug], map[string]string{
			"name": name, "slug": slug, "difficulty": diff,
			"board": bSlug, "board_name": bFull, "level": level,
		})
	}
	return map[string]interface{}{
		"subject":      intent.Subject,
		"boards_found": len(boardMap),
		"coverage":     boardMap,
	}, nil
}

func orchestratePathway(intent ParsedIntent) (map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT t.name, t.slug, t.difficulty, eb.slug, eb.full_name, c.level
		FROM topics t
		JOIN curricula c    ON t.curriculum_id = c.id
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s     ON c.subject_id = s.id
		WHERE s.slug = $1
		ORDER BY
			COALESCE((SELECT order_val FROM (VALUES
				('bece',1),('waec',4),('jamb',5),('nbte',6),('nuc',11)
			) AS o(slug, order_val) WHERE slug = eb.slug), 99),
			t.order_index
	`, intent.Subject)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pathway []map[string]interface{}
	stage := 0
	for rows.Next() {
		var name, slug, diff, bSlug, bFull, level string
		if err := rows.Scan(&name, &slug, &diff, &bSlug, &bFull, &level); err != nil {
			continue
		}
		stage++
		pathway = append(pathway, map[string]interface{}{
			"stage": stage, "topic": name, "board": bSlug,
			"board_name": bFull, "level": level, "difficulty": diff,
		})
	}
	return map[string]interface{}{
		"subject":     intent.Subject,
		"total_steps": len(pathway),
		"pathway":     pathway,
	}, nil
}

func orchestrateLLMPrompt(intent ParsedIntent) (map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT t.name, s.name, eb.full_name, c.level
		FROM topics t
		JOIN curricula c    ON t.curriculum_id = c.id
		JOIN exam_boards eb ON c.exam_board_id = eb.id
		JOIN subjects s     ON c.subject_id = s.id
		WHERE eb.slug = $1 AND s.slug = $2
		ORDER BY t.order_index
	`, intent.Board, intent.Subject)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topicNames []string
	var boardFull, subjectName, level string
	for rows.Next() {
		var tName, sName, bFull, lvl string
		if err := rows.Scan(&tName, &sName, &bFull, &lvl); err != nil {
			continue
		}
		topicNames = append(topicNames, tName)
		boardFull, subjectName, level = bFull, sName, lvl
	}
	if len(topicNames) == 0 {
		return nil, fmt.Errorf("no curriculum found for %s/%s", intent.Board, intent.Subject)
	}

	sysPrompt := fmt.Sprintf(
		"You are an expert AI Tutor for %s (%s) %s. "+
			"You have comprehensive knowledge of all %d official topics. "+
			"Teach clearly, use Nigerian examples, and follow the official curriculum order: %s.",
		boardFull, strings.ToUpper(intent.Board), subjectName,
		len(topicNames), strings.Join(topicNames[:min(10, len(topicNames))], ", "),
	)

	return map[string]interface{}{
		"board":         intent.Board,
		"subject":       subjectName,
		"level":         level,
		"total_topics":  len(topicNames),
		"system_prompt": sysPrompt,
		"topics":        topicNames,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ──────────────────────────────────────────────────────────────────────────────
// Main Handler
// ──────────────────────────────────────────────────────────────────────────────

type QueryRequest struct {
	Question  string `json:"question" binding:"required"`
	SessionID string `json:"session_id"`
}

// HandleCurriculumQuery is the Curriculum Query Brain
// POST /api/v1/query
func HandleCurriculumQuery(c *gin.Context) {
	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Request body must be JSON with a 'question' field",
		})
		return
	}

	rawQuestion := strings.TrimSpace(req.Question)
	normalized := normalizeQuery(rawQuestion)
	hash := hashQuery(normalized)

	// ── Stage 1: Cache check ─────────────────────────────────────────────────
	if cached, found := cacheGet(hash); found {
		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Message: "Answered from cache (pattern match)",
			Data: gin.H{
				"cache_hit":    true,
				"query":        rawQuestion,
				"result":       cached,
			},
			Meta: &models.Meta{Version: "v1"},
		})
		return
	}

	// ── Stage 2: Parse intent (merge with session if provided) ───────────────
	var existingIntent *ParsedIntent
	if req.SessionID != "" {
		if session, ok := getSession(req.SessionID); ok {
			existingIntent = &session.Intent
		}
	}
	intent := parseIntent(rawQuestion, existingIntent)

	// ── Stage 3: Clarification check ─────────────────────────────────────────
	missing := missingFields(intent)
	if len(missing) > 0 {
		// Store partial intent in session for follow-up
		sessionID := hash[:16]
		setSession(sessionID, QuerySession{
			Intent:    intent,
			ExpiresAt: time.Now().Add(10 * time.Minute),
		})

		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Message: "I need a bit more information to answer your question accurately.",
			Data: gin.H{
				"needs_clarification":    true,
				"your_question":          rawQuestion,
				"clarification_required": missing,
				"session_id":             sessionID,
				"hint":                   "Send your answers with the same session_id so I can remember the context.",
				"example": gin.H{
					"question":   "JAMB, Physics",
					"session_id": sessionID,
				},
			},
			Meta: &models.Meta{Version: "v1"},
		})
		return
	}

	// Clear session if used
	if req.SessionID != "" {
		deleteSession(req.SessionID)
	}

	// ── Stage 4: Orchestrate query ───────────────────────────────────────────
	result, err := orchestrate(intent)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: fmt.Sprintf("Could not retrieve curriculum: %v. Try rephrasing your question.", err),
		})
		return
	}

	// ── Stage 5: Cache the successful response ───────────────────────────────
	intentTags := []string{
		"board:" + intent.Board,
		"subject:" + intent.Subject,
		"action:" + intent.Action,
	}
	go cachePut(hash, rawQuestion, normalized, intentTags, result)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: fmt.Sprintf("Curriculum intelligence response for: %s", rawQuestion),
		Data: gin.H{
			"cache_hit": false,
			"query":     rawQuestion,
			"intent": gin.H{
				"board":   intent.Board,
				"subject": intent.Subject,
				"action":  intent.Action,
			},
			"result": result,
		},
		Meta: &models.Meta{Version: "v1"},
	})
}
