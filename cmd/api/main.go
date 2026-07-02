package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/handlers"
	"github.com/afrilearn/curriculum-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found — using system environment variables")
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Printf("⚠️  Database not connected: %v\n", err)
		log.Println("   Running without database — only static endpoints will work.")
		log.Println("   Set up PostgreSQL and add DB_* vars to .env to enable full functionality.")
	}
	defer database.Close()

	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// ── Public Web Routes (Root, Health, Portal Dashboard, Swagger Docs) ──────
	router.GET("/", handlers.ServeDeveloperPortal)
	router.GET("/health", handlers.HealthCheck)
	router.GET("/portal", handlers.ServeDeveloperPortal)
	router.GET("/docs", handlers.ServeSwaggerUI)
	router.GET("/swagger", handlers.ServeSwaggerUI)
	router.GET("/docs/openapi.json", handlers.ServeOpenAPISpec)

	// ── API v1 routes ─────────────────────────────────────────────────────────
	v1 := router.Group("/api/v1")
	{
		// Unprotected API Key generation for self-service portal
		v1.POST("/keys/generate", handlers.GenerateAPIKey)
	}

	// Protected API v1 endpoints (Require/validate X-API-Key)
	v1Protected := router.Group("/api/v1")
	v1Protected.Use(middleware.APIKeyAuth())
	{
		// Root info
		v1Protected.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"service":     "AfriLearn Curriculum API",
				"version":     "v1",
				"description": "African Curriculum Infrastructure — BECE, WAEC, JAMB, NUC University Degrees, NBTE Polytechnics",
				"portal_url":  "/portal",
				"docs_url":    "/docs",
				"authentication": gin.H{
					"header":      "X-API-Key",
					"query_param": "api_key",
					"demo_key":    "afr_live_demo_9f8e2b7a",
				},
				"endpoints": gin.H{
					"subjects":     "GET /api/v1/subjects",
					"subject":      "GET /api/v1/subjects/:slug",
					"boards":       "GET /api/v1/exam-boards",
					"curriculum":   "GET /api/v1/curriculum/:board/:subject",
					"llm_prompt":   "GET /api/v1/curriculum/:board/:subject/llm-prompt",
					"search":       "GET /api/v1/search?q=<query>",
					"generate_key": "POST /api/v1/keys/generate",
				},
				"example_calls": []string{
					"/api/v1/curriculum/waec/mathematics",
					"/api/v1/curriculum/waec/physics/llm-prompt",
					"/api/v1/curriculum/nuc/computer-science/llm-prompt",
					"/api/v1/curriculum/yabatech/computer-engineering-tech/llm-prompt",
					"/api/v1/search?q=quadratic equations",
				},
			})
		})

		// Subjects
		v1Protected.GET("/subjects", handlers.GetAllSubjects)
		v1Protected.GET("/subjects/:slug", handlers.GetSubjectBySlug)

		// Exam Boards & Institutions
		v1Protected.GET("/exam-boards", handlers.GetAllExamBoards)

		// Curriculum — main endpoints
		v1Protected.GET("/curriculum/:board/:subject", handlers.GetCurriculum)
		v1Protected.GET("/curriculum/:board/:subject/llm-prompt", handlers.GetLLMPrompt)

		// Search
		v1Protected.GET("/search", handlers.SearchTopics)
	}

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("🚀 AfriLearn Curriculum API running on http://localhost:%s\n", port)
	log.Printf("🔑 Developer Portal:          http://localhost:%s/\n", port)
	log.Printf("📖 Interactive Swagger Docs: http://localhost:%s/docs\n", port)
	log.Printf("💚 Health Check:              http://localhost:%s/health\n", port)
	log.Println("──────────────────────────────────────────────────────")
	log.Printf("📡 Try: curl -H 'X-API-Key: afr_live_demo_9f8e2b7a' http://localhost:%s/api/v1/curriculum/waec/physics/llm-prompt\n", port)

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("\n🔴 Shutting down server...")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
