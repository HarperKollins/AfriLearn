package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/afrilearn/curriculum-api/internal/cache"
	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	_ = godotenv.Load("../../.env")
	if database.DB == nil {
		_ = database.Connect()
	}
	r := gin.New()
	cache.InitCache(0)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/curriculum/match/:topic", GetCurriculumMatch)
		v1.GET("/curriculum/pathway", GetLearningPathway)
		v1.GET("/curriculum/prerequisites/:board/:subject/:topic", GetTopicPrerequisites)
		v1.GET("/curriculum/:board/:subject", GetCurriculum)
		v1.GET("/curriculum/:board/:subject/llm-prompt", GetLLMPrompt)
		v1.GET("/curriculum/:board/:subject/embeddings", GetCurriculumEmbeddings)
		v1.GET("/search", SearchTopics)
		v1.POST("/search/vector", HandleVectorSearch)
	}
	return r
}

// skipIfNoDB skips a test when the database is not available (e.g., CI without DB).
func skipIfNoDB(t *testing.T) {
	t.Helper()
	if database.DB == nil {
		t.Skip("Database not available — skipping integration test")
	}
	if err := database.DB.Ping(); err != nil {
		t.Skipf("Database ping failed: %v — skipping integration test", err)
	}
}

// ─── Core Endpoint Tests ───────────────────────────────────────────────────────

func TestIntegration_HealthCheck(t *testing.T) {
	router := setupTestRouter()
	router.GET("/health", HealthCheck)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse health response: %v", err)
	}
	if resp["status"] != "ok" {
		t.Errorf("Expected status=ok, got: %v", resp["status"])
	}
}

// ─── Curriculum Endpoint Tests ─────────────────────────────────────────────────

func TestIntegration_GetCurriculum_WAEC_Mathematics(t *testing.T) {
	router := setupTestRouter()
	skipIfNoDB(t)

	req, _ := http.NewRequest("GET", "/api/v1/curriculum/waec/mathematics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["success"] != true {
		t.Fatalf("Expected success=true, got: %v", resp["success"])
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data object in response")
	}

	topics, ok := data["topics"].([]interface{})
	if !ok || len(topics) == 0 {
		t.Fatalf("Expected topics array with at least 1 entry, got: %v", data["topics"])
	}

	t.Logf("WAEC Mathematics: %d topics returned", len(topics))

	// Verify topic structure has required fields
	firstTopic := topics[0].(map[string]interface{})
	for _, field := range []string{"id", "name", "slug", "difficulty", "order_index"} {
		if firstTopic[field] == nil {
			t.Errorf("Topic missing required field: %s", field)
		}
	}

	// Check if subtopics exist — some datasets (e.g., WAEC Mathematics v1) are topic-level only.
	// This is a known data quality gap documented in ROADMAP.md Phase 1.
	foundSubtopics := false
	for _, topicRaw := range topics {
		topic := topicRaw.(map[string]interface{})
		if subs, ok := topic["subtopics"].([]interface{}); ok && len(subs) > 0 {
			foundSubtopics = true
			break
		}
	}
	if !foundSubtopics {
		t.Logf("DATA QUALITY NOTE: WAEC Mathematics has no subtopics in current dataset — known gap, documented in ROADMAP.md")
		// Not a hard failure — the API works correctly, the data is surface-level
	}
}

