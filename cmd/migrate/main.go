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
		{"Create indexes", `
			CREATE INDEX IF NOT EXISTS idx_curricula_exam_board ON curricula(exam_board_id);
			CREATE INDEX IF NOT EXISTS idx_curricula_subject    ON curricula(subject_id);
			CREATE INDEX IF NOT EXISTS idx_topics_curriculum    ON topics(curriculum_id);
			CREATE INDEX IF NOT EXISTS idx_subtopics_topic      ON subtopics(topic_id);
			CREATE INDEX IF NOT EXISTS idx_objectives_subtopic  ON learning_objectives(subtopic_id)`},
		{"Seed exam boards", `
			INSERT INTO exam_boards (slug, name, full_name, country, description, website) VALUES
				('waec',  'WAEC',  'West African Examinations Council',                    'Nigeria', 'Conducts the WASSCE and other examinations across West Africa.', 'https://waec.org.ng'),
				('jamb',  'JAMB',  'Joint Admissions and Matriculation Board',             'Nigeria', 'Responsible for university entrance examinations in Nigeria.',   'https://jamb.gov.ng'),
				('neco',  'NECO',  'National Examinations Council',                        'Nigeria', 'Conducts the SSCE and other examinations in Nigeria.',           'https://neco.gov.ng'),
				('nerdc', 'NERDC', 'Nigerian Educational Research and Development Council','Nigeria', 'Responsible for developing the national curriculum.',            'https://nerdc.gov.ng')
			ON CONFLICT (slug) DO NOTHING`},
		{"Seed subjects", `
			INSERT INTO subjects (slug, name, description, category) VALUES
				('mathematics',          'Mathematics',          'The study of numbers, quantities, shapes, and patterns.',                  'science'),
				('english-language',     'English Language',     'The study of the English language, literature, and communication.',         'arts'),
				('physics',              'Physics',              'The study of matter, energy, and fundamental forces.',                      'science'),
				('chemistry',            'Chemistry',            'The study of substances, their properties, and reactions.',                 'science'),
				('biology',              'Biology',              'The study of living organisms and their interactions.',                     'science'),
				('economics',            'Economics',            'The study of how societies allocate scarce resources.',                     'commercial'),
				('government',           'Government',           'The study of political systems and civic responsibilities.',                'arts'),
				('literature',           'Literature in English','The study and analysis of literary texts in English.',                     'arts'),
				('geography',            'Geography',            'The study of physical features of the earth and human activity.',          'arts'),
				('further-mathematics',  'Further Mathematics',  'Advanced mathematics covering topics beyond the standard syllabus.',        'science')
			ON CONFLICT (slug) DO NOTHING`},
	}

	for _, stmt := range statements {
		log.Printf("  → %s...", stmt.name)
		if _, err := db.Exec(stmt.sql); err != nil {
			log.Fatalf("  ❌ Failed: %v", err)
		}
		log.Printf("  ✅ Done")
	}

	// Verify
	var boardCount, subjectCount int
	db.QueryRow("SELECT COUNT(*) FROM exam_boards").Scan(&boardCount)
	db.QueryRow("SELECT COUNT(*) FROM subjects").Scan(&subjectCount)

	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║      ✅ Schema deployed to Neon!         ║")
	fmt.Printf("║   Exam Boards: %-4d  Subjects: %-4d       ║\n", boardCount, subjectCount)
	fmt.Println("╠══════════════════════════════════════════╣")
	fmt.Println("║  Next: go run cmd/seeder/main.go         ║")
	fmt.Println("╚══════════════════════════════════════════╝")
}
