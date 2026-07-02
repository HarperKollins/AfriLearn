package scraper

// NUCElectricalEngScraper implements Scraper for NUC CCMAS B.Eng. Electrical & Electronic Engineering (100L - 500L)
type NUCElectricalEngScraper struct{}

func NewNUCElectricalEngScraper() *NUCElectricalEngScraper {
	return &NUCElectricalEngScraper{}
}

func (s *NUCElectricalEngScraper) BoardSlug() string   { return "nuc" }
func (s *NUCElectricalEngScraper) SubjectSlug() string { return "electrical-engineering" }
func (s *NUCElectricalEngScraper) Level() string       { return "tertiary-degree" }
func (s *NUCElectricalEngScraper) SourceURL() string {
	return "https://nuc.edu.ng/ccmas/engineering/electrical-engineering"
}

func (s *NUCElectricalEngScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 - 200 Level: Engineering Mathematics and Circuit Fundamentals",
			Description: "Engineering Mathematics, Applied Mechanics, Workshop Practice, and Circuit Theory I.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "ENG 101 & MTH 201: Engineering Mathematics and Mechanics", Objectives: []string{
					"Apply differential equations, vector calculus, linear algebra, and complex numbers to engineering problems",
					"Analyze statics, dynamics, stress-strain relations, and engineering workshop practice",
				}},
				{Name: "EEE 201: Circuit Theory I and Physical Electronics", Objectives: []string{
					"Apply Kirchhoff's laws, Thevenin's, Norton's, and Superposition theorems to DC and AC circuits",
					"Understand semiconductor physics, p-n junction diodes, and bipolar junction transistors (BJT)",
				}},
			},
		},
		{
			Name:        "300 Level: Electromagnetic Fields, Signals, Analog Electronics, and SIWES",
			Description: "EEE 301 Circuit Theory II, Electromagnetic Fields, Analog Electronics, Digital Systems, and SIWES.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "EEE 301 & EEE 302: Electromagnetic Fields and Signal Analysis", Objectives: []string{
					"State Maxwell's equations and analyze electromagnetic wave propagation in waveguides and transmission lines",
					"Apply Laplace transforms, Fourier series, and Z-transforms to continuous and discrete signals",
				}},
				{Name: "EEE 303 & EEE 304: Analog Electronics and Microprocessors", Objectives: []string{
					"Design operational amplifier (Op-Amp) circuits, active filters, and power amplifiers",
					"Program 8086/ARM microprocessors and microcontrollers (Arduino/PIC) for embedded control",
				}},
			},
		},
		{
			Name:        "400 - 500 Level: Power Systems, Control Engineering, Telecommunications, and B.Eng. Project",
			Description: "Power Systems Analysis, Control Theory, Telecommunications, High Voltage, and Final Project.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "EEE 401 & EEE 402: Power Systems and Control Engineering", Objectives: []string{
					"Analyze power generation, transmission line parameters, fault calculations, and grid stability",
					"Apply root locus, Bode plots, and state-space analysis to feedback control systems",
				}},
				{Name: "EEE 501: Telecommunications Engineering and Power Electronics", Objectives: []string{
					"Analyze AM, FM, digital modulations (QPSK, QAM), cellular networks (4G/5G), and fiber optics",
					"Design power electronic converters (rectifiers, inverters, choppers) for renewable energy integration",
				}},
				{Name: "EEE 599: Final Year Engineering Design Project", Objectives: []string{
					"Design, fabricate, test, and document an original electrical/electronic engineering device or system",
					"Defend design project before COREN and NSE professional engineering accreditation panels",
				}},
			},
		},
	}
}
