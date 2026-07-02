package scraper

// NUCNursingScraper implements Scraper for NUC CCMAS B.N.Sc. Nursing Science (100L - 500L)
type NUCNursingScraper struct{}

func NewNUCNursingScraper() *NUCNursingScraper {
	return &NUCNursingScraper{}
}

func (s *NUCNursingScraper) BoardSlug() string   { return "nuc" }
func (s *NUCNursingScraper) SubjectSlug() string { return "nursing-science" }
func (s *NUCNursingScraper) Level() string       { return "tertiary-degree" }
func (s *NUCNursingScraper) SourceURL() string {
	return "https://nuc.edu.ng/ccmas/allied-health/nursing-science"
}

func (s *NUCNursingScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 - 200 Level: Foundation Science, Anatomy, and Fundamentals of Nursing",
			Description: "Anatomy, Physiology, Biochemistry, Microbiology, and Foundations of Professional Nursing.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "NUR 201: Foundations of Nursing and Health Assessment", Objectives: []string{
					"Understand nursing history, nursing process (Assessment, Diagnosis, Planning, Implementation, Evaluation)",
					"Perform vital signs measurement, physical health assessment, and aseptic technique infection control",
				}},
				{Name: "ANAT 201 & PHYS 201: Human Anatomy and Physiology for Nurses", Objectives: []string{
					"Study structural anatomy and physiological mechanisms of human organ systems",
				}},
			},
		},
		{
			Name:        "300 Level: Medical-Surgical Nursing and Pharmacology",
			Description: "Medical-Surgical Nursing I & II, Pharmacology, Pathology, and Clinical Practicum.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "NUR 301: Medical-Surgical Nursing", Objectives: []string{
					"Manage nursing care for cardiovascular, respiratory, endocrine, renal, and gastrointestinal diseases",
					"Perform peri-operative nursing care, wound management, and emergency clinical procedures",
				}},
				{Name: "PHARM 301: Pharmacology in Nursing Practice", Objectives: []string{
					"Calculate drug dosages, administer medications via oral, IV, IM routes safely, and monitor side effects",
				}},
			},
		},
		{
			Name:        "400 - 500 Level: Maternal, Paediatric, Mental Health Nursing, and Final Board Defense",
			Description: "Maternal & Child Health Nursing, Psychiatric Nursing, Community Health Nursing, and B.N.Sc. Thesis.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "NUR 401: Maternal & Child Health Nursing (Midwifery)", Objectives: []string{
					"Provide comprehensive antenatal, intrapartum, postpartum nursing care, and newborn care",
				}},
				{Name: "NUR 501: Mental Health, Community Nursing, and Research Project", Objectives: []string{
					"Manage psychiatric nursing interventions for mental health disorders",
					"Execute community health assessment, epidemiology, disease surveillance, and defend original B.N.Sc. thesis",
				}},
			},
		},
	}
}
