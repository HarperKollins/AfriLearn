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

		// NUC CCMAS (University National Core Standards - 100L to 500L)
		scraper.NewNUCComputerScienceScraper(),
		scraper.NewNUCMedicineScraper(),
		scraper.NewNUCElectricalEngScraper(),
		scraper.NewNUCLawScraper(),
		scraper.NewNUCAccountingScraper(),
		scraper.NewNUCBusinessAdminScraper(),
		scraper.NewNUCNursingScraper(),
		scraper.NewNUCMechanicalEngScraper(),
		scraper.NewNUCMassCommScraper(),

		// Individual Universities (EBSU, AE-FUNAI, UNEC, UNN, UNILAG, FUTO)
		scraper.NewEBSUComputerScienceScraper(),
		scraper.NewFUNAIComputerScienceScraper(),
		scraper.NewUNECLawScraper(),
		scraper.NewUNNComputerScienceScraper(),
		scraper.NewUNILAGComputerScienceScraper(),
		scraper.NewFUTOPetroleumEngScraper(),

		// NBTE Polytechnics (YABATECH, IMT Enugu)
		scraper.NewYABATECHComputerEngScraper(),
		scraper.NewIMTSLTScraper(),
	}

	for i, s := range scrapers {
		log.Printf("------------------------------------------------------------ [%d/%d]", i+1, len(scrapers))
		if err := engine.Execute(s); err != nil {
			log.Printf("❌ Ingestion error for [%s/%s]: %v\n", s.BoardSlug(), s.SubjectSlug(), err)
		}
	}

	log.Println("------------------------------------------------------------")
	log.Println("✅ All 36 Curricula Ingested Successfully!")
	log.Println()
	log.Println("   BECE (JSS1-3):")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/bece/mathematics")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/bece/basic-science")
	log.Println()
	log.Println("   WAEC (SS1-3):")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/waec/mathematics")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/waec/physics")
	log.Println()
	log.Println("   JAMB (UTME):")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/jamb/mathematics")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/jamb/physics")
	log.Println()
	log.Println("   NUC CCMAS National Degrees (100L - 500L):")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/nuc/computer-science")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/nuc/medicine-and-surgery")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/nuc/electrical-engineering")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/nuc/law")
	log.Println()
	log.Println("   Individual Universities (EBSU, FUNAI, UNEC, UNN, UNILAG, FUTO):")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/ebsu/computer-science")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/funai/computer-science")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/unec/law")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/unn/computer-science")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/unilag/computer-science")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/futo/petroleum-engineering")
	log.Println()
	log.Println("   NBTE Polytechnics (YABATECH, IMT Enugu):")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/yabatech/computer-engineering-tech")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/imt/science-laboratory-tech")
}
