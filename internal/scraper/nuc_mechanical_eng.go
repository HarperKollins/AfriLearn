package scraper

// NUCMechanicalEngScraper implements Scraper for NUC CCMAS B.Eng. Mechanical Engineering (100L - 500L)
type NUCMechanicalEngScraper struct{}

func NewNUCMechanicalEngScraper() *NUCMechanicalEngScraper {
	return &NUCMechanicalEngScraper{}
}

func (s *NUCMechanicalEngScraper) BoardSlug() string   { return "nuc" }
func (s *NUCMechanicalEngScraper) SubjectSlug() string { return "mechanical-engineering" }
func (s *NUCMechanicalEngScraper) Level() string       { return "tertiary-degree" }
func (s *NUCMechanicalEngScraper) SourceURL() string {
	return "https://nuc.edu.ng/ccmas/engineering/mechanical-engineering"
}

func (s *NUCMechanicalEngScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 - 200 Level: Engineering Mechanics and Thermodynamics",
			Description: "Applied Mechanics, Engineering Materials, Engineering Drawing, and Thermodynamics I.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "MEE 201: Engineering Thermodynamics I", Objectives: []string{
					"State 1st and 2nd Laws of Thermodynamics and apply energy balance to open and closed systems",
					"Analyze Carnot, Rankine, Otto, Diesel, and Brayton thermodynamic cycles",
				}},
				{Name: "MEE 202: Fluid Mechanics I and Strength of Materials", Objectives: []string{
					"Apply hydrostatics, Bernoulli's equation, and continuity equation to fluid flow",
					"Calculate shear force, bending moment, stress, strain, and torsion in structural elements",
				}},
			},
		},
		{
			Name:        "300 Level: Machine Design, Heat Transfer, Manufacturing, and SIWES",
			Description: "MEE 301 Machine Component Design, Heat Transfer, Manufacturing Technology, and SIWES.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "MEE 301: Machine Element Design", Objectives: []string{
					"Design gears, shafts, bearings, belts, clutches, brakes, and threaded fasteners",
				}},
				{Name: "MEE 302: Heat and Mass Transfer", Objectives: []string{
					"Calculate steady and unsteady conduction, free/forced convection, and radiation heat transfer",
					"Design heat exchangers using Log Mean Temperature Difference (LMTD) and NTU methods",
				}},
			},
		},
		{
			Name:        "400 - 500 Level: Control Engineering, Mechatronics, CAD/CAM, and B.Eng. Project",
			Description: "Internal Combustion Engines, Control Systems, CAD/CAM, Renewable Energy, and Thesis.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "MEE 501: CAD/CAM, Robotics, and Mechatronics", Objectives: []string{
					"Use 3D CAD modeling (SolidWorks/AutoCAD) and Computer Aided Manufacturing (CNC programming)",
					"Integrate hydraulic/pneumatic actuators, sensors, and microcontrollers for mechatronic automation",
				}},
				{Name: "MEE 599: Final Year Mechanical Engineering Design Project", Objectives: []string{
					"Design, fabricate, and test original mechanical machine prototype and defend before NSE/COREN panel",
				}},
			},
		},
	}
}
