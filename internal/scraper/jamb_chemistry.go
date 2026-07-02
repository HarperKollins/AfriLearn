package scraper

// JAMBChemistryScraper implements Scraper for JAMB Chemistry
type JAMBChemistryScraper struct{}

func NewJAMBChemistryScraper() *JAMBChemistryScraper {
	return &JAMBChemistryScraper{}
}

func (s *JAMBChemistryScraper) BoardSlug() string   { return "jamb" }
func (s *JAMBChemistryScraper) SubjectSlug() string { return "chemistry" }
func (s *JAMBChemistryScraper) Level() string       { return "tertiary-entry" }
func (s *JAMBChemistryScraper) SourceURL() string {
	return "https://ibass.jamb.gov.ng/syllabus/chemistry"
}

func (s *JAMBChemistryScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Separation of Mixtures and Atomic Structure",
			Description: "JAMB UTME Topic 1: Physical vs chemical changes, separation techniques, atomic structure, and periodic trends.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Separation Techniques and Purity", Objectives: []string{
					"Select appropriate separation techniques: filtration, crystallization, distillation, chromatography, sublimation",
					"Determine criteria of purity for chemical substances (melting point, boiling point)",
				}},
				{Name: "Atomic Structure and Chemical Bonding", Objectives: []string{
					"Determine atomic numbers, mass numbers, electron configurations, and isotopes",
					"Compare electrovalent, covalent, dative, metallic, and hydrogen bonding",
				}},
			},
		},
		{
			Name:        "Gas Laws, Stoichiometry, and Solutions",
			Description: "JAMB UTME Topic 2: Kinetic theory of gases, ideal gas laws, mole concept, and solubility.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Gas Laws and Kinetic Theory", Objectives: []string{
					"Apply Boyle's Law, Charles' Law, Graham's Law of Diffusion, and Gay-Lussac's Law",
					"Calculate gas volume at STP and molar volume (22.4 dm³)",
				}},
				{Name: "Mole Concept, Solubility, and Titration", Objectives: []string{
					"Calculate empirical formulas, molecular formulas, and stoichiometry of equations",
					"Plot solubility curves and perform volumetric acid-base titration calculations",
				}},
			},
		},
		{
			Name:        "Energetics, Equilibrium, Electrochemistry, and Organic Chemistry",
			Description: "JAMB UTME Topic 3: Enthalpy, rates of reaction, chemical equilibrium, redox, electrolysis, and organic chemistry.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "Chemical Energetics, Kinetics, and Equilibrium", Objectives: []string{
					"Calculate enthalpy change ΔH for exothermic and endothermic reactions",
					"Explain factors affecting reaction rates (temperature, concentration, catalyst)",
					"State Le Chatelier's principle and analyze equilibrium constant Kc shifts",
				}},
				{Name: "Redox Reactions and Electrolysis", Objectives: []string{
					"Determine oxidation states and identify oxidizing/reducing agents",
					"Calculate mass of elements liberated during electrolysis using Faraday's laws",
				}},
				{Name: "Organic Chemistry and IUPAC Nomenclature", Objectives: []string{
					"Name and draw structures for alkanes, alkenes, alkynes, alkanols, alkanoic acids, esters, and benzene",
					"Explain petroleum cracking, octane number, polymerisation, and saponification",
				}},
			},
		},
	}
}
