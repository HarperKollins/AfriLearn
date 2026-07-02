package scraper

// NUCMedicineScraper implements Scraper for NUC CCMAS M.B.B.S. Medicine and Surgery (100L - 500L)
type NUCMedicineScraper struct{}

func NewNUCMedicineScraper() *NUCMedicineScraper {
	return &NUCMedicineScraper{}
}

func (s *NUCMedicineScraper) BoardSlug() string   { return "nuc" }
func (s *NUCMedicineScraper) SubjectSlug() string { return "medicine-and-surgery" }
func (s *NUCMedicineScraper) Level() string       { return "tertiary-degree" }
func (s *NUCMedicineScraper) SourceURL() string {
	return "https://nuc.edu.ng/ccmas/medicine/medicine-and-surgery"
}

func (s *NUCMedicineScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 Level: Pre-Medical Foundation Sciences",
			Description: "General Biology, General Chemistry, General Physics, and Medical Mathematics foundation courses.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "BIO 101 & CHE 101: Cell Biology and Organic Chemistry", Objectives: []string{
					"Understand cell ultrastructure, biomolecules, genetics, and metabolic pathways",
					"Study functional groups, stereochemistry, and organic reactions relevant to biological systems",
				}},
				{Name: "PHY 101: Physics for Medical Sciences", Objectives: []string{
					"Apply fluid mechanics, optics, wave motion, radiation physics, and thermodynamics to human physiology",
				}},
			},
		},
		{
			Name:        "200 Level: Basic Medical Sciences (1st Professional MBBS Part 1)",
			Description: "Gross Anatomy, Embryology, Histology, Medical Biochemistry, and Human Physiology.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "ANAT 201: Gross Anatomy, Histology, and Embryology", Objectives: []string{
					"Dissect and identify anatomical structures of upper limbs, lower limbs, thorax, abdomen, pelvis, head, and neck",
					"Identify microscopic tissue architecture under light microscopy (Histology)",
					"Trace human embryonic development from fertilization to organogenesis (Embryology)",
				}},
				{Name: "BIOC 201 & PHYS 201: Medical Biochemistry and Human Physiology", Objectives: []string{
					"Study carbohydrate, lipid, protein, and nucleic acid metabolism and enzymatic regulation",
					"Analyze physiological mechanisms of cardiovascular, respiratory, renal, gastrointestinal, and endocrine systems",
				}},
			},
		},
		{
			Name:        "300 Level: Pathology and Pharmacology (2nd Professional MBBS Part 2)",
			Description: "Anatomical Pathology, Haematology, Chemical Pathology, Medical Microbiology, and Pharmacology.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "PATH 301: General and Systemic Pathology", Objectives: []string{
					"Understand mechanisms of cellular injury, inflammation, tissue repair, neoplasia, and hemodynamic disorders",
					"Analyze disease etiologies and histopathological changes in major organ systems",
				}},
				{Name: "PHARM 301: Pharmacology and Therapeutics", Objectives: []string{
					"Study pharmacokinetics (absorption, distribution, metabolism, excretion) and pharmacodynamics (drug-receptor interactions)",
					"Analyze therapeutic actions and toxicities of autonomic, cardiovascular, antimicrobial, and CNS drugs",
				}},
			},
		},
		{
			Name:        "400 - 500 Level: Clinical Rotations and Final MBBS",
			Description: "Obstetrics & Gynaecology, Paediatrics, Community Medicine, Internal Medicine, General Surgery, and Final Board Defense.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "CLIN 401: Obstetrics, Gynaecology, and Paediatrics", Objectives: []string{
					"Manage antenatal care, labor, delivery, obstetric emergencies, and gynaecological disorders",
					"Diagnose and treat neonatal disorders, childhood infections, malnutrition, and developmental pediatrics",
				}},
				{Name: "CLIN 501: Internal Medicine and General Surgery", Objectives: []string{
					"Conduct clinical history taking, physical examinations, and differential diagnosis in internal medicine",
					"Perform surgical ward rounds, emergency resuscitation, and basic operative surgical procedures",
					"Pass Final MBBS clinical ward examinations and oral viva defense before Medical and Dental Council of Nigeria (MDCN) examiners",
				}},
			},
		},
	}
}
