package scraper

// BECESocialStudiesScraper implements Scraper for BECE Social Studies (JSS1 - JSS3)
type BECESocialStudiesScraper struct{}

func NewBECESocialStudiesScraper() *BECESocialStudiesScraper {
	return &BECESocialStudiesScraper{}
}

func (s *BECESocialStudiesScraper) BoardSlug() string   { return "bece" }
func (s *BECESocialStudiesScraper) SubjectSlug() string { return "social-studies" }
func (s *BECESocialStudiesScraper) Level() string       { return "junior-secondary" }
func (s *BECESocialStudiesScraper) SourceURL() string {
	return "https://nerdc.gov.ng/curriculum/junior-secondary/social-studies"
}

func (s *BECESocialStudiesScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Man, Culture, and Social Environment (JSS1 - JSS3)",
			Description: "Physical and social environment, culture, socialization, family types, marriage, and leadership.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Social Environment, Family, and Culture", Objectives: []string{
					"Explain man's interaction with physical and social environments",
					"Compare nuclear family and extended family roles and responsibilities",
					"Describe cultural elements, diversity, preservation, and cultural changes in Nigeria",
				}},
				{Name: "Socialization, Marriage, and Gender Roles", Objectives: []string{
					"Explain agents of socialization: family, school, peer group, religious institutions, mass media",
					"Describe types of marriage (monogamy, polygamy) and conditions for successful marriage",
				}},
			},
		},
		{
			Name:        "Social Issues, Safety, and National Unity (JSS1 - JSS3)",
			Description: "Social problems, crime, road safety, national unity, integration, and conflict resolution.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Contemporary Social Problems and Safety", Objectives: []string{
					"Identify causes, consequences, and solutions for cultism, examination malpractice, juvenile delinquency",
					"Demonstrate road safety rules, traffic signs, and roles of FRSC in accident prevention",
				}},
				{Name: "National Identity, Integration, and Conflict", Objectives: []string{
					"Identify national symbols: coat of arms, national flag, national anthem, pledge, currency",
					"Explain causes of communal/ethnic conflicts and peaceful conflict resolution mechanisms",
				}},
			},
		},
		{
			Name:        "Civic Education, Governance, and Rights (JSS1 - JSS3)",
			Description: "Civic values, citizenship, human rights, constitution, democracy, and governance in Nigeria.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Civic Values and Citizenship", Objectives: []string{
					"Define civic values: honesty, discipline, self-reliance, integrity, cooperation",
					"Explain ways of acquiring Nigerian citizenship and duties/responsibilities of citizens",
				}},
				{Name: "Human Rights, Constitution, and Governance", Objectives: []string{
					"Identify fundamental human rights and agencies protecting human rights in Nigeria",
					"Explain democratic principles, free elections, rule of law, and arms of government (executive, legislature, judiciary)",
				}},
			},
		},
	}
}
