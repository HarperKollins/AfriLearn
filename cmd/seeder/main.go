package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/ingestion"
	"github.com/joho/godotenv"
)

func main() {
	validateOnly := flag.Bool("validate-only", false, "Scan and validate all JSON files without writing to database")
	flag.Parse()

	// Resolve the data/curricula directory relative to where the binary runs.
	dataDir := "data/curricula"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		dataDir = filepath.Join("..", "..", "data", "curricula")
	}

	engine := ingestion.NewEngine(dataDir)

	if *validateOnly {
		log.Println("🔍 Running AfriLearn Curriculum Data Quality & Schema Validation...")
		count, valErrs, err := engine.ValidateAll()
		if err != nil {
			log.Fatalf("❌ Validation error: %v", err)
		}
		log.Printf("📊 Scanned %d curriculum JSON files.", count)
		if len(valErrs) == 0 {
			log.Println("✨ PERFECT! All curriculum files passed 100% schema & quality checks!")
			return
		}

		log.Printf("⚠️  Found %d validation warning(s)/error(s):", len(valErrs))
		for i, ve := range valErrs {
			log.Printf("  [%d] File: %s | Field: %s -> %s", i+1, ve.FilePath, ve.Field, ve.Message)
		}
		os.Exit(1)
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v\nMake sure DATABASE_URL is set in .env", err)
	}
	defer database.Close()

	log.Println("╔════════════════════════════════════════════════════════════╗")
	log.Println("║     AfriLearn — Curriculum Infrastructure Data Engine      ║")
	log.Println("║     Phase 1: JSON-driven ingestion (no hardcoded Go)       ║")
	log.Println("╚════════════════════════════════════════════════════════════╝")
	log.Println()

	if err := engine.Run(); err != nil {
		log.Fatalf("❌ Ingestion failed: %v", err)
	}

	log.Println()
	log.Println("────────────────────────────────────────────────────────────")
	log.Println("🎉 All curricula ingested! Try a few endpoints:")
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
	log.Println("   NUC CCMAS National Degrees:")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/nuc/computer-science")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/nuc/medicine-and-surgery")
	log.Println()
	log.Println("   Cross-board match:")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/match/quadratic-equations")
	log.Println("     curl http://localhost:8080/api/v1/curriculum/match/photosynthesis")
}
