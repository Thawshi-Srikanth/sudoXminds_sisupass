package middleware

import (
	"net/http"
	"strconv"
	"time"

	"sisupass.com/sisupass/internal/jsonlog"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	bytesWritten, err := rw.ResponseWriter.Write(data)
	rw.bytesWritten += bytesWritten
	return bytesWritten, err
}

func RequestLogger(logger *jsonlog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// wrap the response writer
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			logger.PrintInfo("request started", map[string]string{
				"method":     r.Method,
				"path":       r.URL.Path,
				"remote_ip":  r.RemoteAddr,
				"user_agent": r.UserAgent(),
			})

			next.ServeHTTP(rw, r)

			// Log request completion
			duration := time.Since(start)
			logger.PrintInfo("request completed", map[string]string{
				"method":   r.Method,
				"path":     r.URL.Path,
				"status":   strconv.Itoa(rw.statusCode),
				"duration": duration.String(),
				"bytes":    strconv.Itoa(rw.bytesWritten),
			})
		})
	}
}
