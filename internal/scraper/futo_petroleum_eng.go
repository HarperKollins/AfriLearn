package scraper

// FUTOPetroleumEngScraper implements Scraper for FUTO B.Eng. Petroleum Engineering (School of Engineering, Owerri)
type FUTOPetroleumEngScraper struct{}

func NewFUTOPetroleumEngScraper() *FUTOPetroleumEngScraper {
	return &FUTOPetroleumEngScraper{}
}

func (s *FUTOPetroleumEngScraper) BoardSlug() string   { return "futo" }
func (s *FUTOPetroleumEngScraper) SubjectSlug() string { return "petroleum-engineering" }
func (s *FUTOPetroleumEngScraper) Level() string       { return "tertiary-degree" }
func (s *FUTOPetroleumEngScraper) SourceURL() string {
	return "https://futo.edu.ng/academics/seet/petroleum-engineering"
}

func (s *FUTOPetroleumEngScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 - 200 Level: FUTO Engineering Foundation & Geology",
			Description: "ENG 101, ENG 201 Thermodynamics, Fluid Mechanics, Geology for Petroleum Engineers at FUTO.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "PET 201: Introduction to Petroleum Engineering & Geology (FUTO)", Objectives: []string{
					"Understand origin of hydrocarbons, sedimentary basins, petroleum reservoir rocks, and fluid dynamics",
				}},
			},
		},
		{
			Name:        "300 - 500 Level: Reservoir, Drilling, Production, SIWES, and FUTO B.Eng. Project",
			Description: "Reservoir Engineering, Drilling Hydraulics, Well Logging, Enhanced Oil Recovery, SIWES, and Thesis.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "PET 301 & PET 401: Drilling Technology, Reservoir Simulation, and Well Logging", Objectives: []string{
					"Design rotary drilling bits, mud hydraulics, pressure control, and reservoir fluid phase behavior",
					"Interpret wireline logs (Gamma Ray, Resistivity, Density, Neutron) for hydrocarbon pay zones",
				}},
				{Name: "PET 501 & PET 599: Production Engineering, Gas Technology, and B.Eng. Project", Objectives: []string{
					"Design artificial lift systems (gas lift, ESP), natural gas processing, and defend FUTO B.Eng. thesis",
				}},
			},
		},
	}
}
