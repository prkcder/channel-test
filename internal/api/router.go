package api

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// NewRouter creates and configures the HTTP router
func NewRouter(handler *Handler) http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/health", handler.HealthCheck)
	mux.HandleFunc("/students/", handleStudentsRoutes(handler))
	mux.HandleFunc("/exams/", handleExamsRoutes(handler))

	mux.HandleFunc("/", handler.Index)

	// Wrap with logging middleware
	return loggingMiddleware(mux)
}

// handleStudentsRoutes routes requests for /students and /students/{id}
func handleStudentsRoutes(handler *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Exact match for /students
		if path == "/students" || path == "/students/" {
			handler.ListStudents(w, r)
			return
		}

		// Match /students/{id}
		if strings.HasPrefix(path, "/students/") {
			handler.GetStudent(w, r)
			return
		}

		handler.GetStudent(w, r)
	}
}

// handleExamsRoutes routes requests for /exams and /exams/{number}
func handleExamsRoutes(handler *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Exact match for /exams
		if path == "/exams" || path == "/exams/" {
			handler.ListExams(w, r)
			return
		}

		// Match /exams/{number}
		if strings.HasPrefix(path, "/exams/") {
			handler.GetExam(w, r)
			return
		}

		handler.GetExam(w, r)
	}
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, wrapped.statusCode, duration)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
