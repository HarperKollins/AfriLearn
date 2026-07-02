package scraper

// UNNComputerScienceScraper implements Scraper for UNN B.Sc. Computer Science (Faculty of Physical Sciences, Nsukka)
type UNNComputerScienceScraper struct{}

func NewUNNComputerScienceScraper() *UNNComputerScienceScraper {
	return &UNNComputerScienceScraper{}
}

func (s *UNNComputerScienceScraper) BoardSlug() string   { return "unn" }
func (s *UNNComputerScienceScraper) SubjectSlug() string { return "computer-science" }
func (s *UNNComputerScienceScraper) Level() string       { return "tertiary-degree" }
func (s *UNNComputerScienceScraper) SourceURL() string {
	return "https://unn.edu.ng/academics/faculties/physical-sciences/computer-science"
}

func (s *UNNComputerScienceScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 - 200 Level: UNN Science Foundation & Core Programming",
			Description: "COS 101, COS 102, CSC 201, MTH 101, PHY 101, and General Studies at UNN Nsukka.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "COS 101 & COS 102: Computing Fundamentals & Problem Solving (UNN)", Objectives: []string{
					"Study computer architectures, algorithms, flowcharts, and Python programming at UNN",
				}},
				{Name: "CSC 201 & CSC 202: OOP with Java and Digital Logic Design", Objectives: []string{
					"Master Java object-oriented programming, digital circuits, combinational logic, and assembly",
				}},
			},
		},
		{
			Name:        "300 - 400 Level: Advanced Computing, Software Engineering, SIWES, and UNN Project",
			Description: "Data Structures, Database Management, Software Engineering, AI, and UNN Thesis.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "CSC 301 & CSC 302: Data Structures, DBMS, and SIWES Internship", Objectives: []string{
					"Implement data structures, SQL database modeling, and complete 6-month industrial training",
				}},
				{Name: "CSC 401 & CSC 499: Artificial Intelligence, Compiler Construction, and Thesis", Objectives: []string{
					"Implement neural networks, compiler parsing, and defend original B.Sc. research thesis at UNN",
				}},
			},
		},
	}
}
