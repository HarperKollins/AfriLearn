package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL not set in .env")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Cannot connect to Neon: %v", err)
	}
	log.Println("✅ Connected to Neon PostgreSQL!")

	// Run schema statements one by one
	statements := []struct {
		name string
		sql  string
	}{
		{"Enable UUID extension", `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`},
		{"Create exam_boards table", `
			CREATE TABLE IF NOT EXISTS exam_boards (
				id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				slug        VARCHAR(100) NOT NULL UNIQUE,
				name        VARCHAR(100) NOT NULL,
				full_name   VARCHAR(255) NOT NULL,
				country     VARCHAR(100) NOT NULL DEFAULT 'Nigeria',
				description TEXT,
				website     VARCHAR(255),
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			)`},
		{"Create subjects table", `
			CREATE TABLE IF NOT EXISTS subjects (
				id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				slug        VARCHAR(100) NOT NULL UNIQUE,
				name        VARCHAR(255) NOT NULL,
				description TEXT,
				category    VARCHAR(50) NOT NULL DEFAULT 'general',
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			)`},
		{"Create curricula table", `
			CREATE TABLE IF NOT EXISTS curricula (
				id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				exam_board_id   UUID NOT NULL REFERENCES exam_boards(id) ON DELETE CASCADE,
				subject_id      UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
				year            INTEGER NOT NULL DEFAULT EXTRACT(YEAR FROM NOW()),
				level           VARCHAR(50) NOT NULL DEFAULT 'senior-secondary',
				source_url      VARCHAR(500),
				created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				UNIQUE(exam_board_id, subject_id, year)
			)`},
		{"Create topics table", `
			CREATE TABLE IF NOT EXISTS topics (
				id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				curriculum_id   UUID NOT NULL REFERENCES curricula(id) ON DELETE CASCADE,
				slug            VARCHAR(200) NOT NULL,
				name            VARCHAR(255) NOT NULL,
				description     TEXT,
				order_index     INTEGER NOT NULL DEFAULT 0,
				difficulty      VARCHAR(20) NOT NULL DEFAULT 'medium',
				created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				UNIQUE(curriculum_id, slug)
			)`},
		{"Create subtopics table", `
			CREATE TABLE IF NOT EXISTS subtopics (
				id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				topic_id    UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
				slug        VARCHAR(200) NOT NULL,
				name        VARCHAR(255) NOT NULL,
				description TEXT,
				order_index INTEGER NOT NULL DEFAULT 0,
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				UNIQUE(topic_id, slug)
			)`},
		{"Create learning_objectives table", `
			CREATE TABLE IF NOT EXISTS learning_objectives (
				id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				subtopic_id UUID NOT NULL REFERENCES subtopics(id) ON DELETE CASCADE,
				description TEXT NOT NULL,
				verb        VARCHAR(50) NOT NULL DEFAULT 'understand',
				order_index INTEGER NOT NULL DEFAULT 0,
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			)`},
		{"Create api_keys table", `
			CREATE TABLE IF NOT EXISTS api_keys (
				id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				api_key        VARCHAR(100) NOT NULL UNIQUE,
				developer_name VARCHAR(255) NOT NULL,
				email          VARCHAR(255) NOT NULL,
				tier           VARCHAR(50) NOT NULL DEFAULT 'free',
				requests_count BIGINT NOT NULL DEFAULT 0,
				is_active      BOOLEAN NOT NULL DEFAULT true,
				created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
			)`},
		{"Create indexes", `
			CREATE INDEX IF NOT EXISTS idx_curricula_exam_board ON curricula(exam_board_id);
			CREATE INDEX IF NOT EXISTS idx_curricula_subject    ON curricula(subject_id);
			CREATE INDEX IF NOT EXISTS idx_topics_curriculum    ON topics(curriculum_id);
			CREATE INDEX IF NOT EXISTS idx_subtopics_topic      ON subtopics(topic_id);
			CREATE INDEX IF NOT EXISTS idx_objectives_subtopic  ON learning_objectives(subtopic_id);
			CREATE INDEX IF NOT EXISTS idx_api_keys_lookup      ON api_keys(api_key)`},
		{"Seed API keys", `
			INSERT INTO api_keys (api_key, developer_name, email, tier) VALUES
				('afr_live_demo_9f8e2b7a', 'AfriLearn Demo Developer', 'demo@afrilearn.org', 'free'),
				('afr_live_pro_8372bf91',  'AfriLearn Pro Partner',    'pro@afrilearn.org',  'pro')
			ON CONFLICT (api_key) DO NOTHING`},
	}

	for _, stmt := range statements {
		log.Printf("  → %s...", stmt.name)
		if _, err := db.Exec(stmt.sql); err != nil {
			log.Fatalf("  ❌ Failed: %v", err)
		}
		log.Printf("  ✅ Done")
	}

	// Verify
	var boardCount, subjectCount, keyCount int
	db.QueryRow("SELECT COUNT(*) FROM exam_boards").Scan(&boardCount)
	db.QueryRow("SELECT COUNT(*) FROM subjects").Scan(&subjectCount)
	db.QueryRow("SELECT COUNT(*) FROM api_keys").Scan(&keyCount)

	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║      ✅ Schema deployed to Neon!         ║")
	fmt.Printf("║  Institutions: %-3d  Subjects/Degrees: %-3d  ║\n", boardCount, subjectCount)
	fmt.Printf("║  API Keys Registered: %-18d ║\n", keyCount)
	fmt.Println("╚══════════════════════════════════════════╝")
}
