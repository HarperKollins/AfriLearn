package scraper

// NUCAccountingScraper implements Scraper for NUC CCMAS B.Sc. Accounting (100L - 400L)
type NUCAccountingScraper struct{}

func NewNUCAccountingScraper() *NUCAccountingScraper {
	return &NUCAccountingScraper{}
}

func (s *NUCAccountingScraper) BoardSlug() string   { return "nuc" }
func (s *NUCAccountingScraper) SubjectSlug() string { return "accounting" }
func (s *NUCAccountingScraper) Level() string       { return "tertiary-degree" }
func (s *NUCAccountingScraper) SourceURL() string {
	return "https://nuc.edu.ng/ccmas/administration/accounting"
}

func (s *NUCAccountingScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 - 200 Level: Principles of Accounting and Business Mathematics",
			Description: "ACC 101, ACC 201 Financial Accounting, Business Law, and Quantitative Analysis.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "ACC 101: Introduction to Financial Accounting I & II", Objectives: []string{
					"Understand accounting principles, double-entry bookkeeping, ledger accounts, and trial balance",
					"Prepare trading, profit and loss accounts, and balance sheet (Statement of Financial Position)",
				}},
				{Name: "ACC 201: Intermediate Financial Accounting", Objectives: []string{
					"Account for partnership creation, admission/retirement of partners, and partnership dissolution",
					"Prepare manufacturing accounts, departmental accounts, branch accounts, and incomplete records",
				}},
			},
		},
		{
			Name:        "300 Level: Cost Accounting, Management Accounting, and Auditing",
			Description: "ACC 301 Cost Accounting, Management Accounting, Auditing Principles, and Financial Reporting.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "ACC 301 & ACC 302: Cost and Management Accounting", Objectives: []string{
					"Analyze material, labor, and overhead costing methods (job, batch, process costing)",
					"Apply marginal costing, absorption costing, CVP analysis, and budgetary control systems",
				}},
				{Name: "ACC 303: Auditing Principles and Practice", Objectives: []string{
					"Understand audit objectives, internal control systems, audit sampling, and ICAN/IFAC ethical standards",
					"Conduct audit verification of assets/liabilities and draft independent auditor's report",
				}},
			},
		},
		{
			Name:        "400 Level: Advanced Accounting, Taxation, and B.Sc. Thesis",
			Description: "ACC 401 Advanced Financial Accounting, Taxation, Corporate Finance, and Research Project.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "ACC 401 & ACC 402: Advanced Financial Reporting and Taxation", Objectives: []string{
					"Prepare consolidated financial statements for parent-subsidiary group companies under IFRS",
					"Calculate personal income tax, company income tax (CIT), petroleum profit tax, and VAT under FIRS rules",
				}},
				{Name: "ACC 499: Final Year Accounting Research Project", Objectives: []string{
					"Write and defend original accounting research thesis adhering to ICAN and NUC academic standards",
				}},
			},
		},
	}
}
