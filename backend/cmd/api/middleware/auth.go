package middleware

import (
	"errors"
	"net/http"
	"strings"

	"sisupass.com/sisupass/cmd/api/app"
	"sisupass.com/sisupass/internal/types"
)

func RequireAuthentication(app *app.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Extract Authorization header
			authorizationHeader := r.Header.Get("Authorization")

			if authorizationHeader == "" {
				r = app.ContextSetUser(r, &types.User{})
				next.ServeHTTP(w, r)
				return
			}

			// Parse Bearer token
			headerParts := strings.Split(authorizationHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				app.InvalidAuthenticationTokenResponse(w, r)
				return
			}

			token := headerParts[1]

			user, err := app.Models.Users.GetForToken("authentication", token)
			if err != nil {
				switch {
				case errors.Is(err, types.ErrRecordNotFound):
					app.InvalidAuthenticationTokenResponse(w, r)
				default:
					app.ServerErrorResponse(w, r, err)
				}
				return
			}

			// Add user to request context
			r = app.ContextSetUser(r, user)
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuthenticatedUser(app *app.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := app.ContextGetUser(r)

			if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
				app.AuthenticationRequiredResponse(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireActivatedUser(app *app.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := app.ContextGetUser(r)

			if !user.Activated {
				app.InactiveAccountResponse(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
