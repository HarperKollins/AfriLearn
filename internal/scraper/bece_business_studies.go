package scraper

// BECEBusinessStudiesScraper implements Scraper for BECE Business Studies (JSS1 - JSS3)
type BECEBusinessStudiesScraper struct{}

func NewBECEBusinessStudiesScraper() *BECEBusinessStudiesScraper {
	return &BECEBusinessStudiesScraper{}
}

func (s *BECEBusinessStudiesScraper) BoardSlug() string   { return "bece" }
func (s *BECEBusinessStudiesScraper) SubjectSlug() string { return "business-studies" }
func (s *BECEBusinessStudiesScraper) Level() string       { return "junior-secondary" }
func (s *BECEBusinessStudiesScraper) SourceURL() string {
	return "https://nerdc.gov.ng/curriculum/junior-secondary/business-studies"
}

func (s *BECEBusinessStudiesScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Office Practice and Commerce (JSS1 - JSS3)",
			Description: "Office organization, clerical staff duties, filing systems, correspondence, commerce, and trade.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Office Organization and Clerical Duties", Objectives: []string{
					"Identify administrative departments and duties of a receptionist and clerk",
					"Explain filing systems (alphabetical, numerical, subject) and office equipment",
				}},
				{Name: "Commerce and Trade", Objectives: []string{
					"Distinguish home trade (wholesale, retail) and foreign trade (import, export)",
					"Describe roles of commercial banks, central bank, insurance, and advertising in trade",
				}},
			},
		},
		{
			Name:        "Bookkeeping and Financial Records (JSS1 - JSS3)",
			Description: "Source documents, cash book, ledgers, trial balance, and double-entry bookkeeping system.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Source Documents and Journals", Objectives: []string{
					"Identify invoice, receipt, voucher, debit note, and credit note",
					"Record transactions in Sales Journal, Purchases Journal, and Petty Cash Book",
				}},
				{Name: "Double-Entry Bookkeeping and Trial Balance", Objectives: []string{
					"Apply double-entry bookkeeping principle (debit receiver, credit giver)",
					"Post transactions to ledger accounts and extract a Trial Balance",
				}},
			},
		},
	}
}
