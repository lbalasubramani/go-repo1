package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lbalasubramani/go-repo1/pkg/logger"
)

func LoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			next.ServeHTTP(lw, r)
			
			log.Info(fmt.Sprintf(
				"%s %s %d %s",
				r.Method,
				r.URL.Path,
				lw.statusCode,
				time.Since(start),
			))
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

func RecoveryMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error(fmt.Sprintf("Panic recovered: %v", err))
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
