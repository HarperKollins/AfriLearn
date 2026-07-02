package scraper

// UNILAGComputerScienceScraper implements Scraper for UNILAG B.Sc. Computer Science (Faculty of Science)
type UNILAGComputerScienceScraper struct{}

func NewUNILAGComputerScienceScraper() *UNILAGComputerScienceScraper {
	return &UNILAGComputerScienceScraper{}
}

func (s *UNILAGComputerScienceScraper) BoardSlug() string   { return "unilag" }
func (s *UNILAGComputerScienceScraper) SubjectSlug() string { return "computer-science" }
func (s *UNILAGComputerScienceScraper) Level() string       { return "tertiary-degree" }
func (s *UNILAGComputerScienceScraper) SourceURL() string {
	return "https://unilag.edu.ng/academics/science/computer-science"
}

func (s *UNILAGComputerScienceScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 - 200 Level: UNILAG Foundation Sciences & Programming",
			Description: "CSC 101, CSC 102, CSC 201 Object-Oriented Programming, MAT 101, and General Studies at UNILAG.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "CSC 101 & CSC 102: Introduction to Computer Science (UNILAG)", Objectives: []string{
					"Master computing fundamentals, Python programming, logic design, and Linux CLI at UNILAG Yaba",
				}},
				{Name: "CSC 201 & CSC 202: OOP with Java & Discrete Mathematics", Objectives: []string{
					"Implement data structures, OOP design patterns, Boolean algebra, graph theory, and automata",
				}},
			},
		},
		{
			Name:        "300 - 400 Level: UNILAG Systems, Cloud, AI, SIWES, and B.Sc. Thesis",
			Description: "Operating Systems, DBMS, Software Engineering, Cloud Computing, AI, SIWES, and UNILAG Thesis.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "CSC 301 & CSC 302: Operating Systems, Databases, and SIWES", Objectives: []string{
					"Analyze process scheduling, memory virtualization, SQL database design, and complete 6-month Yaba Tech Hub SIWES",
				}},
				{Name: "CSC 401 & CSC 499: Artificial Intelligence & B.Sc. Research Project", Objectives: []string{
					"Develop deep learning models, cloud backend systems, and defend UNILAG B.Sc. research project",
				}},
			},
		},
	}
}
