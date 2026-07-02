package scraper

// FUNAIComputerScienceScraper implements Scraper for AE-FUNAI B.Sc. Computer Science (Faculty of Science)
type FUNAIComputerScienceScraper struct{}

func NewFUNAIComputerScienceScraper() *FUNAIComputerScienceScraper {
	return &FUNAIComputerScienceScraper{}
}

func (s *FUNAIComputerScienceScraper) BoardSlug() string   { return "funai" }
func (s *FUNAIComputerScienceScraper) SubjectSlug() string { return "computer-science" }
func (s *FUNAIComputerScienceScraper) Level() string       { return "tertiary-degree" }
func (s *FUNAIComputerScienceScraper) SourceURL() string {
	return "https://funai.edu.ng/academics/science/computer-science"
}

func (s *FUNAIComputerScienceScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 - 200 Level: AE-FUNAI Science Foundation & Programming",
			Description: "COS 101, COS 102, CSC 201, MTH 101, and General Studies at AE-FUNAI Ikwo.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "COS 101 & COS 102: Introduction to Computing Sciences (FUNAI)", Objectives: []string{
					"Understand computing history, operating systems, and basic problem solving in Python",
				}},
				{Name: "CSC 201 & CSC 202: Object-Oriented Programming & Discrete Structures", Objectives: []string{
					"Master OOP principles in Java/C++, discrete mathematics, logic gates, and Boolean algebra",
				}},
			},
		},
		{
			Name:        "300 - 400 Level: Advanced Systems, Cybersecurity, SIWES, and Project",
			Description: "Data Structures, Database Systems, Computer Networks, Cybersecurity, and FUNAI Thesis.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "CSC 301 & CSC 302: Data Structures and Industrial Training", Objectives: []string{
					"Implement advanced trees, graphs, sorting algorithms, and complete 6-month SIWES internship",
				}},
				{Name: "CSC 401 & CSC 499: Cybersecurity, Machine Learning, and Research Project", Objectives: []string{
					"Study network security, cryptography, machine learning, and defend final year FUNAI B.Sc. thesis",
				}},
			},
		},
	}
}
