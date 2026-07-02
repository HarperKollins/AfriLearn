package scraper

// NUCBusinessAdminScraper implements Scraper for NUC CCMAS B.Sc. Business Administration (100L - 400L)
type NUCBusinessAdminScraper struct{}

func NewNUCBusinessAdminScraper() *NUCBusinessAdminScraper {
	return &NUCBusinessAdminScraper{}
}

func (s *NUCBusinessAdminScraper) BoardSlug() string   { return "nuc" }
func (s *NUCBusinessAdminScraper) SubjectSlug() string { return "business-administration" }
func (s *NUCBusinessAdminScraper) Level() string       { return "tertiary-degree" }
func (s *NUCBusinessAdminScraper) SourceURL() string {
	return "https://nuc.edu.ng/ccmas/administration/business-administration"
}

func (s *NUCBusinessAdminScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 - 200 Level: Principles of Management and Business Environment",
			Description: "BUS 101 Introduction to Business, BUS 201 Principles of Management, and Business Environment.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "BUS 101 & BUS 201: Fundamentals of Business and Management", Objectives: []string{
					"Understand evolution of management thought (Taylor, Fayol, Weber) and management functions",
					"Analyze political, economic, socio-cultural, and technological (PEST) business environment in Nigeria",
				}},
			},
		},
		{
			Name:        "300 Level: Organizational Behavior, Marketing, and Human Resource Management",
			Description: "BUS 301 Organizational Behavior, Marketing Management, and Human Resource Management.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "BUS 301 & MKT 301: Organizational Behavior and Marketing", Objectives: []string{
					"Analyze individual behavior, group dynamics, motivation theories (Maslow, Herzberg), and leadership styles",
					"Apply 4 Ps of marketing mix (Product, Price, Place, Promotion) and market segmentation",
				}},
				{Name: "HRM 301: Human Resource Management", Objectives: []string{
					"Execute manpower planning, recruitment, selection, performance appraisal, and industrial relations",
				}},
			},
		},
		{
			Name:        "400 Level: Strategic Management, Corporate Governance, and B.Sc. Project",
			Description: "BUS 401 Strategic Management, Business Policy, International Business, and Thesis.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "BUS 401: Strategic Management and Business Policy", Objectives: []string{
					"Perform SWOT analysis, Porter's Five Forces analysis, and formulate corporate strategies",
					"Understand corporate governance, business ethics, and corporate social responsibility (CSR)",
				}},
				{Name: "BUS 499: Final Year Business Research Project", Objectives: []string{
					"Conduct empirical business research and present undergraduate research defense",
				}},
			},
		},
	}
}
