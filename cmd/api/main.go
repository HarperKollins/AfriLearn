package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/handlers"
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

	// ── Health check ─────────────────────────────────────────────────────────
	router.GET("/health", handlers.HealthCheck)

	// ── API v1 routes ─────────────────────────────────────────────────────────
	v1 := router.Group("/api/v1")
	{
		// Root info
		v1.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"service":     "AfriLearn Curriculum API",
				"version":     "v1",
				"description": "African Curriculum Infrastructure — BECE, WAEC, JAMB, NUC University Degrees, NBTE Polytechnics",
				"endpoints": gin.H{
					"subjects":   "GET /api/v1/subjects",
					"subject":    "GET /api/v1/subjects/:slug",
					"boards":     "GET /api/v1/exam-boards",
					"curriculum": "GET /api/v1/curriculum/:board/:subject",
					"llm_prompt": "GET /api/v1/curriculum/:board/:subject/llm-prompt",
					"search":     "GET /api/v1/search?q=<query>",
				},
				"example_calls": []string{
					"/api/v1/curriculum/waec/mathematics",
					"/api/v1/curriculum/waec/physics/llm-prompt",
					"/api/v1/curriculum/nuc/computer-science/llm-prompt",
					"/api/v1/curriculum/yabatech/computer-engineering-tech/llm-prompt",
					"/api/v1/search?q=quadratic equations",
					"/api/v1/subjects",
				},
			})
		})

		// Subjects
		v1.GET("/subjects", handlers.GetAllSubjects)
		v1.GET("/subjects/:slug", handlers.GetSubjectBySlug)

		// Exam Boards & Institutions
		v1.GET("/exam-boards", handlers.GetAllExamBoards)

		// Curriculum — main endpoints
		v1.GET("/curriculum/:board/:subject", handlers.GetCurriculum)
		v1.GET("/curriculum/:board/:subject/llm-prompt", handlers.GetLLMPrompt)

		// Search
		v1.GET("/search", handlers.SearchTopics)
	}

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("🚀 AfriLearn Curriculum API running on http://localhost:%s\n", port)
	log.Printf("📖 API Documentation: http://localhost:%s/api/v1/\n", port)
	log.Printf("💚 Health Check:      http://localhost:%s/health\n", port)
	log.Println("──────────────────────────────────────────────────────")
	log.Printf("📡 Try: curl http://localhost:%s/api/v1/curriculum/waec/physics/llm-prompt\n", port)

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
