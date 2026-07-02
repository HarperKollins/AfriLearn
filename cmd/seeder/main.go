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
		// BECE (Junior Secondary JSS1 - JSS3)
		scraper.NewBECEMathScraper(),
		scraper.NewBECEBasicScienceScraper(),
		scraper.NewBECEBasicTechScraper(),
		scraper.NewBECEEnglishScraper(),
		scraper.NewBECESocialStudiesScraper(),
		scraper.NewBECEBusinessStudiesScraper(),

		// WAEC (Senior Secondary SS1 - SS3)
		scraper.NewWAECMathScraper(),
		scraper.NewWAECPhysicsScraper(),
		scraper.NewWAECBiologyScraper(),
		scraper.NewWAECChemistryScraper(),
		scraper.NewWAECEconomicsScraper(),
		scraper.NewWAECGovernmentScraper(),
		scraper.NewWAECLiteratureScraper(),

		// JAMB (UTME Tertiary Entry)
		scraper.NewJAMBMathScraper(),
		scraper.NewJAMBPhysicsScraper(),
		scraper.NewJAMBChemistryScraper(),
		scraper.NewJAMBBiologyScraper(),
		scraper.NewJAMBEconomicsScraper(),
		scraper.NewJAMBGovernmentScraper(),
	}

	for i, s := range scrapers {
		log.Printf("------------------------------------------------------------ [%d/%d]", i+1, len(scrapers))
		if err := engine.Execute(s); err != nil {
			log.Printf("❌ Ingestion error for [%s/%s]: %v\n", s.BoardSlug(), s.SubjectSlug(), err)
		}
	}

	log.Println("------------------------------------------------------------")
	log.Println("✅ All 19 Curricula Ingested Successfully!")
	log.Println()
	log.Println("   BECE (JSS1-3):")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/bece/mathematics")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/bece/basic-science")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/bece/basic-technology")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/bece/english-language")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/bece/social-studies")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/bece/business-studies")
	log.Println()
	log.Println("   WAEC (SS1-3):")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/waec/mathematics")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/waec/physics")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/waec/biology")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/waec/chemistry")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/waec/economics")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/waec/government")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/waec/literature")
	log.Println()
	log.Println("   JAMB (UTME):")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/jamb/mathematics")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/jamb/physics")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/jamb/chemistry")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/jamb/biology")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/jamb/economics")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/jamb/government")
}
