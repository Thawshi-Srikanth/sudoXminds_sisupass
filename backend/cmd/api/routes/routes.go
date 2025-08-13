package routes

import (
	"github.com/go-chi/chi/v5"
	"sisupass.com/sisupass/cmd/api/app"
	"sisupass.com/sisupass/cmd/api/handlers"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", handlers.Healthcheck(app))

	r.Route("/api/v1/users", func(r chi.Router) {
		r.Post("/register", handlers.RegisterUser(app))
		r.Get("/", handlers.GetUser(app))
	})

	return r
}