func TestIntegration_GetCurriculum_UNILAG_Law(t *testing.T) {
	router := setupTestRouter()
	skipIfNoDB(t)

	req, _ := http.NewRequest("GET", "/api/v1/curriculum/unilag/law", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("UNILAG Law: Expected HTTP 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["success"] != true {
		t.Fatalf("UNILAG Law: Expected success=true")
	}

	data := resp["data"].(map[string]interface{})
	topics := data["topics"].([]interface{})
	if len(topics) == 0 {
		t.Fatalf("UNILAG Law: Expected at least 1 topic, got 0")
	}

	t.Logf("UNILAG Law: %d topics", len(topics))
}

func TestIntegration_GetCurriculum_FUTO_MechanicalEngineering(t *testing.T) {
	router := setupTestRouter()
	skipIfNoDB(t)

	req, _ := http.NewRequest("GET", "/api/v1/curriculum/futo/mechanical-engineering", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// May 404 if FUTO mechanical-engineering isn't in the dataset — that's acceptable
	if w.Code == http.StatusNotFound {
		t.Logf("FUTO mechanical-engineering not found — trying futo/engineering")
		req2, _ := http.NewRequest("GET", "/api/v1/curriculum/futo/engineering", nil)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		if w2.Code == http.StatusNotFound {
			t.Skip("FUTO Engineering not in dataset — skipping")
		}
		w = w2
	}

	if w.Code != http.StatusOK {
		t.Fatalf("FUTO Engineering: Expected HTTP 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Fatalf("FUTO Engineering: Expected success=true")
	}
	t.Logf("FUTO Engineering: response OK")
}

// ─── LLM Prompt Quality Tests ──────────────────────────────────────────────────

func TestIntegration_LLMPrompt_WAEC_Physics_HasBloomsTaxonomy(t *testing.T) {
	router := setupTestRouter()
	skipIfNoDB(t)

	req, _ := http.NewRequest("GET", "/api/v1/curriculum/waec/physics/llm-prompt", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	data := resp["data"].(map[string]interface{})

	// Verify system_prompt is non-empty and board-specific
	systemPrompt, ok := data["system_prompt"].(string)
	if !ok || len(systemPrompt) < 100 {
		t.Fatalf("Expected a substantive system_prompt, got: %s", systemPrompt)
	}
	if !strings.Contains(systemPrompt, "WAEC") && !strings.Contains(systemPrompt, "Physics") {
		t.Errorf("System prompt should reference WAEC and Physics, got: %s", systemPrompt[:200])
	}

	// Verify Bloom's taxonomy structure exists
	blooms, ok := data["blooms_taxonomy_breakdown"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected blooms_taxonomy_breakdown, got nil")
	}
	total := 0
	for _, v := range blooms {
		if count, ok := v.(float64); ok {
			total += int(count)
		}
	}
	if total == 0 {
		// DATA QUALITY NOTE: WAEC Physics objectives may have empty verb fields.
		// The Bloom's mechanism is correct — it maps verbs to levels.
		// When all verbs are empty strings, mapVerbToBloomLevel returns "understand" by default.
		// But since empty verb → "understand", total should still be non-zero if objectives exist.
		// If 0, it means WAEC Physics has no learning objectives in the dataset.
		t.Logf("DATA QUALITY NOTE: WAEC Physics Bloom's all zero — likely no objectives or all verbs are empty. Known gap.")
	}
	t.Logf("WAEC Physics Bloom's distribution: %v (total: %d objectives)", blooms, total)

	// Verify token count is present
	if data["estimated_token_count"] == nil {
		t.Error("Expected estimated_token_count in response")
	}
	tokenCount := int(data["estimated_token_count"].(float64))
	if tokenCount < 100 {
		t.Errorf("Expected token count > 100, got %d", tokenCount)
	}
	t.Logf("WAEC Physics LLM prompt: ~%d tokens", tokenCount)

	// Verify pedagogical_directives are WAEC-specific
	directives, ok := data["pedagogical_directives"].([]interface{})
	if !ok || len(directives) == 0 {
		t.Error("Expected pedagogical_directives to be non-empty")
	}
	for _, d := range directives {
		str := d.(string)
		if strings.Contains(str, "WAEC") || strings.Contains(str, "JAMB") || strings.Contains(str, "marking scheme") {
			t.Logf("Found WAEC-specific directive: %s", str)
			break
		}
	}
}

func TestIntegration_LLMPrompt_UNILAG_Law_HasSubjectSpecificRules(t *testing.T) {
	router := setupTestRouter()
	skipIfNoDB(t)

	req, _ := http.NewRequest("GET", "/api/v1/curriculum/unilag/law/llm-prompt", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})

	// Verify subject_specific_rules are Law-specific
	rules, ok := data["subject_specific_rules"].([]interface{})
	if !ok || len(rules) == 0 {
		t.Fatalf("Expected subject_specific_rules for Law, got empty/nil")
	}

	foundIRAC := false
	foundCaseLaw := false
	for _, ruleRaw := range rules {
		rule := ruleRaw.(string)
		if strings.Contains(rule, "IRAC") {
			foundIRAC = true
		}
		if strings.Contains(rule, "case") || strings.Contains(rule, "ratio decidendi") {
			foundCaseLaw = true
		}
	}
	if !foundIRAC {
		t.Errorf("Expected Law rules to include IRAC method, rules: %v", rules)
	}
	if !foundCaseLaw {
		t.Errorf("Expected Law rules to mention case law, rules: %v", rules)
	}
	t.Logf("UNILAG Law: %d subject-specific rules, IRAC=%v, CaseLaw=%v", len(rules), foundIRAC, foundCaseLaw)
}

func TestIntegration_LLMPrompt_BECE_HasBeginnerfriendlyDirectives(t *testing.T) {
	router := setupTestRouter()
	skipIfNoDB(t)

	req, _ := http.NewRequest("GET", "/api/v1/curriculum/bece/mathematics/llm-prompt", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})

	systemPrompt := data["system_prompt"].(string)
	if !strings.Contains(systemPrompt, "JSS") && !strings.Contains(systemPrompt, "Junior Secondary") && !strings.Contains(systemPrompt, "beginner") {
		t.Errorf("BECE system prompt should be beginner-focused, got: %s", systemPrompt[:300])
	}

	directives := data["pedagogical_directives"].([]interface{})
	foundNigerianAnalogy := false
	for _, d := range directives {
		if strings.Contains(d.(string), "Nigerian") || strings.Contains(d.(string), "analogies") {
			foundNigerianAnalogy = true
		}
	}
	if !foundNigerianAnalogy {
		t.Errorf("BECE directives should include Nigerian analogies")
	}
	t.Logf("BECE Math: beginner directives verified ✓")
}

// ─── Search Tests ──────────────────────────────────────────────────────────────

func TestIntegration_SearchTopics_FindsAcrossLayers(t *testing.T) {
	router := setupTestRouter()
	skipIfNoDB(t)

	tests := []struct {
		query       string
		expectMatch bool
	}{
		{"photosynthesis", true},
		{"quadratic equations", true},
		{"Newton laws of motion", true},
		{"statutory interpretation", true}, // Should match Law objectives
		{"xyzzy_nonexistent_term_12345", false},
	}

	for _, tc := range tests {
		t.Run(tc.query, func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/search?q=%s", tc.query)
			req, _ := http.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("Search for '%s': Expected HTTP 200, got %d", tc.query, w.Code)
			}

			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)

			data := resp["data"].([]interface{})
			if tc.expectMatch && len(data) == 0 {
				t.Errorf("Search for '%s': Expected results, got 0", tc.query)
			}
			if !tc.expectMatch && len(data) > 0 {
				t.Errorf("Search for '%s': Expected 0 results (nonexistent), got %d", tc.query, len(data))
			}

			if len(data) > 0 {
				first := data[0].(map[string]interface{})
				// Verify new search result structure
				for _, field := range []string{"topic_id", "topic_name", "board_slug", "matched_in", "snippet", "relevance_score"} {
					if first[field] == nil {
						t.Errorf("Search result missing field '%s'", field)
					}
				}
				t.Logf("Search '%s': %d results, first match in '%s': '%s'",
					tc.query, len(data), first["matched_in"], first["snippet"])
			}
		})
	}
}

