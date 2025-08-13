package handlers

import (
	"errors"
	"net/http"

	appPkg "sisupass.com/sisupass/cmd/api/app"
	"sisupass.com/sisupass/internal/services"
	"sisupass.com/sisupass/internal/types"
	"sisupass.com/sisupass/internal/validator"
)

func RegisterUser(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input services.CreateUserRequest

		// 1. Parse request
		err := app.ReadJSON(w, r, &input)
		if err != nil {
			app.BadRequestResponse(w, r, err)
			return
		}

		// 2. Call service (all business logic handled here)
		userResponse, err := app.Services.Users.CreateUser(r.Context(), input)
		if err != nil {
			switch {
			case validator.IsValidationError(err):
				validationErr := err.(validator.ValidationError)
				app.FailedValidationResponse(w, r, validationErr.Errors)
			case errors.Is(err, types.ErrDuplicateEmail):
				app.FailedValidationResponse(w, r, map[string]string{
					"email": "a user with this email address already exists",
				})
			default:
				app.ServerErrorResponse(w, r, err)
			}
			return
		}

		// 3. Send response
		err = app.WriteJSON(w, http.StatusCreated, appPkg.Envelope{
			"user": userResponse,
		}, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}

func GetUser(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			app.BadRequestResponse(w, r, errors.New("email parameter is required"))
			return
		}

		userResponse, err := app.Services.Users.GetUserByEmail(r.Context(), email)
		if err != nil {
			switch {
			case validator.IsValidationError(err):
				validationErr := err.(validator.ValidationError)
				app.FailedValidationResponse(w, r, validationErr.Errors)
			case errors.Is(err, types.ErrRecordNotFound):
				app.NotFoundResponse(w, r)
			default:
				app.ServerErrorResponse(w, r, err)
			}
			return
		}

		err = app.WriteJSON(w, http.StatusOK, appPkg.Envelope{
			"user": userResponse,
		}, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}
