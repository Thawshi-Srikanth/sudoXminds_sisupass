package handlers

import (
	"net/http"

	"sisupass.com/sisupass/cmd/api/app"
)

// Healthcheck endpoint for API health monitoring
// @Summary		Health check
// @Description	Check API health and system status
// @Tags			System
// @Produce		json
// @Success		200	{object}	map[string]interface{}	"API is healthy"
// @Router			/health [get]
func Healthcheck(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		env := map[string]any{
			"status": "available",
			"system_info": map[string]string{
				"environment": app.Config.Env,
				"version":     app.Config.Version,
			},
		}
		err := app.WriteJSON(w, http.StatusOK, env, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}
