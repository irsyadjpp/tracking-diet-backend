package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
)

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// ErrorHandler handles panics and returns proper error responses
type ErrorHandler struct{}

// NewErrorHandler creates a new error handler middleware
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// Recover handles panics and converts them to error responses
func (e *ErrorHandler) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the stack trace
				debug.PrintStack()
				
				// Return internal server error
				e.writeError(w, http.StatusInternalServerError, "Internal server error", "INTERNAL_ERROR")
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

// writeError writes an error response in JSON format
func (e *ErrorHandler) writeError(w http.ResponseWriter, status int, message string, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	response := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Code:    code,
	}
	
	json.NewEncoder(w).Encode(response)
}

// ErrorResponder is a helper function to write error responses
func ErrorResponder(w http.ResponseWriter, status int, message string, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	response := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Code:    code,
	}
	
	json.NewEncoder(w).Encode(response)
}

// JSONResponder is a helper function to write JSON responses
func JSONResponder(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}