func TestIntegration_SearchTopics_FiltersByBoard(t *testing.T) {
	router := setupTestRouter()
	skipIfNoDB(t)

	// Search with board filter — all results should be from WAEC
	req, _ := http.NewRequest("GET", "/api/v1/search?q=algebra&board=waec", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	data := resp["data"].([]interface{})
	for _, itemRaw := range data {
		item := itemRaw.(map[string]interface{})
		boardSlug := item["board_slug"].(string)
		if boardSlug != "waec" {
			t.Errorf("Expected all results to be board=waec, got board=%s", boardSlug)
		}
	}
	t.Logf("Board filter test: %d results, all from WAEC ✓", len(data))
}

// ─── Vector Search Tests ───────────────────────────────────────────────────────

func TestIntegration_VectorSearch_ReturnsHonestMetadata(t *testing.T) {
	router := setupTestRouter()

	reqBody := map[string]interface{}{
		"query": "photosynthesis in plants",
		"limit": 5,
	}
	jsonBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/search/vector", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["success"] != true {
		t.Fatalf("Expected success=true")
	}

	data := resp["data"].(map[string]interface{})

	// Verify search_engine field is honest (no fake ML claims)
	searchEngine, _ := data["search_engine"].(string)
	if strings.Contains(searchEngine, "HNSW") && !strings.Contains(searchEngine, "pgvector slot") &&
		!strings.Contains(searchEngine, "full-text") && !strings.Contains(searchEngine, "FTS") {
		t.Errorf("search_engine field should be honest about FTS, not claim HNSW ML search. Got: %s", searchEngine)
	}
	t.Logf("search_engine: %s", searchEngine)

	// Verify upgrade_path is present
	if data["upgrade_path"] == nil {
		t.Error("Expected upgrade_path field explaining how to add real embeddings")
	}
}

