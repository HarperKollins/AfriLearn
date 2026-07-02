package main

import (
	"log"

	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/scraper"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found")
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v\nMake sure PostgreSQL is running and .env is configured.", err)
	}
	defer database.Close()

	log.Println("╔════════════════════════════════════════════╗")
	log.Println("║   AfriLearn — WAEC Mathematics Seeder     ║")
	log.Println("╚════════════════════════════════════════════╝")
	log.Println()

	s := &scraper.WAECMathScraper{}
	if err := s.Run(); err != nil {
		log.Fatalf("❌ Scraper failed: %v", err)
	}

	log.Println()
	log.Println("✅ All done! Run the API server to query the data:")
	log.Println("   go run cmd/api/main.go")
	log.Println()
	log.Println("   Then try:")
	log.Println("   curl http://localhost:8080/api/v1/curriculum/waec/mathematics")
}
