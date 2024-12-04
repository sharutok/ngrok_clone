package main

import (
	"log"
	"net/http"
	"time"
)

type Applicationlist []string

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request details
		log.Printf("Incoming request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Use a ResponseWriter wrapper to capture the status code
		lrw := &LoggingResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(lrw, r)

		// Log response details after the request is processed
		log.Printf(
			"Completed: %s %s with status %d in %v",
			r.Method, r.URL.Path, lrw.StatusCode, time.Since(start),
		)
	})
}

// LoggingResponseWriter is a wrapper for http.ResponseWriter to capture the status code
type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader captures the status code
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
