package store

import (
	"channel-test/pkg/models"
	"errors"
	"sort"
	"sync"
	"time"
)

var (
	// ErrStudentNotFound is returned when a student ID is not found
	ErrStudentNotFound = errors.New("student not found")

	// ErrExamNotFound is returned when an exam number is not found
	ErrExamNotFound = errors.New("exam not found")
)

// MemoryStore implements the Store interface using in-memory storage
type MemoryStore struct {
	mu     sync.RWMutex
	scores map[string]map[int]models.StudentScore // studentID -> examNumber -> score
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		scores: make(map[string]map[int]models.StudentScore),
	}
}

// AddScore adds a new score event to the store
func (s *MemoryStore) AddScore(event models.ScoreEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.scores[event.StudentID] == nil {
		s.scores[event.StudentID] = make(map[int]models.StudentScore)
	}

	s.scores[event.StudentID][event.Exam] = models.StudentScore{
		Exam:      event.Exam,
		Score:     event.Score,
		Timestamp: time.Now(),
	}

	return nil
}

// GetAllStudents returns a sorted list of all student IDs
func (s *MemoryStore) GetAllStudents() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	students := make([]string, 0, len(s.scores))
	for studentID := range s.scores {
		students = append(students, studentID)
	}

	sort.Strings(students)
	return students
}

// GetStudent returns detailed information about a specific student
func (s *MemoryStore) GetStudent(id string) (*models.Student, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	exams, exists := s.scores[id]
	if !exists {
		return nil, ErrStudentNotFound
	}

	scores := make([]models.StudentScore, 0, len(exams))
	var totalScore float64

	for _, score := range exams {
		scores = append(scores, score)
		totalScore += score.Score
	}

	// Sort scores by exam number
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Exam < scores[j].Exam
	})

	averageScore := 0.0
	if len(scores) > 0 {
		averageScore = totalScore / float64(len(scores))
	}

	return &models.Student{
		ID:           id,
		Scores:       scores,
		AverageScore: averageScore,
	}, nil
}

// GetAllExams returns a sorted list of all exam numbers
func (s *MemoryStore) GetAllExams() []int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	examSet := make(map[int]bool)
	for _, exams := range s.scores {
		for examNum := range exams {
			examSet[examNum] = true
		}
	}

	exams := make([]int, 0, len(examSet))
	for examNum := range examSet {
		exams = append(exams, examNum)
	}

	sort.Ints(exams)
	return exams
}

// GetExam returns detailed information about a specific exam
func (s *MemoryStore) GetExam(number int) (*models.Exam, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]models.ExamResult, 0)
	var totalScore float64

	for studentID, exams := range s.scores {
		if score, exists := exams[number]; exists {
			results = append(results, models.ExamResult{
				StudentID: studentID,
				Score:     score.Score,
			})
			totalScore += score.Score
		}
	}

	if len(results) == 0 {
		return nil, ErrExamNotFound
	}

	// Sort results by student ID
	sort.Slice(results, func(i, j int) bool {
		return results[i].StudentID < results[j].StudentID
	})

	averageScore := totalScore / float64(len(results))

	return &models.Exam{
		Number:       number,
		Results:      results,
		AverageScore: averageScore,
	}, nil
}
