-- AfriLearn Curriculum API - Database Schema
-- Run this against your PostgreSQL database to set up the schema

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- EXAM BOARDS (WAEC, JAMB, NECO, NERDC, etc.)
-- ============================================================
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
);

-- ============================================================
-- SUBJECTS (Mathematics, Physics, Biology, etc.)
-- ============================================================
CREATE TABLE IF NOT EXISTS subjects (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slug        VARCHAR(100) NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    category    VARCHAR(50) NOT NULL DEFAULT 'general', -- science, arts, commercial, general
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- CURRICULA (links an exam board to a subject)
-- ============================================================
CREATE TABLE IF NOT EXISTS curricula (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    exam_board_id   UUID NOT NULL REFERENCES exam_boards(id) ON DELETE CASCADE,
    subject_id      UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    year            INTEGER NOT NULL DEFAULT EXTRACT(YEAR FROM NOW()),
    level           VARCHAR(50) NOT NULL DEFAULT 'senior-secondary', -- primary, junior-secondary, senior-secondary
    source_url      VARCHAR(500),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(exam_board_id, subject_id, year)
);

-- ============================================================
-- TOPICS (major sections, e.g. "Algebra", "Mensuration")
-- ============================================================
CREATE TABLE IF NOT EXISTS topics (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    curriculum_id   UUID NOT NULL REFERENCES curricula(id) ON DELETE CASCADE,
    slug            VARCHAR(200) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    order_index     INTEGER NOT NULL DEFAULT 0,
    difficulty      VARCHAR(20) NOT NULL DEFAULT 'medium', -- easy, medium, hard
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(curriculum_id, slug)
);

-- ============================================================
-- SUBTOPICS (specific areas, e.g. "Quadratic Equations")
-- ============================================================
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
);

-- ============================================================
-- LEARNING OBJECTIVES (what students should be able to do)
-- ============================================================
CREATE TABLE IF NOT EXISTS learning_objectives (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    subtopic_id UUID NOT NULL REFERENCES subtopics(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    verb        VARCHAR(50) NOT NULL DEFAULT 'understand', -- Bloom's taxonomy verbs
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- INDEXES for performance
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_curricula_exam_board ON curricula(exam_board_id);
CREATE INDEX IF NOT EXISTS idx_curricula_subject    ON curricula(subject_id);
CREATE INDEX IF NOT EXISTS idx_topics_curriculum    ON topics(curriculum_id);
CREATE INDEX IF NOT EXISTS idx_subtopics_topic      ON subtopics(topic_id);
CREATE INDEX IF NOT EXISTS idx_objectives_subtopic  ON learning_objectives(subtopic_id);

-- ============================================================
-- SEED DATA - Exam Boards
-- ============================================================
INSERT INTO exam_boards (slug, name, full_name, country, description, website) VALUES
    ('waec',  'WAEC',  'West African Examinations Council',         'Nigeria', 'The body that conducts the WASSCE and other examinations across West Africa.', 'https://waec.org.ng'),
    ('jamb',  'JAMB',  'Joint Admissions and Matriculation Board',  'Nigeria', 'The body responsible for university entrance examinations in Nigeria.',          'https://jamb.gov.ng'),
    ('neco',  'NECO',  'National Examinations Council',             'Nigeria', 'The body that conducts the SSCE and other examinations in Nigeria.',             'https://neco.gov.ng'),
    ('nerdc', 'NERDC', 'Nigerian Educational Research and Development Council', 'Nigeria', 'The body responsible for developing the national curriculum.', 'https://nerdc.gov.ng')
ON CONFLICT (slug) DO NOTHING;

-- ============================================================
-- SEED DATA - Core Subjects
-- ============================================================
INSERT INTO subjects (slug, name, description, category) VALUES
    ('mathematics',     'Mathematics',     'The study of numbers, quantities, shapes, and patterns.',                             'science'),
    ('english-language','English Language','The study of the English language, literature, and communication.',                    'arts'),
    ('physics',         'Physics',         'The study of matter, energy, and the fundamental forces of the universe.',            'science'),
    ('chemistry',       'Chemistry',       'The study of substances, their properties, and the reactions between them.',          'science'),
    ('biology',         'Biology',         'The study of living organisms and their interactions with the environment.',           'science'),
    ('economics',       'Economics',       'The study of how societies allocate scarce resources.',                                'commercial'),
    ('government',      'Government',      'The study of political systems, governance, and civic responsibilities.',             'arts'),
    ('literature',      'Literature in English', 'The study and analysis of literary texts in English.',                         'arts'),
    ('geography',       'Geography',       'The study of the physical features of the earth and human activity.',                 'arts'),
    ('further-mathematics', 'Further Mathematics', 'Advanced mathematics covering topics beyond the standard syllabus.',         'science')
ON CONFLICT (slug) DO NOTHING;
