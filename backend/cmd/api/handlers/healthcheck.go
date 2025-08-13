package handlers

import (
	"net/http"

	"sisupass.com/sisupass/cmd/api/app"
)

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
