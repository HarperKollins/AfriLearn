-- AfriLearn Curriculum API - Database Schema
-- Run this against your PostgreSQL database to set up the schema

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- EXAM BOARDS & INSTITUTIONS (WAEC, JAMB, NECO, BECE, NERDC, NUC, NBTE, Polytechnics & Universities)
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
-- SUBJECTS & DEGREE/DIPLOMA PROGRAMS
-- ============================================================
CREATE TABLE IF NOT EXISTS subjects (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slug        VARCHAR(100) NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    category    VARCHAR(50) NOT NULL DEFAULT 'general', -- science, arts, commercial, basic, computing, engineering, medical, law, pharmacy, agriculture, environment, polytechnic
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- CURRICULA (links an exam board / regulatory body / university / polytechnic to a subject or diploma program)
-- ============================================================
CREATE TABLE IF NOT EXISTS curricula (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    exam_board_id   UUID NOT NULL REFERENCES exam_boards(id) ON DELETE CASCADE,
    subject_id      UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    year            INTEGER NOT NULL DEFAULT EXTRACT(YEAR FROM NOW()),
    level           VARCHAR(50) NOT NULL DEFAULT 'senior-secondary', -- junior-secondary, senior-secondary, tertiary-entry, tertiary-degree, polytechnic-nd, polytechnic-hnd
    source_url      VARCHAR(500),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(exam_board_id, subject_id, year)
);

-- ============================================================
-- TOPICS / COURSE MODULES (major sections or course codes e.g. "CTE 111", "SLT 121")
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
-- SEED DATA - Exam Boards, Regulatory Bodies, Polytechnics & Universities
-- ============================================================
INSERT INTO exam_boards (slug, name, full_name, country, description, website) VALUES
    -- Secondary Exam Boards
    ('bece',        'BECE',        'Basic Education Certificate Examination',   'Nigeria', 'The national exam for Junior Secondary School graduation (JSS3).', 'https://neco.gov.ng'),
    ('waec',        'WAEC',        'West African Examinations Council',         'Nigeria', 'The body that conducts the WASSCE across West Africa.',            'https://waec.org.ng'),
    ('jamb',        'JAMB',        'Joint Admissions and Matriculation Board',  'Nigeria', 'The body responsible for university entrance exams in Nigeria.',     'https://jamb.gov.ng'),
    ('neco',        'NECO',        'National Examinations Council',             'Nigeria', 'The national body that conducts SSCE and BECE exams in Nigeria.',   'https://neco.gov.ng'),
    ('nerdc',       'NERDC',       'Nigerian Educational Research & Dev Council','Nigeria', 'The statutory body that develops the national curriculum.',       'https://nerdc.gov.ng'),

    -- Tertiary Regulatory Bodies
    ('nuc',         'NUC',         'National Universities Commission (CCMAS)',  'Nigeria', 'Regulates university education and sets the 70% core CCMAS standards for all 270+ Nigerian universities.', 'https://nuc.edu.ng'),
    ('nbte',        'NBTE',        'National Board for Technical Education',    'Nigeria', 'Regulates polytechnic and monotechnic ND/HND education in Nigeria.', 'https://nbte.gov.ng'),

    -- Polytechnics (ND & HND)
    ('yabatech',    'YABATECH',    'Yaba College of Technology',               'Nigeria', 'Nigeria''s premier polytechnic located in Yaba, Lagos State.',    'https://yabatech.edu.ng'),
    ('imt',         'IMT',         'Institute of Management and Technology',   'Nigeria', 'Leading polytechnic located in Enugu, Enugu State.',             'https://imt.edu.ng'),
    ('auchi',       'AUCHI',       'Auchi Polytechnic',                        'Nigeria', 'Federal polytechnic located in Auchi, Edo State.',               'https://auchipoly.edu.ng'),
    ('fedpoly-nek', 'NEKEDEPOLY',  'Federal Polytechnic, Nekede',              'Nigeria', 'Federal polytechnic located in Owerri, Imo State.',              'https://fpno.edu.ng'),

    -- Federal & State Universities
    ('ebsu',        'EBSU',        'Ebonyi State University',                   'Nigeria', 'State university located in Abakaliki, Ebonyi State.',            'https://ebsu.edu.ng'),
    ('funai',       'AE-FUNAI',    'Alex Ekwueme Federal University, Ndufu-Alike','Nigeria','Federal university located in Ikwo, Ebonyi State.',             'https://funai.edu.ng'),
    ('unn',         'UNN',         'University of Nigeria, Nsukka',             'Nigeria', 'First autonomous federal university located in Nsukka, Enugu State.', 'https://unn.edu.ng'),
    ('unec',        'UNEC',        'University of Nigeria, Enugu Campus',       'Nigeria', 'Enugu campus of UNN housing Law, Business Admin, and Medical Sciences.', 'https://unn.edu.ng'),
    ('unilag',      'UNILAG',      'University of Lagos',                       'Nigeria', 'Premier federal university in Yaba, Lagos State.',                'https://unilag.edu.ng'),
    ('ui',          'UI',          'University of Ibadan',                      'Nigeria', 'Nigeria''s first federal university located in Ibadan, Oyo State.', 'https://ui.edu.ng'),
    ('oau',         'OAU',         'Obafemi Awolowo University',                'Nigeria', 'Premier federal university located in Ile-Ife, Osun State.',       'https://oauife.edu.ng'),
    ('abu',         'ABU',         'Ahmadu Bello University',                   'Nigeria', 'Premier federal university located in Zaria, Kaduna State.',        'https://abu.edu.ng'),
    ('futo',        'FUTO',        'Federal University of Technology, Owerri',  'Nigeria', 'Premier university of technology in Owerri, Imo State.',          'https://futo.edu.ng'),
    ('futa',        'FUTA',        'Federal University of Technology, Akure',   'Nigeria', 'Premier university of technology in Akure, Ondo State.',          'https://futa.edu.ng'),
    ('covenant',    'CU',          'Covenant University',                       'Nigeria', 'Leading private university located in Ota, Ogun State.',          'https://covenantuniversity.edu.ng')
ON CONFLICT (slug) DO NOTHING;

-- ============================================================
-- SEED DATA - Subjects, University Degrees & Polytechnic Diplomas
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

    -- NUC CCMAS Computing Discipline
    ('computer-science',     'B.Sc. Computer Science',               'NUC CCMAS degree program covering programming, algorithms, systems, AI, and software engineering.', 'computing'),
    ('software-engineering', 'B.Sc. Software Engineering',           'NUC CCMAS degree program covering software architecture, testing, DevOps, and project management.',  'computing'),
    ('cyber-security',       'B.Sc. Cybersecurity',                  'NUC CCMAS degree program covering network security, cryptography, forensics, and ethical hacking.',  'computing'),
    ('data-science',         'B.Sc. Data Science',                   'NUC CCMAS degree program covering machine learning, big data, data mining, and analytics.',          'computing'),

    -- NUC CCMAS Medical & Allied Health Disciplines
    ('medicine-and-surgery', 'M.B.B.S. Medicine and Surgery',        'NUC CCMAS medical degree program covering anatomy, physiology, pathology, and clinical medicine.',  'medical'),
    ('nursing-science',      'B.N.Sc. Nursing Science',              'NUC CCMAS degree program covering clinical nursing, anatomy, pharmacology, and patient care.',      'medical'),
    ('pharmacy',             'Pharm.D. Doctor of Pharmacy',          'NUC CCMAS professional pharmacy program covering pharmaceutics, pharmacology, and clinical pharmacy.','pharmacy'),
    ('microbiology',         'B.Sc. Microbiology',                   'NUC CCMAS degree program covering bacteriology, virology, immunology, and industrial microbiology.','science'),
    ('biochemistry',         'B.Sc. Biochemistry',                   'NUC CCMAS degree program covering enzymology, metabolism, molecular biology, and clinical biochemistry.','science'),

    -- NUC CCMAS Engineering Discipline
    ('electrical-engineering','B.Eng. Electrical & Electronic Eng',  'NUC CCMAS engineering program covering circuits, power systems, electronics, and telecoms.',        'engineering'),
    ('mechanical-engineering','B.Eng. Mechanical Engineering',       'NUC CCMAS engineering program covering thermodynamics, fluid mechanics, and machine design.',       'engineering'),
    ('civil-engineering',    'B.Eng. Civil Engineering',             'NUC CCMAS engineering program covering structures, hydraulics, soil mechanics, and highway eng.',   'engineering'),
    ('chemical-engineering', 'B.Eng. Chemical Engineering',          'NUC CCMAS engineering program covering unit operations, transport phenomena, and reaction eng.',    'engineering'),
    ('petroleum-engineering','B.Eng. Petroleum Engineering',         'NUC CCMAS engineering program covering reservoir engineering, drilling, and production.',          'engineering'),

    -- NUC CCMAS Law Discipline
    ('law',                  'LL.B. Bachelor of Laws',               'NUC CCMAS law degree program covering constitutional law, criminal law, contract, and CAMA.',       'law'),

    -- NUC CCMAS Administration & Social Sciences Disciplines
    ('accounting',           'B.Sc. Accounting',                     'NUC CCMAS degree program covering financial accounting, auditing, taxation, and ICAN standards.',   'commercial'),
    ('business-administration','B.Sc. Business Administration',      'NUC CCMAS degree program covering management, marketing, organizational behavior, and strategy.',   'commercial'),
    ('economics-degree',     'B.Sc. Economics (University)',         'NUC CCMAS university economics program covering econometrics, macro, micro, and public finance.',    'commercial'),
    ('political-science',    'B.Sc. Political Science',              'NUC CCMAS degree program covering political theory, public admin, and international relations.',    'arts'),
    ('mass-communication',   'B.Sc. Mass Communication',             'NUC CCMAS degree program covering journalism, broadcasting, public relations, and digital media.',          'arts'),

    -- NUC CCMAS Architecture & Agriculture Disciplines
    ('architecture',         'B.Sc. Architecture',                   'NUC CCMAS degree program covering architectural design studio, building construction, and CAD.',    'environment'),
    ('agriculture-degree',   'B.Agric. Agriculture',                 'NUC CCMAS degree program covering crop science, animal science, soil science, and ag-economics.',  'agriculture'),

    -- NBTE Polytechnic National Diploma (ND) & Higher National Diploma (HND) Programs
    ('computer-engineering-tech', 'ND/HND Computer Engineering Tech', 'NBTE polytechnic diploma program covering digital systems, hardware repair, microprocessors, and networking.', 'polytechnic'),
    ('science-laboratory-tech',   'ND/HND Science Laboratory Tech (SLT)', 'NBTE polytechnic diploma program covering analytical chemistry, biochemistry, and lab techniques.',  'polytechnic'),
    ('electrical-telecoms-tech',  'ND/HND Electrical & Telecoms Tech',   'NBTE polytechnic diploma program covering power electronics, telecommunications, and high voltage.',  'polytechnic')
ON CONFLICT (slug) DO NOTHING;
