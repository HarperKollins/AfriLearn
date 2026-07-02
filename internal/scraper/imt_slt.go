package scraper

// IMTSLTScraper implements Scraper for IMT Enugu ND & HND Science Laboratory Technology (NBTE Standard)
type IMTSLTScraper struct{}

func NewIMTSLTScraper() *IMTSLTScraper {
	return &IMTSLTScraper{}
}

func (s *IMTSLTScraper) BoardSlug() string   { return "imt" }
func (s *IMTSLTScraper) SubjectSlug() string { return "science-laboratory-tech" }
func (s *IMTSLTScraper) Level() string       { return "polytechnic-nd" }
func (s *IMTSLTScraper) SourceURL() string {
	return "https://imt.edu.ng/academics/applied-sciences/science-laboratory-technology"
}

func (s *IMTSLTScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "ND I & ND II: IMT Laboratory Techniques, Chemistry, and Biology Options",
			Description: "SLT 111 General Laboratory Techniques, Analytical Chemistry, Microbiology, and Physics Techniques.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "SLT 111 & SLT 121: General Laboratory Management & Safety (IMT)", Objectives: []string{
					"Understand laboratory safety regulations, chemical storage, sterilization, and equipment calibration",
					"Perform titrimetric, gravimetric, and spectrophotometric laboratory analyses",
				}},
				{Name: "SLT 211: Biological and Chemical Laboratory Techniques", Objectives: []string{
					"Prepare culture media, slide staining, microscopy, and chemical reagent standardization",
				}},
			},
		},
		{
			Name:        "HND I & HND II: Advanced Instrumentation, Biochemistry Option, and IMT Project",
			Description: "Chromatography, Spectroscopy, Clinical Biochemistry, Industrial Microbiology, and HND Thesis.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "SLT 311 & SLT 411: Advanced Analytical Instrumentation", Objectives: []string{
					"Operate HPLC, Gas Chromatography (GC), Atomic Absorption Spectroscopy (AAS), and UV-Vis spectrophotometers",
				}},
				{Name: "SLT 499: Final HND Science Laboratory Research Project", Objectives: []string{
					"Conduct independent laboratory research and defend thesis before IMT Applied Sciences board",
				}},
			},
		},
	}
}
