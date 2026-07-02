package scraper

// NUCLawScraper implements Scraper for NUC CCMAS LL.B. Law (100L - 500L)
type NUCLawScraper struct{}

func NewNUCLawScraper() *NUCLawScraper {
	return &NUCLawScraper{}
}

func (s *NUCLawScraper) BoardSlug() string   { return "nuc" }
func (s *NUCLawScraper) SubjectSlug() string { return "law" }
func (s *NUCLawScraper) Level() string       { return "tertiary-degree" }
func (s *NUCLawScraper) SourceURL() string {
	return "https://nuc.edu.ng/ccmas/law/bachelor-of-laws"
}

func (s *NUCLawScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 - 200 Level: Legal Foundations, Constitutional Law, and Contracts",
			Description: "Nigerian Legal System, Legal Methods, Constitutional Law, Law of Contract, and Customary Law.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "PUL 101 & PUL 102: Nigerian Legal System and Legal Methods", Objectives: []string{
					"Trace sources of Nigerian law: English law, customary law, Islamic law, statutes, judicial precedents",
					"Analyze hierarchy of courts in Nigeria: Supreme Court, Court of Appeal, Federal/State High Courts",
					"Master legal reasoning, case analysis, statutory interpretation, and legal citation",
				}},
				{Name: "PUL 201 & BUL 201: Constitutional Law and Law of Contract", Objectives: []string{
					"Analyze fundamental rights, separation of powers, judicial review, and executive emergency powers",
					"Study elements of valid contract: offer, acceptance, consideration, intention to create legal relations, capacity",
					"Analyze remedies for breach of contract: damages, specific performance, injunctions",
				}},
			},
		},
		{
			Name:        "300 Level: Criminal Law, Torts, and Commercial Law",
			Description: "Criminal Law, Law of Torts, Commercial Law, and Land Law.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "PUL 301 & BUL 301: Criminal Law and Law of Torts", Objectives: []string{
					"Analyze actus reus, mens rea, general defenses, murder, manslaughter, stealing, fraud, and armed robbery",
					"Study torts of negligence, nuisance, defamation, trespass to person, and trespass to land",
				}},
				{Name: "PUL 302 & BUL 302: Land Law and Commercial Law", Objectives: []string{
					"Analyze Land Use Act 1978, customary land tenure, Certificate of Occupancy (C of O), and mortgages",
					"Study agency law, sale of goods, hire purchase, and negotiable instruments",
				}},
			},
		},
		{
			Name:        "400 - 500 Level: Equity, Evidence, Jurisprudence, Company Law, and LL.B. Thesis",
			Description: "Equity & Trusts, Law of Evidence, Jurisprudence & Legal Theory, Company Law, and Thesis.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "PUL 401 & BUL 401: Equity & Trusts and Law of Evidence", Objectives: []string{
					"Apply maxims of equity, equitable remedies, and express, implied, and constructive trusts",
					"Study Evidence Act 2011: admissibility of evidence, burden of proof, hearsay, and electronic evidence",
				}},
				{Name: "PUL 501 & BUL 501: Jurisprudence, Company Law, and LL.B. Thesis", Objectives: []string{
					"Analyze natural law, positivism, historical, sociological, and realist legal philosophies",
					"Study Companies and Allied Matters Act (CAMA 2020): incorporation, corporate governance, shares, and winding up",
					"Write and defend independent LL.B. legal research project thesis before Faculty Board",
				}},
			},
		},
	}
}
