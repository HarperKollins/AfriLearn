package scraper

// WAECChemistryScraper implements Scraper for WAEC Chemistry (SS1 - SS3)
type WAECChemistryScraper struct{}

func NewWAECChemistryScraper() *WAECChemistryScraper {
	return &WAECChemistryScraper{}
}

func (s *WAECChemistryScraper) BoardSlug() string   { return "waec" }
func (s *WAECChemistryScraper) SubjectSlug() string { return "chemistry" }
func (s *WAECChemistryScraper) Level() string       { return "senior-secondary" }
func (s *WAECChemistryScraper) SourceURL() string {
	return "https://waecsyllabus.com/chemistry-syllabus/"
}

func (s *WAECChemistryScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Atomic Structure, Periodic Table, and Bonding",
			Description: "Subatomic particles, isotopes, electron configuration, periodic trends, electrovalent, covalent, metallic, and hydrogen bonding.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Atomic Structure and Isotopy", Objectives: []string{
					"Identify protons, neutrons, electrons, atomic number Z, and mass number A",
					"Write s, p, d electron configurations for elements 1 to 30",
					"Calculate relative atomic mass RAM from isotopic abundances",
				}},
				{Name: "Periodic Table and Periodic Trends", Objectives: []string{
					"Classify elements into Groups I-VIII and Periods 1-4",
					"Explain periodic trends: atomic radius, ionization energy, electronegativity, electron affinity",
				}},
				{Name: "Chemical Bonding", Objectives: []string{
					"Compare electrovalent (ionic), covalent, co-ordinate (dative), metallic, and hydrogen bonds",
					"Relate bond type to physical properties (melting point, solubility, electrical conductivity)",
				}},
			},
		},
		{
			Name:        "Stoichiometry, Gas Laws, and Solutions",
			Description: "Mole concept, chemical equations, gas laws, solutions, solubility, and volumetric analysis.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Mole Concept and Chemical Equations", Objectives: []string{
					"Calculate molar mass, empirical formula, molecular formula, and percentage composition",
					"Balance chemical equations and perform stoichiometric calculations",
				}},
				{Name: "Gas Laws and Kinetic Theory", Objectives: []string{
					"Apply Boyle's Law, Charles' Law, Graham's Law of Diffusion, and Dalton's Law of Partial Pressures",
					"Calculate gas volume at standard temperature and pressure (STP = 22.4 dm³/mol)",
				}},
				{Name: "Solubility and Volumetric Analysis", Objectives: []string{
					"Define solute, solvent, saturated, unsaturated, and supersaturated solutions",
					"Plot and interpret solubility curves",
					"Perform acid-base titration calculations (molarity, concentration in g/dm³, purity)",
				}},
			},
		},
		{
			Name:        "Acids, Bases, Salts, and Electrochemistry",
			Description: "Acids, bases, pH scale, salt preparation, oxidation-reduction (redox) reactions, and electrolysis.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "Acids, Bases, and Salts", Objectives: []string{
					"Compare Arrhenius, Bronsted-Lowry, and Lewis concepts of acids and bases",
					"Calculate pH = -log[H+] and pOH of solutions",
					"Classify salts (normal, acid, basic, double, complex) and describe preparation methods",
				}},
				{Name: "Redox Reactions and Electrochemistry", Objectives: []string{
					"Assign oxidation numbers and identify oxidizing and reducing agents",
					"Balance redox equations using half-reaction method",
					"State Faraday's 1st and 2nd Laws of Electrolysis and calculate mass deposited m = Q M / (z F)",
				}},
			},
		},
		{
			Name:        "Organic Chemistry and Macromolecules",
			Description: "Hydrocarbons, alkanols, alkanoic acids, esters, fats, oils, polymers, and synthetic materials.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "Hydrocarbons: Alkanes, Alkenes, Alkynes, and Benzene", Objectives: []string{
					"Name organic compounds using IUPAC nomenclature",
					"Describe fractional distillation of crude oil and cracking of petroleum fractions",
					"Compare structural isomerism and chemical reactions of alkanes, alkenes, alkynes, and benzene",
				}},
				{Name: "Alkanols, Alkanoic Acids, and Esters", Objectives: []string{
					"Describe preparation, properties, and reactions of primary, secondary, and tertiary alkanols",
					"Explain esterification reaction and saponification of fats and oils to produce soap",
				}},
				{Name: "Polymers and Synthetic Materials", Objectives: []string{
					"Distinguish addition polymerization (polyethene, PVC) and condensation polymerization (nylon, terylene)",
					"Describe structure and properties of carbohydrates, proteins, and synthetic polymers",
				}},
			},
		},
	}
}
