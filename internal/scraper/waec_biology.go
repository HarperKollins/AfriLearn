package scraper

// WAECBiologyScraper implements Scraper for WAEC Biology
type WAECBiologyScraper struct{}

func NewWAECBiologyScraper() *WAECBiologyScraper {
	return &WAECBiologyScraper{}
}

func (s *WAECBiologyScraper) BoardSlug() string   { return "waec" }
func (s *WAECBiologyScraper) SubjectSlug() string { return "biology" }
func (s *WAECBiologyScraper) Level() string       { return "senior-secondary" }
func (s *WAECBiologyScraper) SourceURL() string {
	return "https://waecsyllabus.com/biology-syllabus/"
}

func (s *WAECBiologyScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Concept of Living and Organization of Life",
			Description: "Characteristics of living things, classification of 5 kingdoms, organization level, cell structure.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Characteristics of Living Things", Objectives: []string{
					"Compare living and non-living things using MR NIGER D characteristics",
					"Explain virus characteristics as intermediate between living and non-living",
				}},
				{Name: "Classification of Living Things", Objectives: []string{
					"Classify organisms into five kingdoms: Monera, Protista, Fungi, Plantae, Animalia",
					"Apply binomial nomenclature rules (genus and species naming)",
					"Identify key structural features of major plant phyla and animal phyla",
				}},
				{Name: "Organization of Life", Objectives: []string{
					"Trace organizational hierarchy: cell -> tissue -> organ -> system -> organism",
					"Analyze single-celled, colonial, and filament organisms (Amoeba, Volvox, Spirogyra)",
				}},
				{Name: "Cell Structure and Function", Objectives: []string{
					"Identify plant and animal cell organelles and their functions using light/electron microscopes",
					"Compare structural differences between plant and animal cells",
				}},
			},
		},
		{
			Name:        "Cellular and Organismal Life Processes",
			Description: "Diffusion, osmosis, autotrophic and heterotrophic nutrition, respiration, excretion, transport, and growth.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Cellular Transport Mechanisms", Objectives: []string{
					"Explain diffusion, osmosis, and active transport across cell membranes",
					"Demonstrate plasmolysis, turgidity, and osmoregulation in plant and animal cells",
				}},
				{Name: "Plant and Animal Nutrition", Objectives: []string{
					"Describe light and dark reactions of photosynthesis in chloroplasts",
					"Identify digestive system structures, enzymes, and nutrient absorption in mammals",
					"Perform chemical tests for carbohydrates, proteins, lipids, and vitamins",
				}},
				{Name: "Respiration and Gas Exchange", Objectives: []string{
					"Compare aerobic respiration and anaerobic respiration (glycolysis, Krebs cycle, fermentation)",
					"Describe respiratory surfaces and breathing mechanisms in humans, fish, insects, and plants",
				}},
				{Name: "Excretion and Osmoregulation", Objectives: []string{
					"Identify excretory organs and metabolic waste products in mammals, insects, and plants",
					"Explain nephron structure and urine formation (ultrafiltration and reabsorption) in mammalian kidneys",
				}},
				{Name: "Transport Systems in Animals and Plants", Objectives: []string{
					"Describe double circulation, blood composition, and heart structure in mammals",
					"Explain xylem water transport (transpiration pull, cohesion-tension) and phloem translocation",
				}},
				{Name: "Growth and Hormonal Regulation", Objectives: []string{
					"Measure plant and animal growth (sigmoid growth curve)",
					"Explain roles of plant hormones (auxins, gibberellins) and mammalian endocrine hormones",
				}},
			},
		},
		{
			Name:        "Reproduction and Heredity",
			Description: "Asexual and sexual reproduction, genetics, cell division (mitosis/meiosis), variation, and molecular genetics.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "Reproductive Systems and Types", Objectives: []string{
					"Compare asexual reproduction (binary fission, budding, vegetative propagation) and sexual reproduction",
					"Describe structure and function of male and female reproductive systems in mammals and flowering plants",
					"Explain pollination, fertilization, seed formation, and germination mechanisms",
				}},
				{Name: "Cell Division: Mitosis and Meiosis", Objectives: []string{
					"Describe stages of mitosis (prophase, metaphase, anaphase, telophase) and its significance in growth",
					"Describe stages of meiosis and its role in gamete formation and genetic variation",
				}},
				{Name: "Mendelian Genetics and Inheritance", Objectives: []string{
					"Apply Mendel's First and Second Laws of Inheritance using Punnett squares",
					"Calculate monohybrid cross and dihybrid cross genotypic and phenotypic ratios",
					"Explain sex determination, sex-linked traits (color blindness, hemophilia), and ABO blood grouping",
				}},
				{Name: "Variation and Genetic Applications", Objectives: []string{
					"Distinguish continuous variation (height, mass) and discontinuous variation (blood group, fingerprint)",
					"Explain DNA structure, gene mutations, chromosome mutations, and biotechnology applications (paternity tests, DNA fingerprinting)",
				}},
			},
		},
		{
			Name:        "Ecology, Adaptation, and Evolution",
			Description: "Ecosystem dynamics, nutrient cycles, food chains, ecological adaptation, conservation, and evolutionary theories.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Ecosystem Components and Energy Flow", Objectives: []string{
					"Identify biotic factors, abiotic factors, biomes, and habitats in terrestrial and aquatic ecosystems",
					"Construct food chains, food webs, and ecological pyramids (number, biomass, energy)",
					"Describe carbon cycle, nitrogen cycle, and water cycle in nature",
				}},
				{Name: "Ecological Adaptations and Succession", Objectives: []string{
					"Describe structural and behavioral adaptations of plants and animals for food, protection, and climate",
					"Explain primary and secondary ecological succession in habitats",
				}},
				{Name: "Pollution and Resource Conservation", Objectives: []string{
					"Identify causes, effects, and control of air, water, and soil pollution",
					"Discuss natural resource conservation, forest reserves, and climate change mitigation strategies",
				}},
				{Name: "Evolutionary Theories", Objectives: []string{
					"Compare Lamarck's theory of use and disuse and Darwin's theory of natural selection",
					"Identify evidence of evolution: fossils, comparative anatomy, embryology, and molecular biology",
				}},
			},
		},
	}
}
