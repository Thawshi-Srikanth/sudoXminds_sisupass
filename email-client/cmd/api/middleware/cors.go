package middleware

import (
	"net/http"

	"sisupass.com/email-client/cmd/api/app"
)

func EnableCORS(app *app.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add Vary header to indicate that the response varies based on Origin

			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Method")

			origin := r.Header.Get("Origin")

			if origin != "" {
				// Check if origin is in trusted origins
				for i := range app.Config.CORS.TrustedOrigins {
					if origin == app.Config.CORS.TrustedOrigins[i] {
						w.Header().Set("Access-Control-Allow-Origin", origin)

						// Handle preflight requests
						if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
							w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
							w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
							w.WriteHeader(http.StatusOK)
							return

						}
						break
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
