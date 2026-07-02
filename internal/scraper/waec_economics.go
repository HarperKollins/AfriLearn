package scraper

// WAECEconomicsScraper implements Scraper for WAEC Economics (SS1 - SS3)
type WAECEconomicsScraper struct{}

func NewWAECEconomicsScraper() *WAECEconomicsScraper {
	return &WAECEconomicsScraper{}
}

func (s *WAECEconomicsScraper) BoardSlug() string   { return "waec" }
func (s *WAECEconomicsScraper) SubjectSlug() string { return "economics" }
func (s *WAECEconomicsScraper) Level() string       { return "senior-secondary" }
func (s *WAECEconomicsScraper) SourceURL() string {
	return "https://waecsyllabus.com/economics-syllabus/"
}

func (s *WAECEconomicsScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Basic Economic Concepts and Production",
			Description: "Scarcity, choice, scale of preference, opportunity cost, production, division of labor, and business units.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Fundamental Economic Concepts", Objectives: []string{
					"Explain basic economic concepts: wants, scarcity, choice, scale of preference, opportunity cost",
					"Distinguish microeconomics and macroeconomics",
				}},
				{Name: "Production and Factors of Production", Objectives: []string{
					"Describe primary, secondary, and tertiary production stages",
					"Analyze land, labor, capital, and entrepreneurship as factors of production and their rewards",
					"Explain division of labor, specialization, and localization of industries",
				}},
				{Name: "Business Organizations", Objectives: []string{
					"Compare sole proprietorships, partnerships, joint-stock companies, co-operatives, and public enterprises",
				}},
			},
		},
		{
			Name:        "Microeconomics: Demand, Supply, and Markets",
			Description: "Laws of demand and supply, price equilibrium, elasticity, consumer behavior, and market structures.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Demand and Supply", Objectives: []string{
					"State laws of demand and supply and construct demand and supply schedules/curves",
					"Calculate equilibrium price and equilibrium quantity algebraically and graphically",
					"Distinguish changes in demand/supply vs. changes in quantity demanded/supplied",
				}},
				{Name: "Elasticity of Demand and Supply", Objectives: []string{
					"Calculate Price Elasticity of Demand (PED), Income Elasticity (YED), and Cross Elasticity (XED)",
					"Calculate Price Elasticity of Supply (PES) and explain factors affecting elasticity",
				}},
				{Name: "Market Structures", Objectives: []string{
					"Compare perfect competition, monopoly, monopolistic competition, and oligopoly price/output determination",
				}},
			},
		},
		{
			Name:        "Macroeconomics, Money, and Public Finance",
			Description: "National income accounting, money, banking, inflation, public finance, taxation, and fiscal policy.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "National Income Accounting", Objectives: []string{
					"Define Gross Domestic Product (GDP), Gross National Product (GNP), Net National Product (NNP), and Per Capita Income",
					"Calculate national income using Output, Income, and Expenditure methods",
				}},
				{Name: "Money, Banking, and Inflation", Objectives: []string{
					"Explain functions of money, commercial banks, and central bank monetary policy tools",
					"Analyze causes, effects, and control of inflation and deflation",
				}},
				{Name: "Public Finance and Taxation", Objectives: []string{
					"Explain government revenue sources, recurrent/capital expenditure, and budget types (balanced, surplus, deficit)",
					"Distinguish direct taxation and indirect taxation and principles of taxation (canons of taxation)",
				}},
			},
		},
	}
}
