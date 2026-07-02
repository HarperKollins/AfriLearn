package scraper

// NUCComputerScienceScraper implements Scraper for NUC CCMAS B.Sc. Computer Science (100L - 400L)
type NUCComputerScienceScraper struct{}

func NewNUCComputerScienceScraper() *NUCComputerScienceScraper {
	return &NUCComputerScienceScraper{}
}

func (s *NUCComputerScienceScraper) BoardSlug() string   { return "nuc" }
func (s *NUCComputerScienceScraper) SubjectSlug() string { return "computer-science" }
func (s *NUCComputerScienceScraper) Level() string       { return "tertiary-degree" }
func (s *NUCComputerScienceScraper) SourceURL() string {
	return "https://nuc.edu.ng/ccmas/computing/computer-science"
}

func (s *NUCComputerScienceScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 Level: Foundation Computing and Natural Sciences",
			Description: "COS 101, COS 102, General Studies, Mathematics, and Physics foundation courses.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "COS 101: Introduction to Computing Sciences", Objectives: []string{
					"Trace history of computing, computer generations, and classification of computers",
					"Identify computer hardware components, central processing unit CPU, and peripheral devices",
					"Distinguish system software, operating systems, and application software",
				}},
				{Name: "COS 102: Problem Solving and Programming I", Objectives: []string{
					"Formulate algorithms using flowcharts and pseudocode",
					"Understand variables, data types, control structures (if-else, loops), and functions",
					"Write basic programs in Python/C++ to solve computational problems",
				}},
				{Name: "MTH 101 & PHY 101: Mathematics and Physics for Computing", Objectives: []string{
					"Apply algebra, trigonometry, coordinate geometry, and differential calculus",
					"Understand mechanics, heat, electricity, magnetism, and wave motion principles",
				}},
			},
		},
		{
			Name:        "200 Level: Core Computer Science Fundamentals",
			Description: "CSC 201 Programming II, Discrete Structures, Digital Logic, and Computer Architecture.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "CSC 201: Computer Programming II", Objectives: []string{
					"Implement object-oriented programming concepts: classes, objects, encapsulation, inheritance, polymorphism",
					"Handle file I/O operations, exceptions, and recursion in C++/Java",
				}},
				{Name: "CSC 203: Discrete Structures and Logic", Objectives: []string{
					"Apply set theory, mathematical logic, truth tables, and proof techniques",
					"Analyze relations, functions, graph theory, trees, and combinatorics in computing",
				}},
				{Name: "CSC 204: Digital Logic Design and Assembly Language", Objectives: []string{
					"Design logic gates (AND, OR, NOT, NAND, NOR, XOR) and Boolean algebra simplification",
					"Construct combinational circuits (adders, multiplexers) and sequential circuits (flip-flops, counters)",
					"Understand computer organization, instruction set architecture (ISA), and assembly language programming",
				}},
			},
		},
		{
			Name:        "300 Level: Data Structures, Software Engineering, and SIWES",
			Description: "CSC 301 Data Structures, Database Systems, Software Engineering, Operating Systems, and SIWES Training.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "CSC 301: Data Structures and Algorithms", Objectives: []string{
					"Implement arrays, linked lists, stacks, queues, trees, binary search trees, and graphs",
					"Analyze algorithm complexity using Big-O notation for searching (binary search) and sorting (quicksort, mergesort)",
				}},
				{Name: "CSC 302: Database Management Systems (DBMS)", Objectives: []string{
					"Design relational databases using Entity-Relationship (ER) modeling and 1NF, 2NF, 3NF normalization",
					"Write SQL queries for Data Definition Language (DDL) and Data Manipulation Language (DML)",
					"Explain transaction processing, ACID properties, indexing, and database security",
				}},
				{Name: "CSC 304: Software Engineering & Operating Systems", Objectives: []string{
					"Apply Agile and Waterfall software development life cycle (SDLC) methodologies",
					"Understand process scheduling, memory management (virtual memory, paging), file systems, and concurrency/deadlocks",
				}},
				{Name: "CSC 399: SIWES (Students Industrial Work Experience Scheme)", Objectives: []string{
					"Gain 6 months real-world software development and IT industry internship experience",
					"Write technical industrial training report and defend work before departmental board",
				}},
			},
		},
		{
			Name:        "400 Level: Advanced Specializations and Final Year Project",
			Description: "Artificial Intelligence, Computer Networks, Cybersecurity, Compiler Construction, and Research Project.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "CSC 401: Artificial Intelligence and Machine Learning", Objectives: []string{
					"Implement search algorithms (A* search, minimax) and knowledge representation",
					"Understand supervised learning, unsupervised learning, neural networks, and deep learning basics",
				}},
				{Name: "CSC 402: Computer Networks and Cybersecurity", Objectives: []string{
					"Explain OSI 7-layer model and TCP/IP protocol suite (IP routing, DNS, HTTP, TLS)",
					"Understand network security threats, symmetric/asymmetrical cryptography (RSA, AES), and firewalls",
				}},
				{Name: "CSC 499: Final Year Research Project", Objectives: []string{
					"Conduct independent original software development or theoretical computer science research",
					"Write formal undergraduate thesis and present defense before external examiner",
				}},
			},
		},
	}
}
