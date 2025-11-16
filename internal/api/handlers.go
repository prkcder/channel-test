package api

import (
	"channel-test/internal/store"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// Handler handles HTTP requests for the scores API
type Handler struct {
	store store.Store
}

// NewHandler creates a new API handler
func NewHandler(store store.Store) *Handler {
	return &Handler{
		store: store,
	}
}

// In internal/api/handlers.go, add:
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
    respondJSON(w, http.StatusOK, map[string]interface{}{
        "service": "Test Scores API",
        "version": "1.0.0",
        "endpoints": []string{
            "GET /health",
            "GET /students",
            "GET /students/{id}",
            "GET /exams",
            "GET /exams/{number}",
        },
    })
}

// ListStudents handles GET /students
// Returns all students that have received at least one test score
func (h *Handler) ListStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	students := h.store.GetAllStudents()

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"students": students,
		"count":    len(students),
	})
}

// GetStudent handles GET /students/{id}
// Returns test results and average score for a specific student
func (h *Handler) GetStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract student ID from path
	id := extractPathParam(r.URL.Path, "/students/")
	if id == "" {
		http.Error(w, "Student ID required", http.StatusBadRequest)
		return
	}

	student, err := h.store.GetStudent(id)
	if err != nil {
		if errors.Is(err, store.ErrStudentNotFound) {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, student)
}

// ListExams handles GET /exams
// Returns all exams that have been recorded
func (h *Handler) ListExams(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	exams := h.store.GetAllExams()

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"exams": exams,
		"count": len(exams),
	})
}

// GetExam handles GET /exams/{number}
// Returns all results and average score for a specific exam
func (h *Handler) GetExam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract exam number from path
	numberStr := extractPathParam(r.URL.Path, "/exams/")
	if numberStr == "" {
		http.Error(w, "Exam number required", http.StatusBadRequest)
		return
	}

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		http.Error(w, "Invalid exam number", http.StatusBadRequest)
		return
	}

	exam, err := h.store.GetExam(number)
	if err != nil {
		if errors.Is(err, store.ErrExamNotFound) {
			http.Error(w, "Exam not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, exam)
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// respondJSON writes a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// extractPathParam extracts the parameter from a URL path
// e.g., "/students/foo" with prefix "/students/" returns "foo"
func extractPathParam(path, prefix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	param := strings.TrimPrefix(path, prefix)
	param = strings.TrimSuffix(param, "/")
	return param
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// NotFound handles 404 responses
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusNotFound, ErrorResponse{
		Error:   "not_found",
		Message: "The requested resource was not found",
	})
}
