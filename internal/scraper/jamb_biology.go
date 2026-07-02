package scraper

// JAMBBiologyScraper implements Scraper for JAMB Biology
type JAMBBiologyScraper struct{}

func NewJAMBBiologyScraper() *JAMBBiologyScraper {
	return &JAMBBiologyScraper{}
}

func (s *JAMBBiologyScraper) BoardSlug() string   { return "jamb" }
func (s *JAMBBiologyScraper) SubjectSlug() string { return "biology" }
func (s *JAMBBiologyScraper) Level() string       { return "tertiary-entry" }
func (s *JAMBBiologyScraper) SourceURL() string {
	return "https://ibass.jamb.gov.ng/syllabus/biology"
}

func (s *JAMBBiologyScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Variety of Organisms and Cellular Form",
			Description: "JAMB UTME Section 1: Kingdom classification, cell structure, organelle functions, and cellular transport.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Organism Classification and Cell Structure", Objectives: []string{
					"Classify organisms from Monera to Animalia with key structural features",
					"Compare plant cell and animal cell organelles using light and electron microscopy",
				}},
				{Name: "Cellular Transport Mechanisms", Objectives: []string{
					"Explain diffusion, osmosis, plasmolysis, turgidity, and active transport in biological systems",
				}},
			},
		},
		{
			Name:        "Physiological Processes in Living Organisms",
			Description: "JAMB UTME Section 2: Nutrition, transport systems, respiration, excretion, movement, and growth regulation.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Plant and Animal Nutrition", Objectives: []string{
					"Describe autotrophic (photosynthesis) and heterotrophic nutrition",
					"Identify human digestive system structures, enzymes, and balanced diet requirements",
				}},
				{Name: "Internal Transport and Respiration", Objectives: []string{
					"Describe mammalian blood circulation, heart structure, and plant vascular transport (xylem/phloem)",
					"Compare aerobic and anaerobic respiration in cells and respiratory organs in animals",
				}},
				{Name: "Excretion, Movement, and Hormones", Objectives: []string{
					"Describe mammalian kidney structure, nephron function, and excretion in plants",
					"Explain skeletal support systems and hormonal regulation in plants (auxins) and mammals",
				}},
			},
		},
		{
			Name:        "Ecology, Heredity, and Evolution",
			Description: "JAMB UTME Section 3: Ecosystems, food webs, pollution, Mendelian genetics, variation, and evolution.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "Ecology and Conservation", Objectives: []string{
					"Analyze food chains, food webs, energy flow, ecological pyramids, and nutrient cycles",
					"Identify air, water, and soil pollution causes, effects, and conservation strategies",
				}},
				{Name: "Genetics, Variation, and Evolution", Objectives: []string{
					"Apply Mendel's laws of inheritance to monohybrid and dihybrid crosses",
					"Distinguish continuous and discontinuous variation and human sex-linked traits",
					"Compare Lamarckism, Darwinian natural selection, and evidence of evolution",
				}},
			},
		},
	}
}
