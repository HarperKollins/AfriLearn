package scraper

// BECEEnglishScraper implements Scraper for BECE English Studies (JSS1 - JSS3)
type BECEEnglishScraper struct{}

func NewBECEEnglishScraper() *BECEEnglishScraper {
	return &BECEEnglishScraper{}
}

func (s *BECEEnglishScraper) BoardSlug() string   { return "bece" }
func (s *BECEEnglishScraper) SubjectSlug() string { return "english-language" }
func (s *BECEEnglishScraper) Level() string       { return "junior-secondary" }
func (s *BECEEnglishScraper) SourceURL() string {
	return "https://nerdc.gov.ng/curriculum/junior-secondary/english-studies"
}

func (s *BECEEnglishScraper) Topics() []TopicData {
	return []TopicData{
		{
			Name:        "Reading Comprehension and Summary (JSS1 - JSS3)",
			Description: "Reading strategies, main idea extraction, vocabulary development, and summary writing skills.",
			Difficulty:  "easy",
			Subtopics: []SubtopicData{
				{Name: "Reading Comprehension Skills", Objectives: []string{
					"Identify main ideas and supporting details in narrative, expository, and descriptive passages",
					"Infer word meanings from contextual clues in reading passages",
					"Answer factual, inferential, and evaluative comprehension questions",
				}},
				{Name: "Summary Writing Techniques", Objectives: []string{
					"Extract essential points from passages and eliminate redundant information",
					"Write concise summary sentences using own words",
				}},
			},
		},
		{
			Name:        "Grammar, Mechanics, and Vocabulary (JSS1 - JSS3)",
			Description: "Parts of speech, sentence types, tenses, subject-verb agreement, punctuation, and vocabulary.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Parts of Speech and Sentence Structure", Objectives: []string{
					"Identify and use nouns, pronouns, verbs, adjectives, adverbs, prepositions, conjunctions",
					"Apply subject-verb agreement (concord) rules in sentences",
					"Distinguish simple, compound, and complex sentences",
				}},
				{Name: "Tenses, Voice, and Punctuation", Objectives: []string{
					"Use present, past, future, perfect, and continuous tenses correctly",
					"Convert sentences between active voice and passive voice",
					"Use capital letters, full stops, commas, question marks, quotation marks, and apostrophes correctly",
				}},
			},
		},
		{
			Name:        "Composition and Essay Writing (JSS1 - JSS3)",
			Description: "Formal and informal letters, narrative, descriptive, argumentative, and expository essays.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Letter Writing", Objectives: []string{
					"Write informal letters to friends and family with proper layout and tone",
					"Write formal letters (applications, complaints, requests) with formal layout and language",
				}},
				{Name: "Essay and Narrative Writing", Objectives: []string{
					"Write narrative essays recounting personal experiences and events chronologically",
					"Write descriptive essays describing persons, places, events, and objects vividly",
					"Write argumentative essays presenting logical arguments for or against a motion",
				}},
			},
		},
		{
			Name:        "Literature in English and Oracy (JSS1 - JSS3)",
			Description: "Prose, poetry, drama, literary terms, figures of speech, and oral English phonetics.",
			Difficulty:  "medium",
			Subtopics: []SubtopicData{
				{Name: "Literature Genres and Figures of Speech", Objectives: []string{
					"Analyze plot, themes, characterization, and setting in prose, drama, and poetry texts",
					"Identify figures of speech: simile, metaphor, personification, hyperbole, alliteration, assonance",
				}},
				{Name: "Oral English and Phonetics", Objectives: []string{
					"Identify and pronounce English vowel sounds (monophthongs, diphthongs) and consonant sounds",
					"Identify stress patterns in words and intonation patterns in sentences",
				}},
			},
		},
	}
}
