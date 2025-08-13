package routes

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"sisupass.com/sisupass/cmd/api/app"
	"sisupass.com/sisupass/cmd/api/handlers"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", handlers.Healthcheck(app))

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:4000/swagger/doc.json"),
	))

	r.Route("/api/v1/users", func(r chi.Router) {
		r.Post("/register", handlers.RegisterUser(app))
		r.Get("/", handlers.GetUser(app))
	})

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/password-reset", handlers.SendPasswordResetEmail(app))
		r.Put("/activate", handlers.ActivateUser(app))
		r.Put("/password", handlers.UpdateUserPassword(app))
		r.Get("/google/login", handlers.GoogleLogin(app))
		r.Get("/google/callback", handlers.GoogleCallback(app))
	})

	return r
}
