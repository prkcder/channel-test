package models

import "time"

// ScoreEvent represents an incoming SSE score event
type ScoreEvent struct {
	Exam      int     `json:"exam"`
	StudentID string  `json:"studentId"`
	Score     float64 `json:"score"`
}

// StudentScore represents a single test score for a student
type StudentScore struct {
	Exam      int       `json:"exam"`
	Score     float64   `json:"score"`
	Timestamp time.Time `json:"timestamp"`
}

// Student represents a student with all their scores
type Student struct {
	ID           string         `json:"id"`
	Scores       []StudentScore `json:"scores"`
	AverageScore float64        `json:"averageScore"`
}

// ExamResult represents a single student's result on an exam
type ExamResult struct {
	StudentID string  `json:"studentId"`
	Score     float64 `json:"score"`
}

// Exam represents an exam with all results
type Exam struct {
	Number       int          `json:"number"`
	Results      []ExamResult `json:"results"`
	AverageScore float64      `json:"averageScore"`
}
