package store

import (
	"channel-test/pkg/models"
	"testing"
)

func TestMemoryStore_AddScore(t *testing.T) {
	store := NewMemoryStore()

	event := models.ScoreEvent{
		Exam:      1,
		StudentID: "student1",
		Score:     0.85,
	}

	err := store.AddScore(event)
	if err != nil {
		t.Fatalf("AddScore failed: %v", err)
	}

	students := store.GetAllStudents()
	if len(students) != 1 {
		t.Errorf("Expected 1 student, got %d", len(students))
	}

	if students[0] != "student1" {
		t.Errorf("Expected student1, got %s", students[0])
	}
}

func TestMemoryStore_GetStudent(t *testing.T) {
	store := NewMemoryStore()

	// Add multiple scores for a student
	events := []models.ScoreEvent{
		{Exam: 1, StudentID: "student1", Score: 0.85},
		{Exam: 2, StudentID: "student1", Score: 0.90},
		{Exam: 3, StudentID: "student1", Score: 0.95},
	}

	for _, event := range events {
		store.AddScore(event)
	}

	student, err := store.GetStudent("student1")
	if err != nil {
		t.Fatalf("GetStudent failed: %v", err)
	}

	if student.ID != "student1" {
		t.Errorf("Expected student1, got %s", student.ID)
	}

	if len(student.Scores) != 3 {
		t.Errorf("Expected 3 scores, got %d", len(student.Scores))
	}

	expectedAvg := (0.85 + 0.90 + 0.95) / 3
	if student.AverageScore != expectedAvg {
		t.Errorf("Expected average %.2f, got %.2f", expectedAvg, student.AverageScore)
	}
}

func TestMemoryStore_GetStudent_NotFound(t *testing.T) {
	store := NewMemoryStore()

	_, err := store.GetStudent("nonexistent")
	if err != ErrStudentNotFound {
		t.Errorf("Expected ErrStudentNotFound, got %v", err)
	}
}

func TestMemoryStore_GetAllStudents(t *testing.T) {
	store := NewMemoryStore()

	events := []models.ScoreEvent{
		{Exam: 1, StudentID: "charlie", Score: 0.85},
		{Exam: 1, StudentID: "alice", Score: 0.90},
		{Exam: 1, StudentID: "bob", Score: 0.95},
	}

	for _, event := range events {
		store.AddScore(event)
	}

	students := store.GetAllStudents()
	if len(students) != 3 {
		t.Errorf("Expected 3 students, got %d", len(students))
	}

	// Should be sorted alphabetically
	expected := []string{"alice", "bob", "charlie"}
	for i, s := range students {
		if s != expected[i] {
			t.Errorf("Expected student %s at position %d, got %s", expected[i], i, s)
		}
	}
}

func TestMemoryStore_GetExam(t *testing.T) {
	store := NewMemoryStore()

	// Add scores for exam 1 from multiple students
	events := []models.ScoreEvent{
		{Exam: 1, StudentID: "student1", Score: 0.80},
		{Exam: 1, StudentID: "student2", Score: 0.90},
		{Exam: 1, StudentID: "student3", Score: 1.00},
	}

	for _, event := range events {
		store.AddScore(event)
	}

	exam, err := store.GetExam(1)
	if err != nil {
		t.Fatalf("GetExam failed: %v", err)
	}

	if exam.Number != 1 {
		t.Errorf("Expected exam 1, got %d", exam.Number)
	}

	if len(exam.Results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(exam.Results))
	}

	expectedAvg := (0.80 + 0.90 + 1.00) / 3
	if exam.AverageScore != expectedAvg {
		t.Errorf("Expected average %.2f, got %.2f", expectedAvg, exam.AverageScore)
	}
}

func TestMemoryStore_GetExam_NotFound(t *testing.T) {
	store := NewMemoryStore()

	_, err := store.GetExam(999)
	if err != ErrExamNotFound {
		t.Errorf("Expected ErrExamNotFound, got %v", err)
	}
}

func TestMemoryStore_GetAllExams(t *testing.T) {
	store := NewMemoryStore()

	events := []models.ScoreEvent{
		{Exam: 3, StudentID: "student1", Score: 0.85},
		{Exam: 1, StudentID: "student1", Score: 0.90},
		{Exam: 2, StudentID: "student2", Score: 0.95},
	}

	for _, event := range events {
		store.AddScore(event)
	}

	exams := store.GetAllExams()
	if len(exams) != 3 {
		t.Errorf("Expected 3 exams, got %d", len(exams))
	}

	// Should be sorted numerically
	expected := []int{1, 2, 3}
	for i, e := range exams {
		if e != expected[i] {
			t.Errorf("Expected exam %d at position %d, got %d", expected[i], i, e)
		}
	}
}

func TestMemoryStore_UpdateScore(t *testing.T) {
	store := NewMemoryStore()

	// Add initial score
	event1 := models.ScoreEvent{
		Exam:      1,
		StudentID: "student1",
		Score:     0.85,
	}
	store.AddScore(event1)

	// Update with new score for same exam
	event2 := models.ScoreEvent{
		Exam:      1,
		StudentID: "student1",
		Score:     0.95,
	}
	store.AddScore(event2)

	student, _ := store.GetStudent("student1")
	if len(student.Scores) != 1 {
		t.Errorf("Expected 1 score (updated), got %d", len(student.Scores))
	}

	if student.Scores[0].Score != 0.95 {
		t.Errorf("Expected score 0.95, got %.2f", student.Scores[0].Score)
	}
}

func TestMemoryStore_ConcurrentAccess(t *testing.T) {
	store := NewMemoryStore()

	// Test concurrent writes
	done := make(chan bool)
	
	for i := 0; i < 10; i++ {
		go func(id int) {
			event := models.ScoreEvent{
				Exam:      id % 3,
				StudentID: "student" + string(rune('0'+id)),
				Score:     0.85,
			}
			store.AddScore(event)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	students := store.GetAllStudents()
	if len(students) != 10 {
		t.Errorf("Expected 10 students after concurrent writes, got %d", len(students))
	}
}
