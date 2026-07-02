package scraper

// UNECLawScraper implements Scraper for UNEC LL.B. Law (Faculty of Law, Enugu Campus)
type UNECLawScraper struct{}

func NewUNECLawScraper() *UNECLawScraper {
	return &UNECLawScraper{}
}

func (s *UNECLawScraper) BoardSlug() string   { return "unec" }
func (s *UNECLawScraper) SubjectSlug() string { return "law" }
func (s *UNECLawScraper) Level() string       { return "tertiary-degree" }
func (s *UNECLawScraper) SourceURL() string {
	return "https://unn.edu.ng/academics/faculties/law-enugu-campus"
}

func (s *UNECLawScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 - 200 Level: UNEC Legal System, Constitutional Law & Contracts",
			Description: "Nigerian Legal System, Legal Methods, Constitutional Law I & II, and Law of Contract at UNEC Enugu.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "LAW 101: Nigerian Legal System and Legal Methods (UNEC)", Objectives: []string{
					"Analyze sources of Nigerian law, judicial precedents, and law court structures at UNEC Enugu Campus",
				}},
				{Name: "LAW 201 & LAW 202: Constitutional Law and Law of Contract", Objectives: []string{
					"Study 1999 Nigerian Constitution, human rights enforcement, contract formation, and breach remedies",
				}},
			},
		},
		{
			Name:        "300 - 500 Level: Criminal Law, Land Law, Evidence, Jurisprudence, and UNEC Thesis",
			Description: "Criminal Law, Land Law, Law of Torts, Equity & Trusts, CAMA 2020, and LL.B. Thesis.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "LAW 301 & LAW 401: Criminal Law, Torts, Land Law, and Evidence", Objectives: []string{
					"Analyze Criminal Code/Penal Code offenses, Land Use Act 1978, torts of negligence, and Evidence Act 2011",
				}},
				{Name: "LAW 501 & LAW 599: Equity, CAMA 2020, Jurisprudence, and LL.B. Project", Objectives: []string{
					"Master corporate law, jurisprudence theories, and defend legal thesis before UNEC Faculty Board",
				}},
			},
		},
	}
}
