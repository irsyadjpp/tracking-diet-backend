package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/irsyadjpp/tracking-diet-backend/pkg/logger"
)

// LoggingMiddleware logs HTTP requests with request ID tracking
type LoggingMiddleware struct {
	logger *logger.Logger
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(log *logger.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: log,
	}
}

// RequestLogger logs each HTTP request with timing information
func (m *LoggingMiddleware) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Generate or extract request ID
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// Add request ID to response header
		w.Header().Set("X-Request-ID", requestID)
		
		// Create logger with request ID
		requestLogger := m.logger.WithRequestID(requestID)
		
		// Log request
		requestLogger.Info("Incoming request", map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
			"query":  r.URL.RawQuery,
			"remote": r.RemoteAddr,
		})
		
		// Wrap response writer to capture status code
		wrapped := &responseWriter{w, http.StatusOK}
		
		// Call next handler
		next.ServeHTTP(wrapped, r)
		
		// Log response
		duration := time.Since(start)
		requestLogger.Info("Request completed", map[string]interface{}{
			"method":    r.Method,
			"path":      r.URL.Path,
			"status":    wrapped.status,
			"duration":  duration.Milliseconds(),
			"remote":    r.RemoteAddr,
		})
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}