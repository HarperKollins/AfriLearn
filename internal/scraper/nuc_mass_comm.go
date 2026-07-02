package scraper

// NUCMassCommScraper implements Scraper for NUC CCMAS B.Sc. Mass Communication (100L - 400L)
type NUCMassCommScraper struct{}

func NewNUCMassCommScraper() *NUCMassCommScraper {
	return &NUCMassCommScraper{}
}

func (s *NUCMassCommScraper) BoardSlug() string   { return "nuc" }
func (s *NUCMassCommScraper) SubjectSlug() string { return "mass-communication" }
func (s *NUCMassCommScraper) Level() string       { return "tertiary-degree" }
func (s *NUCMassCommScraper) SourceURL() string {
	return "https://nuc.edu.ng/ccmas/arts/mass-communication"
}

func (s *NUCMassCommScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "100 - 200 Level: Introduction to Mass Communication and News Writing",
			Description: "MAC 101 Introduction to Mass Comm, History of Nigerian Media, and News Reporting.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "MAC 101: Introduction to Mass Communication", Objectives: []string{
					"Understand mass communication theories, models (Lasswell, Shannon-Weaver), and media functions",
					"Trace history of print and broadcast media in Nigeria from 1859 Iwe Irohin to present",
				}},
				{Name: "MAC 201: News Writing and Reporting", Objectives: []string{
					"Apply lead writing, inverted pyramid style, interviewing techniques, and news values",
				}},
			},
		},
		{
			Name:        "300 Level: Public Relations, Advertising, and Broadcasting",
			Description: "MAC 301 Public Relations, Advertising, Radio/TV Production, and Photojournalism.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "MAC 301 & MAC 302: Public Relations and Advertising", Objectives: []string{
					"Design public relations campaigns, press releases, crisis management, and media relations",
					"Execute advertising strategy, copy writing, layout design, and media planning",
				}},
				{Name: "MAC 303: Radio and Television Production", Objectives: []string{
					"Produce radio programs, audio editing, television news scripting, camera operations, and video editing",
				}},
			},
		},
		{
			Name:        "400 Level: Media Law, Ethics, Digital Media, and B.Sc. Project",
			Description: "MAC 401 Media Law & Ethics, Digital Journalism, Development Comm, and Research Project.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "MAC 401: Media Law, Ethics, and Digital Media", Objectives: []string{
					"Understand libel, defamation, sedition, copyright law, Freedom of Information (FOI) Act, and NBC codes",
					"Apply digital journalism techniques, online publishing, data journalism, and social media management",
				}},
				{Name: "MAC 499: Final Year Mass Communication Research Project", Objectives: []string{
					"Execute empirical content analysis, survey, or audience research and defend thesis",
				}},
			},
		},
	}
}
