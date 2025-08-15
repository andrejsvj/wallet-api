package middleware

import (
	"net/http"
	"time"
	"wallet-api/utils/logger"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)

		if wrappedWriter.statusCode >= 400 {
			logger.GlobalLogger.Error(
				"Request failed: %s %s - Status: %d - Duration: %v",
				r.Method, r.URL.Path, wrappedWriter.statusCode, duration,
			)
		}
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
