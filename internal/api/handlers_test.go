package api

import (
	"channel-test/internal/store"
	"channel-test/pkg/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestStore() store.Store {
	s := store.NewMemoryStore()
	
	// Add test data
	testEvents := []models.ScoreEvent{
		{Exam: 1, StudentID: "alice", Score: 0.85},
		{Exam: 2, StudentID: "alice", Score: 0.90},
		{Exam: 1, StudentID: "bob", Score: 0.75},
		{Exam: 1, StudentID: "charlie", Score: 0.95},
		{Exam: 2, StudentID: "bob", Score: 0.80},
	}
	
	for _, event := range testEvents {
		s.AddScore(event)
	}
	
	return s
}

func TestHandler_ListStudents(t *testing.T) {
	handler := NewHandler(setupTestStore())
	
	req := httptest.NewRequest(http.MethodGet, "/students", nil)
	w := httptest.NewRecorder()
	
	handler.ListStudents(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	students, ok := response["students"].([]interface{})
	if !ok {
		t.Fatal("Expected students array in response")
	}
	
	if len(students) != 3 {
		t.Errorf("Expected 3 students, got %d", len(students))
	}
}

func TestHandler_GetStudent(t *testing.T) {
	handler := NewHandler(setupTestStore())
	
	req := httptest.NewRequest(http.MethodGet, "/students/alice", nil)
	w := httptest.NewRecorder()
	
	handler.GetStudent(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var student models.Student
	if err := json.NewDecoder(w.Body).Decode(&student); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if student.ID != "alice" {
		t.Errorf("Expected student alice, got %s", student.ID)
	}
	
	if len(student.Scores) != 2 {
		t.Errorf("Expected 2 scores for alice, got %d", len(student.Scores))
	}
	
	expectedAvg := (0.85 + 0.90) / 2
	if student.AverageScore != expectedAvg {
		t.Errorf("Expected average %.2f, got %.2f", expectedAvg, student.AverageScore)
	}
}

func TestHandler_GetStudent_NotFound(t *testing.T) {
	handler := NewHandler(setupTestStore())
	
	req := httptest.NewRequest(http.MethodGet, "/students/nonexistent", nil)
	w := httptest.NewRecorder()
	
	handler.GetStudent(w, req)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestHandler_GetStudent_EmptyID(t *testing.T) {
	handler := NewHandler(setupTestStore())
	
	req := httptest.NewRequest(http.MethodGet, "/students/", nil)
	w := httptest.NewRecorder()
	
	handler.GetStudent(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandler_ListExams(t *testing.T) {
	handler := NewHandler(setupTestStore())
	
	req := httptest.NewRequest(http.MethodGet, "/exams", nil)
	w := httptest.NewRecorder()
	
	handler.ListExams(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	exams, ok := response["exams"].([]interface{})
	if !ok {
		t.Fatal("Expected exams array in response")
	}
	
	if len(exams) != 2 {
		t.Errorf("Expected 2 exams, got %d", len(exams))
	}
}

func TestHandler_GetExam(t *testing.T) {
	handler := NewHandler(setupTestStore())
	
	req := httptest.NewRequest(http.MethodGet, "/exams/1", nil)
	w := httptest.NewRecorder()
	
	handler.GetExam(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var exam models.Exam
	if err := json.NewDecoder(w.Body).Decode(&exam); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if exam.Number != 1 {
		t.Errorf("Expected exam 1, got %d", exam.Number)
	}
	
	if len(exam.Results) != 3 {
		t.Errorf("Expected 3 results for exam 1, got %d", len(exam.Results))
	}
	
	expectedAvg := (0.85 + 0.75 + 0.95) / 3
	if exam.AverageScore != expectedAvg {
		t.Errorf("Expected average %.4f, got %.4f", expectedAvg, exam.AverageScore)
	}
}

func TestHandler_GetExam_NotFound(t *testing.T) {
	handler := NewHandler(setupTestStore())
	
	req := httptest.NewRequest(http.MethodGet, "/exams/999", nil)
	w := httptest.NewRecorder()
	
	handler.GetExam(w, req)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestHandler_GetExam_InvalidNumber(t *testing.T) {
	handler := NewHandler(setupTestStore())
	
	req := httptest.NewRequest(http.MethodGet, "/exams/invalid", nil)
	w := httptest.NewRecorder()
	
	handler.GetExam(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandler_HealthCheck(t *testing.T) {
	handler := NewHandler(setupTestStore())
	
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	
	handler.HealthCheck(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if response["status"] != "healthy" {
		t.Errorf("Expected status healthy, got %s", response["status"])
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	handler := NewHandler(setupTestStore())
	
	tests := []struct {
		method string
		path   string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{http.MethodPost, "/students", handler.ListStudents},
		{http.MethodPut, "/students/alice", handler.GetStudent},
		{http.MethodDelete, "/exams", handler.ListExams},
		{http.MethodPatch, "/exams/1", handler.GetExam},
	}
	
	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		
		tt.handler(w, req)
		
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s %s: Expected status 405, got %d", tt.method, tt.path, w.Code)
		}
	}
}

func TestExtractPathParam(t *testing.T) {
	tests := []struct {
		path     string
		prefix   string
		expected string
	}{
		{"/students/alice", "/students/", "alice"},
		{"/exams/123", "/exams/", "123"},
		{"/students/alice/", "/students/", "alice"},
		{"/students/", "/students/", ""},
		{"/other/path", "/students/", ""},
	}
	
	for _, tt := range tests {
		result := extractPathParam(tt.path, tt.prefix)
		if result != tt.expected {
			t.Errorf("extractPathParam(%q, %q) = %q, want %q", 
				tt.path, tt.prefix, result, tt.expected)
		}
	}
}
