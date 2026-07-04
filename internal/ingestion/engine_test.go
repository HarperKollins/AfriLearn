package ingestion

import (
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Mathematics & Geometry", "mathematics-and-geometry"},
		{"100 Level: Legal Methods / Reasoning", "100-level-legal-methods-reasoning"},
		{"Plane Geometry & Trigonometry?", "plane-geometry-and-trigonometry"},
		{"   Multiple   Spaces   ", "multiple-spaces"},
		{"200–300 Level: Constitutional Law", "200-300-level-constitutional-law"},
	}

	for _, tt := range tests {
		got := Slugify(tt.input)
		if got != tt.expected {
			t.Errorf("Slugify(%q) = %q; expected %q", tt.input, got, tt.expected)
		}
	}
}

func TestExtractBloomsVerb(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Calculate the acceleration of a free-falling body", "calculate"},
		{"Explain the process of photosynthesis", "explain"},
		{"Design a simple electronic circuit", "design"},
		{"Identify the main clauses of the Nigerian Constitution", "identify"},
		{"Unknown objective text", "understand"},
	}

	for _, tt := range tests {
		got := ExtractBloomsVerb(tt.input)
		if got != tt.expected {
			t.Errorf("ExtractBloomsVerb(%q) = %q; expected %q", tt.input, got, tt.expected)
		}
	}
}

func TestClassifyDifficulty(t *testing.T) {
	tests := []struct {
		topicName string
		expected  string
	}{
		{"Introductory Calculus and Derivatives", "hard"},
		{"Trigonometry and Wave Motion", "medium"},
		{"Basic Counting and Sets", "easy"},
	}

	for _, tt := range tests {
		got := ClassifyDifficulty(tt.topicName)
		if got != tt.expected {
			t.Errorf("ClassifyDifficulty(%q) = %q; expected %q", tt.topicName, got, tt.expected)
		}
	}
}

func TestValidateCurriculumFile(t *testing.T) {
	valid := CurriculumFile{
		Board:     "unilag",
		Subject:   "law",
		Level:     "tertiary-degree",
		SourceURL: "https://unilag.edu.ng",
		Topics: []TopicData{
			{
				Name:        "100 Level Legal Methods",
				Description: "Legal methods description",
				Difficulty:  "easy",
				Subtopics: []SubtopicData{
					{
						Name:       "Legal Reasoning",
						Objectives: []string{"Explain statutory interpretation rules"},
					},
				},
			},
		},
	}

	errs := ValidateCurriculumFile(valid, "test.json")
	if len(errs) != 0 {
		t.Fatalf("Expected 0 validation errors for valid file, got %d", len(errs))
	}

	invalid := CurriculumFile{
		Board:   "",
		Subject: "math",
		Level:   "",
		Topics:  []TopicData{},
	}

	invalidErrs := ValidateCurriculumFile(invalid, "bad.json")
	if len(invalidErrs) == 0 {
		t.Fatalf("Expected validation errors for malformed file, got 0")
	}
}
