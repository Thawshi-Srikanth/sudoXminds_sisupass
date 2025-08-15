package handlers

import (
	"net/http"

	"sisupass.com/email-client/cmd/api/app"
)

// Healthcheck endpoint for API health monitoring
// @Summary		Health check
// @Description	Check API health status and system information
// @Tags			System
// @Produce		json
// @Success		200	{object}	app.HealthResponse	"API is healthy and operational"
// @Router			/health [get]
func Healthcheck(appPtr *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		env := map[string]any{
			"status": "available",
			"system_info": map[string]string{
				"environment": appPtr.Config.Env,
				"version":     appPtr.Config.Version,
				"service":     "email-service",
			},
		}
		err := appPtr.WriteJSON(w, http.StatusOK, env, nil)
		if err != nil {
			appPtr.ServerErrorResponse(w, r, err)
		}
	}
}
