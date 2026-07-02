package scraper

// YABATECHComputerEngScraper implements Scraper for YABATECH ND & HND Computer Engineering Tech (NBTE Standard)
type YABATECHComputerEngScraper struct{}

func NewYABATECHComputerEngScraper() *YABATECHComputerEngScraper {
	return &YABATECHComputerEngScraper{}
}

func (s *YABATECHComputerEngScraper) BoardSlug() string   { return "yabatech" }
func (s *YABATECHComputerEngScraper) SubjectSlug() string { return "computer-engineering-tech" }
func (s *YABATECHComputerEngScraper) Level() string       { return "polytechnic-nd" }
func (s *YABATECHComputerEngScraper) SourceURL() string {
	return "https://yabatech.edu.ng/academics/engineering/computer-engineering"
}

func (s *YABATECHComputerEngScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "ND I & ND II: YABATECH Computer Hardware, Electronics, and Workshop",
			Description: "CTE 111, CTE 121 Digital Logic, Electronic Devices, Computer Hardware Maintenance, and Workshop Practice.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "CTE 111 & CTE 121: Digital Electronics and Logic Design (YABATECH)", Objectives: []string{
					"Study logic gates, flip-flops, counters, registers, and Boolean algebraic simplification",
					"Assemble and troubleshoot digital logic circuits on breadboards in YABATECH labs",
				}},
				{Name: "CTE 211: Computer Hardware System Maintenance & Assembly", Objectives: []string{
					"Diagnose PC motherboard components, power supply units, RAM, storage, and peripheral installation",
					"Perform computer system troubleshooting, OS formatting, and network cabling (RJ45 crimping)",
				}},
			},
		},
		{
			Name:        "HND I & HND II: Microprocessors, Embedded Systems, and YABATECH Project",
			Description: "Microcomputer Systems, Embedded C/Assembly, Interfacing, and Final HND Project.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "CTE 311 & CTE 411: Microprocessor Architecture & Embedded Systems", Objectives: []string{
					"Program 8086/ARM microprocessors and microcontrollers (Arduino/PIC) for industrial automation",
					"Interface analog-to-digital converters (ADC), sensors, motors, and LCD displays",
				}},
				{Name: "CTE 499: Final Year HND Hardware Engineering Project", Objectives: []string{
					"Design, solder, test, and defend hardware project prototype before NBTE and YABATECH academic board",
				}},
			},
		},
	}
}
