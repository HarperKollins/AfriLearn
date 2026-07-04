package handlers

import (
	"testing"
)

func TestParseUserIntent(t *testing.T) {
	tests := []struct {
		input          string
		expectedBoard  string
		expectedSubj   string
		expectedAction string
	}{
		{
			input:          "Give me the JSS3 math syllabus",
			expectedBoard:  "bece",
			expectedSubj:   "mathematics",
			expectedAction: "curriculum",
		},
		{
			input:          "Explain SS2 physics for WAEC",
			expectedBoard:  "waec",
			expectedSubj:   "physics",
			expectedAction: "curriculum",
		},
		{
			input:          "LLM prompt for UNILAG Law",
			expectedBoard:  "unilag",
			expectedSubj:   "law",
			expectedAction: "llm-prompt",
		},
		{
			input:          "What are the prerequisites for quadratic equations in WAEC math?",
			expectedBoard:  "waec",
			expectedSubj:   "mathematics",
			expectedAction: "prerequisites",
		},
		{
			input:          "Learning pathway for WAEC mathematics",
			expectedBoard:  "waec",
			expectedSubj:   "mathematics",
			expectedAction: "pathway",
		},
	}

	for _, tt := range tests {
		intent := ParseUserIntent(tt.input)
		if intent.Board != tt.expectedBoard {
			t.Errorf("ParseUserIntent(%q) Board = %q; expected %q", tt.input, intent.Board, tt.expectedBoard)
		}
		if intent.Subject != tt.expectedSubj {
			t.Errorf("ParseUserIntent(%q) Subject = %q; expected %q", tt.input, intent.Subject, tt.expectedSubj)
		}
		if intent.Action != tt.expectedAction {
			t.Errorf("ParseUserIntent(%q) Action = %q; expected %q", tt.input, intent.Action, tt.expectedAction)
		}
	}
}
