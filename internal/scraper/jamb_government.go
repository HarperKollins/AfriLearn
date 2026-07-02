package scraper

// JAMBGovernmentScraper implements Scraper for JAMB Government
type JAMBGovernmentScraper struct{}

func NewJAMBGovernmentScraper() *JAMBGovernmentScraper {
	return &JAMBGovernmentScraper{}
}

func (s *JAMBGovernmentScraper) BoardSlug() string   { return "jamb" }
func (s *JAMBGovernmentScraper) SubjectSlug() string { return "government" }
func (s *JAMBGovernmentScraper) Level() string       { return "tertiary-entry" }
func (s *JAMBGovernmentScraper) SourceURL() string {
	return "https://ibass.jamb.gov.ng/syllabus/government"
}

func (s *JAMBGovernmentScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Political Theory and Concepts of Government",
			Description: "JAMB UTME Section I: Definition of government, state, sovereignty, power, legitimacy, Rule of Law, and systems.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Basic Political Concepts", Objectives: []string{
					"Analyze government as an institution, process, and field of study",
					"Distinguish state, nation, sovereignty, power, authority, legitimacy, political culture, and political socialization",
				}},
				{Name: "Systems of Government and Ideologies", Objectives: []string{
					"Compare presidential, parliamentary, unitary, federal, confederal, monarchical, and republican systems",
					"Distinguish capitalism, socialism, communism, fascism, and feudalism",
				}},
			},
		},
		{
			Name:        "Organs of Government and Political Processes",
			Description: "JAMB UTME Section II: Legislature, Executive, Judiciary, constitutions, political parties, and elections.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Organs of Government and Rule of Law", Objectives: []string{
					"Describe composition, functions, and relationships of Legislature, Executive, and Judiciary",
					"Analyze Separation of Powers, Checks and Balances, Rule of Law, and Constitutionalism",
				}},
				{Name: "Elections, Political Parties, and Pressure Groups", Objectives: []string{
					"Compare electoral systems: First-Past-The-Post, Proportional Representation, Absolute Majority",
					"Describe political party systems, pressure groups, public opinion, and mass media roles in governance",
				}},
			},
		},
		{
			Name:        "Constitutional History and Foreign Policy of Nigeria",
			Description: "JAMB UTME Section III & IV: Pre-colonial systems, colonial administration, constitutions, federalism, and foreign policy.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "Pre-Colonial Systems and Constitutional History", Objectives: []string{
					"Compare pre-colonial political systems in Hausa/Fulani, Yoruba, and Igbo societies",
					"Analyze Clifford (1922), Richards (1946), Macpherson (1951), Lyttelton (1954), 1960, 1963, 1979, and 1999 Constitutions",
				}},
				{Name: "Nigerian Federalism and Foreign Policy", Objectives: []string{
					"Analyze state creation, revenue allocation, local government reform, and civil service in Nigeria",
					"Describe foreign policy principles, Afrocentricity, non-alignment, and membership in UN, AU, ECOWAS, Commonwealth",
				}},
			},
		},
	}
}
