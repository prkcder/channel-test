package store

import "channel-test/pkg/models"

// Store defines the interface for storing and retrieving test scores
type Store interface {
	// AddScore adds a new score event to the store
	AddScore(event models.ScoreEvent) error

	// GetAllStudents returns a list of all student IDs
	GetAllStudents() []string

	// GetStudent returns detailed information about a specific student
	GetStudent(id string) (*models.Student, error)

	// GetAllExams returns a list of all exam numbers
	GetAllExams() []int

	// GetExam returns detailed information about a specific exam
	GetExam(number int) (*models.Exam, error)
}