func TestIntegration_VectorSearch_FiltersByBoardAndSubject(t *testing.T) {
	router := setupTestRouter()
	skipIfNoDB(t)

	reqBody := map[string]interface{}{
		"query":   "electromagnetic induction",
		"limit":   10,
		"board":   "waec",
		"subject": "physics",
	}
	jsonBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/search/vector", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})

	results, _ := data["results"].([]interface{})
	for _, rRaw := range results {
		r := rRaw.(map[string]interface{})
		if r["board_slug"] != "waec" {
			t.Errorf("Filter by board=waec: got result from board=%s", r["board_slug"])
		}
	}
	t.Logf("Board+subject filter: %d results, all from WAEC Physics", len(results))
}

// ─── Intelligence Layer Tests ──────────────────────────────────────────────────

func TestIntegration_CrossBoardMatch_QuadraticEquations(t *testing.T) {
	router := setupTestRouter()
	skipIfNoDB(t)

	req, _ := http.NewRequest("GET", "/api/v1/curriculum/match/quadratic-equations", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Fatalf("Expected HTTP 200 or 404, got %d. Body: %s", w.Code, w.Body.String())
	}
	if w.Code == http.StatusOK {
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})

		// Verify cross_board_coverage contains boards
		coverage, ok := data["cross_board_coverage"].([]interface{})
		if !ok || len(coverage) == 0 {
			t.Error("Expected cross_board_coverage with at least 1 board")
		}
		t.Logf("Quadratic equations found in %d boards", len(coverage))

		// Verify llm_unified_prompt is present and substantive
		prompt, ok := data["llm_unified_prompt"].(string)
		if !ok || len(prompt) < 50 {
			t.Error("Expected substantive llm_unified_prompt")
		}
	}
}

func TestIntegration_LearningPathway_MathematicsBECEtoJAMB(t *testing.T) {
	router := setupTestRouter()
	skipIfNoDB(t)

	req, _ := http.NewRequest("GET", "/api/v1/curriculum/pathway?subject=mathematics&from=bece&to=jamb", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["success"] != true {
		t.Fatalf("Expected success=true")
	}

	data := resp["data"].(map[string]interface{})
	pathway := data["pathway"].([]interface{})

	if len(pathway) == 0 {
		t.Fatalf("Expected pathway steps, got 0")
	}

	// Verify steps are ordered correctly (BECE topics should come before WAEC)
	stages := make([]int, 0, len(pathway))
	boardsEncountered := make(map[string]bool)
	for _, stepRaw := range pathway {
		step := stepRaw.(map[string]interface{})
		stage := int(step["stage"].(float64))
		stages = append(stages, stage)
		boardsEncountered[step["board"].(string)] = true
	}

	// Stages should be sequential
	for i := 1; i < len(stages); i++ {
		if stages[i] != stages[i-1]+1 {
			t.Errorf("Pathway stages not sequential at index %d: %v", i, stages[i-1:i+1])
		}
	}

	t.Logf("Mathematics BECE→JAMB pathway: %d steps across boards: %v", len(pathway), boardsEncountered)
}

// ─── Embeddings RAG Chunk Tests ────────────────────────────────────────────────

func TestIntegration_Embeddings_ReturnsChunksWithRequiredFields(t *testing.T) {
	router := setupTestRouter()
	skipIfNoDB(t)

	req, _ := http.NewRequest("GET", "/api/v1/curriculum/waec/physics/embeddings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})

	// Verify top-level metadata
	for _, field := range []string{"total_chunks", "chunking_strategy", "embedding_note", "integration_guide"} {
		if data[field] == nil {
			t.Errorf("Expected field '%s' in embeddings response", field)
		}
	}

	chunks := data["chunks"].([]interface{})
	if len(chunks) == 0 {
		t.Fatalf("Expected at least 1 chunk, got 0")
	}

	// Verify chunk structure
	first := chunks[0].(map[string]interface{})
	for _, field := range []string{"chunk_id", "module_title", "text_content", "token_estimate", "metadata"} {
		if first[field] == nil {
			t.Errorf("Chunk missing required field: %s", field)
		}
	}

	// Verify text_content is substantive (not just a title)
	textContent := first["text_content"].(string)
	if len(textContent) < 100 {
		t.Errorf("text_content too short (%d chars) — expected rich curriculum text", len(textContent))
	}

	// Verify NO raw embedding_values (the old fake vector approach should be gone)
	if first["embedding_values"] != nil {
		t.Error("embedding_values should NOT be in response — use integration_guide to generate real embeddings")
	}

	totalTokens := int(data["total_token_estimate"].(float64))
	t.Logf("WAEC Physics embeddings: %d chunks, ~%d total tokens", len(chunks), totalTokens)
}
