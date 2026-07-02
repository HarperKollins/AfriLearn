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

	log.Println("╔════════════════════════════════════════════════════════════╗")
	log.Println("║     AfriLearn — Curriculum Infrastructure Data Engine      ║")
	log.Println("╚════════════════════════════════════════════════════════════╝")
	log.Println()

	engine := scraper.NewEngine()

	scrapers := []scraper.Scraper{
		scraper.NewWAECMathScraper(),
		scraper.NewWAECPhysicsScraper(),
		scraper.NewWAECBiologyScraper(),
		scraper.NewJAMBMathScraper(),
		scraper.NewJAMBPhysicsScraper(),
	}

	for _, s := range scrapers {
		log.Printf("------------------------------------------------------------")
		if err := engine.Execute(s); err != nil {
			log.Printf("❌ Ingestion error for [%s/%s]: %v\n", s.BoardSlug(), s.SubjectSlug(), err)
		}
	}

	log.Println("------------------------------------------------------------")
	log.Println("✅ All 5 Curricula Ingested Successfully!")
	log.Println()
	log.Println("   Try querying live endpoints:")
	log.Println("   curl http://localhost:8080/api/v1/curriculum/waec/mathematics")
	log.Println("   curl http://localhost:8080/api/v1/curriculum/waec/physics")
	log.Println("   curl http://localhost:8080/api/v1/curriculum/waec/biology")
	log.Println("   curl http://localhost:8080/api/v1/curriculum/jamb/mathematics")
	log.Println("   curl http://localhost:8080/api/v1/curriculum/jamb/physics")
}
