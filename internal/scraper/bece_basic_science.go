package scraper

// BECEBasicScienceScraper implements Scraper for BECE Basic Science (JSS1 - JSS3)
type BECEBasicScienceScraper struct{}

func NewBECEBasicScienceScraper() *BECEBasicScienceScraper {
	return &BECEBasicScienceScraper{}
}

func (s *BECEBasicScienceScraper) BoardSlug() string   { return "bece" }
func (s *BECEBasicScienceScraper) SubjectSlug() string { return "basic-science" }
func (s *BECEBasicScienceScraper) Level() string       { return "junior-secondary" }
func (s *BECEBasicScienceScraper) SourceURL() string {
	return "https://nerdc.gov.ng/curriculum/junior-secondary/basic-science"
}

func (s *BECEBasicScienceScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Living Things and Human Health (JSS1 - JSS3)",
			Description: "Living organisms, cellular organization, human body systems, reproductive health, diseases, and drug abuse.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Living Things and Cell Organization", Objectives: []string{
					"Distinguish characteristics of living and non-living things",
					"Identify plant cell and animal cell organelles and their functions",
					"Explain organizational levels of living things from cells to organisms",
				}},
				{Name: "Human Body Systems and Health", Objectives: []string{
					"Describe digestive, circulatory, respiratory, excretory, and reproductive systems in humans",
					"Explain pubertal changes, personal hygiene, and reproductive health in adolescents",
					"Identify transmission, prevention, and control of communicable diseases (malaria, HIV/AIDS, cholera)",
				}},
				{Name: "Drug and Substance Abuse", Objectives: []string{
					"Define drug abuse, misuse, addiction, and prescription drugs",
					"Identify social, psychological, and physiological consequences of drug abuse",
				}},
			},
		},
		{
			Name:        "Matter, Chemical Changes, and Energy (JSS1 - JSS3)",
			Description: "States of matter, elements, compounds, mixtures, chemical symbols, kinetic theory, energy forms, and forces.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Matter and Atomic Structure", Objectives: []string{
					"Describe solid, liquid, and gas states using Kinetic Theory of Matter",
					"Distinguish elements, compounds, and mixtures with examples",
					"Write chemical symbols and formulas for common elements and compounds",
				}},
				{Name: "Chemical Changes and Reactions", Objectives: []string{
					"Distinguish physical changes (melting, boiling) and chemical changes (rusting, burning)",
					"Identify acids, bases, and salts using litmus paper and pH indicators",
				}},
				{Name: "Energy, Work, and Power", Objectives: []string{
					"Identify mechanical, thermal, electrical, chemical, and solar energy forms",
					"Apply law of conservation of energy to energy transformations",
					"Calculate work done (W = F × d) and power (P = W / t)",
				}},
			},
		},
		{
			Name:        "Earth, Space, and Environmental Science (JSS1 - JSS3)",
			Description: "Solar system, Earth rotation/revolution, weather, climate, environmental pollution, and natural resources.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Solar System and Earth Movements", Objectives: []string{
					"Identify planets in the solar system, sun, moon, and satellites",
					"Explain rotation of Earth (day and night) and revolution of Earth (seasons)",
					"Describe solar and lunar eclipses",
				}},
				{Name: "Environmental Pollution and Protection", Objectives: []string{
					"Identify sources and consequences of air, water, and land pollution",
					"Explain deforestation, desertification, soil erosion, and climate change effects",
					"Describe refuse disposal, recycling, and environmental sanitation practices",
				}},
			},
		},
	}
}
