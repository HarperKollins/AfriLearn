package scraper

// EBSUComputerScienceScraper implements Scraper for EBSU B.Sc. Computer Science (Faculty of Physical Sciences)
type EBSUComputerScienceScraper struct{}

func NewEBSUComputerScienceScraper() *EBSUComputerScienceScraper {
	return &EBSUComputerScienceScraper{}
}

func (s *EBSUComputerScienceScraper) BoardSlug() string   { return "ebsu" }
func (s *EBSUComputerScienceScraper) SubjectSlug() string { return "computer-science" }
func (s *EBSUComputerScienceScraper) Level() string       { return "tertiary-degree" }
func (s *EBSUComputerScienceScraper) SourceURL() string {
	return "https://ebsu.edu.ng/academics/physical-sciences/computer-science"
}

func (s *EBSUComputerScienceScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 Level: EBSU Science & Computing Foundation",
			Description: "CSC 101, CSC 102, MAT 101, PHY 101, and GST 101 courses at EBSU Abakaliki.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "CSC 101: Introduction to Computer Science (EBSU)", Objectives: []string{
					"Understand history of computing, computer hardware architectures, and operating systems",
					"Solve basic mathematical and logical problems using Python algorithms",
				}},
				{Name: "MAT 101 & PHY 101: Algebra and Physics for Computing", Objectives: []string{
					"Apply set theory, matrices, calculus, and electricity/magnetism in computer science",
				}},
			},
		},
		{
			Name:        "200 Level: EBSU Core Programming & Digital Logic",
			Description: "CSC 201 Object-Oriented Programming, Discrete Structures, and Data Structures I.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "CSC 201: OOP using C++ / Java (EBSU)", Objectives: []string{
					"Implement OOP encapsulation, inheritance, polymorphism, and exception handling in Java",
				}},
				{Name: "CSC 202: Data Structures & Assembly Language", Objectives: []string{
					"Implement arrays, linked lists, stacks, queues, and 8086 assembly language",
				}},
			},
		},
		{
			Name:        "300 - 400 Level: EBSU Advanced Computing, SIWES, and B.Sc. Project",
			Description: "DBMS, Software Engineering, Artificial Intelligence, SIWES, and EBSU Thesis.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "CSC 301 & CSC 302: Database Systems and SIWES", Objectives: []string{
					"Design relational SQL databases using ER diagrams and normalization",
					"Complete 6 months industrial training at tech hubs/organizations and present defense",
				}},
				{Name: "CSC 401 & CSC 499: Artificial Intelligence and B.Sc. Research Project", Objectives: []string{
					"Implement machine learning algorithms, computer vision, and defend original EBSU thesis",
				}},
			},
		},
	}
}
