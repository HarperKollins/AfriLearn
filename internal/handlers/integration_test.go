package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		v1.POST("/search/vector", HandleVectorSearch)
	}
	return r
}

func TestIntegration_VectorSearchEndpoint(t *testing.T) {
	router := setupTestRouter()

	reqBody := map[string]interface{}{
		"query": "photosynthesis in plants",
		"limit": 3,
	}
	jsonBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/search/vector", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusServiceUnavailable {
		t.Fatalf("Expected HTTP Status 200 or 530, got %d. Body: %s", w.Code, w.Body.String())
	}

	if w.Code == http.StatusOK {
		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}
		if resp["success"] != true {
			t.Fatalf("Expected success = true in response")
		}
	}
}

func TestIntegration_LearningPathwayEndpoint(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/curriculum/pathway?subject=mathematics&from=bece&to=jamb", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusServiceUnavailable {
		t.Fatalf("Expected HTTP Status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}
