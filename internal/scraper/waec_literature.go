package scraper

// WAECLiteratureScraper implements Scraper for WAEC Literature in English (SS1 - SS3)
type WAECLiteratureScraper struct{}

func NewWAECLiteratureScraper() *WAECLiteratureScraper {
	return &WAECLiteratureScraper{}
}

func (s *WAECLiteratureScraper) BoardSlug() string   { return "waec" }
func (s *WAECLiteratureScraper) SubjectSlug() string { return "literature" }
func (s *WAECLiteratureScraper) Level() string       { return "senior-secondary" }
func (s *WAECLiteratureScraper) SourceURL() string {
	return "https://waecsyllabus.com/literature-in-english-syllabus/"
}

func (s *WAECLiteratureScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Introduction to Literature and Literary Appreciation",
			Description: "Literary genres, literary devices, figures of speech, poetic devices, and analytical techniques.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Literary Genres and Concepts", Objectives: []string{
					"Distinguish Prose, Poetry, and Drama genres and sub-genres",
					"Identify plot, theme, characterization, setting, point of view, and atmosphere",
				}},
				{Name: "Figures of Speech and Poetic Devices", Objectives: []string{
					"Identify and interpret simile, metaphor, personification, irony, oxymoron, hyperbole, paradox, symbolism",
					"Analyze rhyme scheme, rhythm, meter, stanza, alliteration, assonance, and imagery in poetry",
				}},
			},
		},
		{
			Name:        "African and Non-African Prose",
			Description: "Critical analysis of set African and Non-African novels, themes, characters, and stylistic devices.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "African Prose Analysis", Objectives: []string{
					"Analyze plot summary, character roles, socio-cultural background, and themes in prescribed African novels",
					"Evaluate narrative techniques, symbolism, and linguistic style in African prose",
				}},
				{Name: "Non-African Prose Analysis", Objectives: []string{
					"Analyze plot, themes, characterization, and historical context in prescribed Non-African novels",
				}},
			},
		},
		{
			Name:        "African and Non-African Drama and Poetry",
			Description: "Critical analysis of set African/Non-African plays (tragedy, comedy) and prescribed poems.",
			Difficulty:  "hard",
			Subtopics: []SubtopicData{
				{Name: "African and Non-African Drama", Objectives: []string{
					"Analyze dramatic structure, acts, scenes, conflict, dialogue, monologue, soliloquy, and dramatic irony in prescribed plays",
				}},
				{Name: "African and Non-African Poetry", Objectives: []string{
					"Analyze theme, tone, mood, structure, imagery, and poetic devices in prescribed African and Non-African poems",
				}},
			},
		},
	}
}
