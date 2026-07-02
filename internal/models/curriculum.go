package models

import (
	"time"
)

// ExamBoard represents an examination body (WAEC, JAMB, NECO, etc.)
type ExamBoard struct {
	ID          string    `json:"id" db:"id"`
	Slug        string    `json:"slug" db:"slug"`
	Name        string    `json:"name" db:"name"`
	FullName    string    `json:"full_name" db:"full_name"`
	Country     string    `json:"country" db:"country"`
	Description string    `json:"description" db:"description"`
	Website     string    `json:"website,omitempty" db:"website"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Subject represents a school subject (Mathematics, Physics, etc.)
type Subject struct {
	ID          string    `json:"id" db:"id"`
	Slug        string    `json:"slug" db:"slug"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"` // science, arts, commercial
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Curriculum links an ExamBoard to a Subject with metadata
type Curriculum struct {
	ID          string    `json:"id" db:"id"`
	ExamBoardID string    `json:"exam_board_id" db:"exam_board_id"`
	SubjectID   string    `json:"subject_id" db:"subject_id"`
	Year        int       `json:"year" db:"year"`
	Level       string    `json:"level" db:"level"` // ss1, ss2, ss3, jss, primary
	SourceURL   string    `json:"source_url" db:"source_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Joined fields for API responses
	ExamBoard *ExamBoard `json:"exam_board,omitempty" db:"-"`
	Subject   *Subject   `json:"subject,omitempty" db:"-"`
	Topics    []Topic    `json:"topics,omitempty" db:"-"`
}

// Topic represents a major topic within a curriculum (e.g. "Algebra")
type Topic struct {
	ID           string    `json:"id" db:"id"`
	CurriculumID string    `json:"curriculum_id" db:"curriculum_id"`
	Slug         string    `json:"slug" db:"slug"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	OrderIndex   int       `json:"order_index" db:"order_index"`
	Difficulty   string    `json:"difficulty" db:"difficulty"` // easy, medium, hard
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`

	// Joined fields
	Subtopics []Subtopic `json:"subtopics,omitempty" db:"-"`
}

// Subtopic is a specific area under a Topic (e.g. "Quadratic Equations")
type Subtopic struct {
	ID          string    `json:"id" db:"id"`
	TopicID     string    `json:"topic_id" db:"topic_id"`
	Slug        string    `json:"slug" db:"slug"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	OrderIndex  int       `json:"order_index" db:"order_index"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Joined fields
	Objectives []LearningObjective `json:"objectives,omitempty" db:"-"`
}

// LearningObjective is what a student should be able to do after studying a subtopic
type LearningObjective struct {
	ID          string    `json:"id" db:"id"`
	SubtopicID  string    `json:"subtopic_id" db:"subtopic_id"`
	Description string    `json:"description" db:"description"`
	Verb        string    `json:"verb" db:"verb"` // Bloom's taxonomy: remember, understand, apply, analyze, evaluate, create
	OrderIndex  int       `json:"order_index" db:"order_index"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// APIResponse wraps all API responses in a consistent format
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta holds pagination and source metadata
type Meta struct {
	Total       int    `json:"total,omitempty"`
	Page        int    `json:"page,omitempty"`
	PerPage     int    `json:"per_page,omitempty"`
	Source      string `json:"source,omitempty"`
	LastUpdated string `json:"last_updated,omitempty"`
	Version     string `json:"version"`
}
