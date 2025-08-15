package routes

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    httpSwagger "github.com/swaggo/http-swagger"
    "sisupass.com/email-client/cmd/api/app"
    "sisupass.com/email-client/cmd/api/handlers"
    custommiddleware "sisupass.com/email-client/cmd/api/middleware"
    "golang.org/x/time/rate"
)

func SetupRoutes(app *app.Application) *chi.Mux {
    r := chi.NewRouter()

    rateLimiter := custommiddleware.NewRateLimiter(
        rate.Limit(app.Config.Limiter.RPS), 
        app.Config.Limiter.Burst,
    )

    r.Use(custommiddleware.Recovery(app.Logger))
    r.Use(custommiddleware.RequestLogger(app.Logger))
    r.Use(custommiddleware.EnableCORS(app))
    r.Use(custommiddleware.RateLimit(app, rateLimiter))
    r.Use(middleware.Heartbeat("/ping"))

    r.Get("/health", handlers.Healthcheck(app))

    r.Get("/swagger/*", httpSwagger.Handler(
        httpSwagger.URL("http://localhost:4000/swagger/doc.json"),
    ))

    emailHandler := handlers.NewEmailHandler(app.Services.Mail)

    r.Route("/api/v1", func(r chi.Router) {
        r.Route("/email", func(r chi.Router) {
            r.Post("/send", emailHandler.SendEmail(app))
            r.Post("/bulk", emailHandler.SendBulkEmail(app))
            r.Get("/{id}/status", emailHandler.GetEmailStatus(app))
            r.Get("/templates", emailHandler.GetEmailTemplates(app))
            r.Get("/stats", emailHandler.GetEmailStats(app))
        })
    })

    return r
}