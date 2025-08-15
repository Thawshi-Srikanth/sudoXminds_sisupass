package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"sisupass.com/email-client/internal/jsonlog"
)

func Recovery(logger *jsonlog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connection", "close")

					// Log the panic with stack trace
					logger.PrintError(fmt.Errorf("panic: %v", err), map[string]string{
						"method":    r.Method,
						"url":       r.URL.String(),
						"remote_ip": r.RemoteAddr,
						"stack":     string(debug.Stack()),
					})

					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
