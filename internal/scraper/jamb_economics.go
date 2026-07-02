package scraper

// JAMBEconomicsScraper implements Scraper for JAMB Economics
type JAMBEconomicsScraper struct{}

func NewJAMBEconomicsScraper() *JAMBEconomicsScraper {
	return &JAMBEconomicsScraper{}
}

func (s *JAMBEconomicsScraper) BoardSlug() string   { return "jamb" }
func (s *JAMBEconomicsScraper) SubjectSlug() string { return "economics" }
func (s *JAMBEconomicsScraper) Level() string       { return "tertiary-entry" }
func (s *JAMBEconomicsScraper) SourceURL() string {
	return "https://ibass.jamb.gov.ng/syllabus/economics"
}

func (s *JAMBEconomicsScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Economics Fundamentals and Production",
			Description: "JAMB UTME Section I: Economic methodology, scarcity, choice, opportunity cost, and production systems.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Economic Concepts and Methodology", Objectives: []string{
					"Explain basic economic tools: tables, graphs, charts, measures of central tendency",
					"Analyze scarcity, choice, scale of preference, opportunity cost, and economic systems (capitalist, socialist, mixed)",
				}},
				{Name: "Production and Business Units", Objectives: []string{
					"Explain factors of production, division of labor, specialization, and production functions",
					"Compare sole proprietorships, partnerships, joint-stock companies, and public corporations",
				}},
			},
		},
		{
			Name:        "Microeconomics: Demand, Supply, Prices, and Markets",
			Description: "JAMB UTME Section II: Demand, supply, price determination, elasticities, utility, and market structures.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Demand, Supply, and Price Equilibrium", Objectives: []string{
					"Calculate market equilibrium price and quantity algebraically and graphically",
					"Calculate Price Elasticity of Demand (PED), Supply (PES), Income (YED), and Cross Elasticity (XED)",
				}},
				{Name: "Utility Theory and Market Structures", Objectives: []string{
					"Apply Cardinal Utility (Law of Diminishing Marginal Utility) and Ordinal Utility (indifference curves)",
					"Compare price and output determination in Perfect Competition, Monopoly, Monopolistic Competition, and Oligopoly",
				}},
			},
		},
		{
			Name:        "Macroeconomics: Money, Public Finance, and International Trade",
			Description: "JAMB UTME Section III & IV: National income, inflation, money, banking, public finance, population, and trade.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "Money, Banking, and Inflation", Objectives: []string{
					"Explain functions of money, commercial bank credit creation, and central bank monetary policy tools",
					"Analyze causes, effects, and fiscal/monetary solutions for inflation and deflation",
				}},
				{Name: "National Income, Public Finance, and Trade", Objectives: []string{
					"Calculate GDP, GNP, NNP, and Per Capita Income using Output, Income, and Expenditure approaches",
					"Analyze government budget, taxation canons, balance of payments, and international trade barriers/agreements",
				}},
			},
		},
	}
}
