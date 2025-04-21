package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
)

// LoggingMiddleware logs the incoming HTTP request and response
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a custom response writer to capture the status code
		rw := &responseWriter{w, http.StatusOK}
		
		// Call the next handler
		next.ServeHTTP(rw, r)
		
		// Log the request details
		log.Printf(
			"%s %s %s %d %s",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			rw.statusCode,
			time.Since(start),
		)
	})
}

// responseWriter is a custom response writer that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// CorsMiddleware returns a middleware that handles CORS
func CorsMiddleware() func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.ExposedHeaders([]string{"Content-Length"}),
		handlers.MaxAge(3600),
	)
} 