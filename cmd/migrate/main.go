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
		{"Enable Vector extension", `CREATE EXTENSION IF NOT EXISTS vector`},
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
				embedding       vector(1536),
				created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				UNIQUE(curriculum_id, slug)
			)`},
		{"Add embedding to topics if missing", `ALTER TABLE topics ADD COLUMN IF NOT EXISTS embedding vector(1536)`},
		{"Create subtopics table", `
			CREATE TABLE IF NOT EXISTS subtopics (
				id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				topic_id     UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
				slug         VARCHAR(200) NOT NULL,
				name         VARCHAR(255) NOT NULL,
				description  TEXT,
				course_code  VARCHAR(50),
				credit_units VARCHAR(20),
				semester     VARCHAR(50),
				order_index  INTEGER NOT NULL DEFAULT 0,
				embedding    vector(1536),
				created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				UNIQUE(topic_id, slug)
			)`},
		{"Add course_code to subtopics", `ALTER TABLE subtopics ADD COLUMN IF NOT EXISTS course_code VARCHAR(50)`},
		{"Add credit_units to subtopics", `ALTER TABLE subtopics ADD COLUMN IF NOT EXISTS credit_units VARCHAR(20)`},
		{"Add semester to subtopics", `ALTER TABLE subtopics ADD COLUMN IF NOT EXISTS semester VARCHAR(50)`},
		{"Add embedding to subtopics", `ALTER TABLE subtopics ADD COLUMN IF NOT EXISTS embedding vector(1536)`},
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
		{"Create topic_prerequisites table", `
			CREATE TABLE IF NOT EXISTS topic_prerequisites (
				id                   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				topic_id             UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
				prerequisite_topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
				order_index          INTEGER NOT NULL DEFAULT 0,
				created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				UNIQUE(topic_id, prerequisite_topic_id)
			)`},
		{"Create query_cache table", `
			CREATE TABLE IF NOT EXISTS query_cache (
				id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				query_hash    TEXT NOT NULL UNIQUE,
				raw_query     TEXT NOT NULL,
				normalized    TEXT NOT NULL,
				intent_tags   TEXT[],
				response_json JSONB NOT NULL,
				hit_count     INTEGER NOT NULL DEFAULT 1,
				last_hit_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
			)`},
		{"Create intelligence layer indexes", `
			CREATE INDEX IF NOT EXISTS idx_topic_prereqs_topic      ON topic_prerequisites(topic_id);
			CREATE INDEX IF NOT EXISTS idx_topic_prereqs_prereq     ON topic_prerequisites(prerequisite_topic_id);
			CREATE INDEX IF NOT EXISTS idx_query_cache_hash         ON query_cache(query_hash);
			CREATE INDEX IF NOT EXISTS idx_query_cache_hits         ON query_cache(hit_count DESC);
			CREATE INDEX IF NOT EXISTS idx_topics_name_search       ON topics USING gin(to_tsvector('english', name));
			CREATE INDEX IF NOT EXISTS idx_subtopics_name_search    ON subtopics USING gin(to_tsvector('english', name));
			CREATE INDEX IF NOT EXISTS idx_objectives_subtopic_desc ON learning_objectives(subtopic_id, description)`},

		// ── Seed Reference Data ─────────────────────────────────────────────────
		{"Seed exam_boards", `
			INSERT INTO exam_boards (slug, name, full_name, country, description, website) VALUES
				('bece',     'BECE',     'Basic Education Certificate Examination',                         'Nigeria', 'Junior Secondary School (JSS1-JSS3) national examinations board.',          'https://waec.org.ng'),
				('waec',     'WAEC',     'West African Examinations Council',                               'Nigeria', 'Senior Secondary Certificate Examination (SSCE) for SS1-SS3.',               'https://waec.org.ng'),
				('jamb',     'JAMB',     'Joint Admissions and Matriculation Board',                        'Nigeria', 'Unified Tertiary Matriculation Examination (UTME) for university entry.',    'https://jamb.gov.ng'),
				('neco',     'NECO',     'National Examinations Council',                                   'Nigeria', 'National SSCE alternative to WAEC for Senior Secondary students.',           'https://neco.gov.ng'),
				('nuc',      'NUC',      'National Universities Commission',                                'Nigeria', 'CCMAS minimum academic standards for Nigerian university degree programmes.', 'https://nuc.edu.ng'),
				('nbte',     'NBTE',     'National Board for Technical Education',                          'Nigeria', 'Polytechnic ND/HND minimum standards for technical and vocational education.','https://nbte.gov.ng'),
				('yabatech', 'YABATECH', 'Yaba College of Technology',                                      'Nigeria', 'Premier polytechnic in Lagos offering ND and HND programmes.',              'https://yabatech.edu.ng'),
				('imt',      'IMT',      'Institute of Management and Technology Enugu',                    'Nigeria', 'State polytechnic in Enugu offering ND and HND programmes.',                'https://imt.edu.ng'),
				('unilag',   'UNILAG',   'University of Lagos',                                             'Nigeria', 'Federal university in Lagos, one of Nigeria''s foremost research universities.','https://unilag.edu.ng'),
				('unn',      'UNN',      'University of Nigeria Nsukka',                                    'Nigeria', 'Premier South-Eastern federal university founded by Nnamdi Azikiwe.',        'https://unn.edu.ng'),
				('unec',     'UNEC',     'University of Nigeria Enugu Campus',                              'Nigeria', 'Enugu campus of University of Nigeria, specialising in law and medical sciences.','https://unn.edu.ng'),
				('ebsu',     'EBSU',     'Ebonyi State University',                                         'Nigeria', 'State university in Abakaliki, Ebonyi State.',                              'https://ebsu.edu.ng'),
				('funai',    'FUNAI',    'Federal University Ndufu-Alike Ikwo',                             'Nigeria', 'Federal university in Ebonyi State (AE-FUNAI).',                            'https://funai.edu.ng'),
				('futo',     'FUTO',     'Federal University of Technology Owerri',                         'Nigeria', 'Technology-focused federal university in Imo State.',                       'https://futo.edu.ng')
			ON CONFLICT (slug) DO NOTHING`},

		{"Seed subjects", `
			INSERT INTO subjects (slug, name, description, category) VALUES
				-- Secondary & UTME subjects
				('mathematics',          'Mathematics',                    'Pure and applied mathematics from arithmetic to calculus.',                    'science'),
				('physics',              'Physics',                        'Study of matter, energy, mechanics, waves, and modern physics.',               'science'),
				('chemistry',            'Chemistry',                      'Study of elements, compounds, reactions, and organic chemistry.',               'science'),
				('biology',              'Biology',                        'Study of living organisms, ecology, genetics, and physiology.',                 'science'),
				('economics',            'Economics',                      'Microeconomics, macroeconomics, market theory, and national development.',       'social-science'),
				('government',           'Government',                     'Nigerian government, constitution, democracy, and political science.',           'social-science'),
				('english-language',     'English Language',               'Reading, writing, comprehension, and grammar for JSS and primary levels.',      'arts'),
				('english-studies',      'English Studies',                'Advanced English language, literature, and communication skills.',               'arts'),
				('literature',           'Literature in English',          'Prose, drama, poetry, and African literary works for secondary students.',       'arts'),
				('social-studies',       'Social Studies',                 'Citizenship, society, environment, and civic education for JSS.',                'social-science'),
				('basic-science',        'Basic Science',                  'Integrated science introducing physics, chemistry, and biology at JSS level.',   'science'),
				('basic-technology',     'Basic Technology',               'Technology fundamentals, workshop practice, and technical drawing at JSS.',      'technology'),
				('business-studies',     'Business Studies',               'Commerce, bookkeeping, and office practice at junior secondary level.',          'business'),
				-- University degree subjects (NUC CCMAS)
				('computer-science',         'Computer Science',                'Algorithms, data structures, software engineering, AI, and networking.',       'technology'),
				('law',                      'Law',                             'Nigerian and international law, jurisprudence, constitutional law.',           'law'),
				('accounting',               'Accounting',                      'Financial accounting, auditing, taxation, and management accounting.',         'business'),
				('business-administration',  'Business Administration',         'Management, marketing, entrepreneurship, and organisational behaviour.',       'business'),
				('nursing-science',          'Nursing Science',                 'Clinical nursing, anatomy, pharmacology, and community health care.',          'health'),
				('medicine-and-surgery',     'Medicine and Surgery',            'MBBS programme: anatomy, physiology, pathology, clinical rotations.',          'health'),
				('mechanical-engineering',   'Mechanical Engineering',          'Thermodynamics, fluid mechanics, manufacturing, and machine design.',          'engineering'),
				('electrical-engineering',   'Electrical Engineering',          'Circuit theory, power systems, electronics, and control engineering.',         'engineering'),
				('petroleum-engineering',    'Petroleum Engineering',           'Reservoir engineering, drilling, production, and petroleum economics.',        'engineering'),
				('mass-communication',       'Mass Communication',              'Journalism, broadcasting, advertising, PR, and digital media.',                'arts'),
				-- Polytechnic / vocational subjects
				('computer-engineering-tech',  'Computer Engineering Technology', 'Digital electronics, microprocessors, hardware maintenance, and networks.',    'technology'),
				('science-laboratory-tech',    'Science Laboratory Technology',   'Analytical chemistry, microbiology, laboratory instrumentation, and QA.',     'science')
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
