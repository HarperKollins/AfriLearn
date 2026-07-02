package scraper

// WAECGovernmentScraper implements Scraper for WAEC Government (SS1 - SS3)
type WAECGovernmentScraper struct{}

func NewWAECGovernmentScraper() *WAECGovernmentScraper {
	return &WAECGovernmentScraper{}
}

func (s *WAECGovernmentScraper) BoardSlug() string   { return "waec" }
func (s *WAECGovernmentScraper) SubjectSlug() string { return "government" }
func (s *WAECGovernmentScraper) Level() string       { return "senior-secondary" }
func (s *WAECGovernmentScraper) SourceURL() string {
	return "https://waecsyllabus.com/government-syllabus/"
}

func (s *WAECGovernmentScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Basic Concepts and Political Ideologies",
			Description: "Definition of government, state, nation, sovereignty, power, authority, legitimacy, democracy, capitalism, socialism.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Fundamental Concepts of Government", Objectives: []string{
					"Define government as an institution of state, process of governing, and academic field",
					"Distinguish state, nation, society, and government",
					"Analyze concepts of sovereignty, power, authority, legitimacy, and political culture",
				}},
				{Name: "Political Ideologies and Systems of Government", Objectives: []string{
					"Compare capitalism, socialism, communism, communalism, fascism, and feudalism",
					"Compare presidential system, parliamentary system, unitary system, federal system, and confederal system",
				}},
			},
		},
		{
			Name:        "Organs of Government and Rule of Law",
			Description: "Legislature, Executive, Judiciary, Separation of Powers, Checks and Balances, Rule of Law, and Citizenship.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Organs of Government", Objectives: []string{
					"Describe composition, functions, and types (unicameral, bicameral) of Legislature",
					"Describe composition, types, and functions of Executive and Judiciary",
					"Analyze Separation of Powers and Checks and Balances in governance",
				}},
				{Name: "Rule of Law, Constitutionalism, and Citizenship", Objectives: []string{
					"Explain principles of Rule of Law and factors limiting Rule of Law",
					"Describe types and characteristics of Constitutions (written, unwritten, rigid, flexible)",
					"Explain rights, duties, and obligations of citizens",
				}},
			},
		},
		{
			Name:        "Constitutional Development and Politics in Nigeria",
			Description: "Pre-colonial political systems, colonial administration, nationalist movements, constitutional history (1922-1999), and Federalism.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "Pre-Colonial Administration and Colonial Rule", Objectives: []string{
					"Compare Hausa/Fulani emirate system, Yoruba Oyo kingdom system, and Igbo decentralized political system",
					"Analyze British Indirect Rule system in Northern, Western, and Eastern Nigeria",
				}},
				{Name: "Constitutional History and Nigerian Federalism", Objectives: []string{
					"Analyze Clifford (1922), Richards (1946), Macpherson (1951), Lyttelton (1954), Independence (1960), Republican (1963), and 1979/1999 Constitutions",
					"Trace development of Nigerian Federalism, state creation, and revenue allocation formulas",
				}},
			},
		},
	}
}
