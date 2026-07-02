package scraper

import (
	"strings"
)

// Scraper is the interface every curriculum dataset provider must implement.
// Each exam board + subject combination gets its own concrete implementation.
type Scraper interface {
	// BoardSlug returns the exam board identifier (e.g. "waec", "jamb")
	BoardSlug() string

	// SubjectSlug returns the subject identifier (e.g. "mathematics", "physics", "biology")
	SubjectSlug() string

	// Level returns the educational level (e.g. "senior-secondary", "tertiary-entry")
	Level() string

	// SourceURL returns the authoritative source URL this curriculum data was derived from
	SourceURL() string

	// Topics returns the complete, verified topic tree for this curriculum.
	// Every scraper implements this as its "Golden Record" — the authoritative,
	// human-reviewed, structured dataset for that exam & subject.
	Topics() []TopicData
}

// TopicData holds a top-level topic and all its children.
type TopicData struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Difficulty  string         `json:"difficulty"` // "easy" | "medium" | "hard"
	Subtopics   []SubtopicData `json:"subtopics"`
}

// SubtopicData holds a subtopic and its learning objectives.
type SubtopicData struct {
	Name       string   `json:"name"`
	Objectives []string `json:"objectives"`
}

// Slugify converts any string to a URL-safe slug format.
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " & ", "-and-")
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ":", "")
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}

// ClassifyDifficulty assigns a difficulty classification based on topic title & keywords.
func ClassifyDifficulty(topicName string, defaultLevel string) string {
	if defaultLevel != "" {
		return defaultLevel
	}
	hardKeywords := []string{"calculus", "matrices", "vectors", "quantum", "nuclear", "electronics", "electromagnetism", "genetics", "organic"}
	mediumKeywords := []string{"trigonometry", "coordinate", "statistics", "probability", "algebra", "waves", "optics", "thermodynamics", "ecology"}

	lower := strings.ToLower(topicName)
	for _, kw := range hardKeywords {
		if strings.Contains(lower, kw) {
			return "hard"
		}
	}
	for _, kw := range mediumKeywords {
		if strings.Contains(lower, kw) {
			return "medium"
		}
	}
	return "easy"
}

// ExtractBloomsVerb extracts the leading action verb from a learning objective.
func ExtractBloomsVerb(objective string) string {
	bloomsVerbs := []string{
		"define", "identify", "calculate", "solve", "apply", "analyze", "analyse",
		"evaluate", "create", "construct", "interpret", "explain", "distinguish",
		"state", "use", "find", "draw", "describe", "determine", "perform", "express",
		"compare", "differentiate", "illustrate", "demonstrate", "relate", "formulate",
	}

	lower := strings.ToLower(strings.TrimSpace(objective))
	for _, verb := range bloomsVerbs {
		if strings.HasPrefix(lower, verb) {
			return verb
		}
	}
	return "understand"
}
