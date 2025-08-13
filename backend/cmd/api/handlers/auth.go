package handlers

import (
	"errors"
	"net/http"

	appPkg "sisupass.com/sisupass/cmd/api/app"
	"sisupass.com/sisupass/internal/services"
	"sisupass.com/sisupass/internal/types"
	"sisupass.com/sisupass/internal/validator"
)

func ActivateUser(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input services.ActivateUserRequest

		err := app.ReadJSON(w, r, &input)
		if err != nil {
			app.BadRequestResponse(w, r, err)
			return
		}

		userResponse, err := app.Services.Users.ActivateUser(r.Context(), input)
		if err != nil {
			switch {
			case validator.IsValidationError(err):
				validationErr := err.(validator.ValidationError)
				app.FailedValidationResponse(w, r, validationErr.Errors)
			case errors.Is(err, types.ErrEditConflict):
				app.EditConflictResponse(w, r)
			default:
				app.ServerErrorResponse(w, r, err)
			}
			return
		}

		response := appPkg.Envelope{
			"user":    userResponse,
			"message": "Account activated successfully",
		}

		err = app.WriteJSON(w, http.StatusOK, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}

func UpdateUserPassword(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input services.UpdatePasswordRequest

		err := app.ReadJSON(w, r, &input)
		if err != nil {
			app.BadRequestResponse(w, r, err)
			return
		}

		err = app.Services.Users.UpdateUserPassword(r.Context(), input)
		if err != nil {
			switch {
			case validator.IsValidationError(err):
				validationErr := err.(validator.ValidationError)
				app.FailedValidationResponse(w, r, validationErr.Errors)
			case errors.Is(err, types.ErrEditConflict):
				app.EditConflictResponse(w, r)
			default:
				app.ServerErrorResponse(w, r, err)
			}
			return
		}

		response := appPkg.Envelope{"message": "your password was successfully reset"}

		err = app.WriteJSON(w, http.StatusOK, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}
