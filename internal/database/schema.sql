-- AfriLearn Curriculum API - Database Schema
-- Run this against your PostgreSQL database to set up the schema

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- EXAM BOARDS & INSTITUTIONS (WAEC, JAMB, NECO, BECE, NERDC, NUC, UNILAG, etc.)
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
-- SUBJECTS & DEGREE PROGRAMS (Secondary Subjects + University Degree Programs)
-- ============================================================
CREATE TABLE IF NOT EXISTS subjects (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slug        VARCHAR(100) NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    category    VARCHAR(50) NOT NULL DEFAULT 'general', -- science, arts, commercial, basic, computing, engineering, medical, law
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- CURRICULA (links an exam board / regulatory body / university to a subject or degree program)
-- ============================================================
CREATE TABLE IF NOT EXISTS curricula (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    exam_board_id   UUID NOT NULL REFERENCES exam_boards(id) ON DELETE CASCADE,
    subject_id      UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    year            INTEGER NOT NULL DEFAULT EXTRACT(YEAR FROM NOW()),
    level           VARCHAR(50) NOT NULL DEFAULT 'senior-secondary', -- junior-secondary, senior-secondary, tertiary-entry, 100-level, 200-level, 300-level, 400-level, 500-level, tertiary-degree
    source_url      VARCHAR(500),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(exam_board_id, subject_id, year)
);

-- ============================================================
-- TOPICS / COURSE MODULES (major sections or university course codes e.g. "COS 101", "CSC 201")
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
-- SUBTOPICS / COURSE UNITS (specific subtopics or lecture modules)
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
-- SEED DATA - Exam Boards & Higher Ed Bodies / Institutions
-- ============================================================
INSERT INTO exam_boards (slug, name, full_name, country, description, website) VALUES
    -- Secondary Exam Boards
    ('bece',    'BECE',    'Basic Education Certificate Examination',   'Nigeria', 'The national exam for Junior Secondary School graduation (JSS3).', 'https://neco.gov.ng'),
    ('waec',    'WAEC',    'West African Examinations Council',         'Nigeria', 'The body that conducts the WASSCE across West Africa.',            'https://waec.org.ng'),
    ('jamb',    'JAMB',    'Joint Admissions and Matriculation Board',  'Nigeria', 'The body responsible for university entrance exams in Nigeria.',     'https://jamb.gov.ng'),
    ('neco',    'NECO',    'National Examinations Council',             'Nigeria', 'The national body that conducts SSCE and BECE exams in Nigeria.',   'https://neco.gov.ng'),
    ('nerdc',   'NERDC',   'Nigerian Educational Research & Dev Council','Nigeria', 'The statutory body that develops the national curriculum.',       'https://nerdc.gov.ng'),

    -- Tertiary Regulatory Bodies & Flagship Universities
    ('nuc',     'NUC',     'National Universities Commission (CCMAS)',  'Nigeria', 'Regulates university education and sets the 70% core CCMAS standards for all Nigerian universities.', 'https://nuc.edu.ng'),
    ('nbte',    'NBTE',    'National Board for Technical Education',    'Nigeria', 'Regulates polytechnic and monotechnic ND/HND education in Nigeria.', 'https://nbte.gov.ng'),
    ('unilag',  'UNILAG',  'University of Lagos',                       'Nigeria', 'Premier federal university in Yaba, Lagos State.',                'https://unilag.edu.ng'),
    ('ui',      'UI',      'University of Ibadan',                      'Nigeria', 'Nigeria''s first federal university located in Ibadan, Oyo State.', 'https://ui.edu.ng'),
    ('oau',     'OAU',     'Obafemi Awolowo University',                'Nigeria', 'Premier federal university located in Ile-Ife, Osun State.',       'https://oauife.edu.ng'),
    ('unn',     'UNN',     'University of Nigeria, Nsukka',             'Nigeria', 'First autonomous federal university located in Nsukka, Enugu State.', 'https://unn.edu.ng'),
    ('abu',     'ABU',     'Ahmadu Bello University',                   'Nigeria', 'Premier federal university located in Zaria, Kaduna State.',        'https://abu.edu.ng'),
    ('covenant','CU',      'Covenant University',                       'Nigeria', 'Leading private university located in Ota, Ogun State.',          'https://covenantuniversity.edu.ng')
ON CONFLICT (slug) DO NOTHING;

-- ============================================================
-- SEED DATA - Subjects & University Degree Programs
-- ============================================================
INSERT INTO subjects (slug, name, description, category) VALUES
    -- Senior Secondary Subjects
    ('mathematics',          'Mathematics',          'Core senior secondary mathematics syllabus.',                                'science'),
    ('english-language',     'English Language',     'Core senior secondary English language and communication.',                   'arts'),
    ('physics',              'Physics',              'Study of matter, energy, mechanics, electricity, and modern physics.',        'science'),
    ('chemistry',            'Chemistry',            'Study of atomic structure, chemical bonding, reactions, and organic chemistry.','science'),
    ('biology',              'Biology',              'Study of living organisms, cellular processes, genetics, and ecology.',       'science'),
    ('economics',            'Economics',            'Microeconomics, macroeconomics, trade, and economic principles.',            'commercial'),
    ('government',           'Government',           'Political institutions, governance, constitutions, and international relations.','arts'),
    ('literature',           'Literature in English','African and non-African prose, poetry, and drama analysis.',                 'arts'),
    ('geography',            'Geography',            'Physical, human, regional, and practical map reading geography.',             'arts'),
    ('further-mathematics',  'Further Mathematics',  'Advanced pure mathematics, mechanics, and statistics.',                       'science'),
    ('agricultural-science', 'Agricultural Science', 'Crop production, animal husbandry, soil science, and farm management.',        'science'),
    ('civic-education',      'Civic Education',      'Values, human rights, citizenship, democracy, and national consciousness.',    'arts'),
    ('commerce',             'Commerce',             'Trade, business operations, banking, insurance, and marketing.',             'commercial'),
    ('financial-accounting', 'Financial Accounting', 'Bookkeeping, financial statements, partnership, and corporate accounting.',  'commercial'),
    ('computer-studies',     'Computer Studies',     'Computer hardware, software, networking, programming, and ICT applications.', 'science'),

    -- Junior Secondary Subjects (JSS1-JSS3 / BECE)
    ('basic-science',        'Basic Science',        'Integrated junior secondary physical, chemical, and biological science.',     'basic'),
    ('basic-technology',     'Basic Technology',     'Materials, tools, woodwork, metalwork, technical drawing, and electronics.',  'basic'),
    ('social-studies',       'Social Studies',       'Human relationships, culture, family, social issues, and environment.',      'basic'),
    ('business-studies',     'Business Studies',     'Office practice, book-keeping, shorthand, commerce, and consumer education.', 'basic'),
    ('cultural-and-creative-arts', 'Cultural & Creative Arts', 'Visual arts, drama, music, and Nigerian cultural heritage.',        'basic'),
    ('physical-and-health-education', 'Physical & Health Education', 'Physical fitness, athletics, games, safety, and health education.', 'basic'),

    -- University Degree Programs (NUC CCMAS Standards)
    ('computer-science',            'B.Sc. Computer Science',               'NUC CCMAS degree program covering programming, algorithms, systems, AI, and software engineering.', 'computing'),
    ('cyber-security',              'B.Sc. Cyber Security',                 'NUC CCMAS degree program covering network security, cryptography, forensics, and ethical hacking.',  'computing'),
    ('software-engineering',        'B.Sc. Software Engineering',           'NUC CCMAS degree program covering software architecture, testing, DevOps, and project management.',  'computing'),
    ('medicine-and-surgery',        'M.B.B.S. Medicine and Surgery',        'NUC CCMAS professional medical degree program covering anatomy, physiology, pathology, and clinical medicine.', 'medical'),
    ('nursing-science',             'B.N.Sc. Nursing Science',              'NUC CCMAS professional degree program covering clinical nursing, anatomy, pharmacology, and patient care.',  'medical'),
    ('electrical-engineering',      'B.Eng. Electrical & Electronic Eng',   'NUC CCMAS engineering program covering circuit theory, power systems, electronics, and telecommunications.','engineering'),
    ('mechanical-engineering',      'B.Eng. Mechanical Engineering',       'NUC CCMAS engineering program covering thermodynamics, fluid mechanics, machine design, and manufacturing.', 'engineering'),
    ('law',                         'LL.B. Bachelor of Laws',               'NUC CCMAS professional law degree program covering constitutional law, criminal law, contract law, and jurisprudence.', 'law'),
    ('accounting',                  'B.Sc. Accounting',                     'NUC CCMAS degree program covering financial accounting, auditing, taxation, and management accounting.',    'commercial'),
    ('mass-communication',          'B.Sc. Mass Communication',             'NUC CCMAS degree program covering journalism, broadcasting, public relations, and digital media.',          'arts')
ON CONFLICT (slug) DO NOTHING;
