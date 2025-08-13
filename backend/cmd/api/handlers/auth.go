package handlers

import (
	"errors"
	"net/http"

	appPkg "sisupass.com/sisupass/cmd/api/app"
	"sisupass.com/sisupass/internal/services"
	"sisupass.com/sisupass/internal/types"
	"sisupass.com/sisupass/internal/validator"
)

// ActivateUser activates a user account using a token
// @Summary		Activate user account
// @Description	Activate a user account using the activation token sent via email
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			request	body		services.ActivateUserRequest	true	"Activation token"
// @Success		200		{object}	map[string]interface{}			"User activated successfully"
// @Failure		400		{object}	map[string]interface{}			"Bad request"
// @Failure		422		{object}	map[string]interface{}			"Validation failed"
// @Failure		409		{object}	map[string]interface{}			"Edit conflict"
// @Failure		500		{object}	map[string]interface{}			"Internal server error"
// @Router			/auth/activate [put]
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

// UpdateUserPassword updates a user's password using a reset token
// @Summary		Update user password
// @Description	Update a user's password using the password reset token sent via email
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			request	body		services.UpdatePasswordRequest	true	"Password reset data"
// @Success		200		{object}	map[string]interface{}			"Password updated successfully"
// @Failure		400		{object}	map[string]interface{}			"Bad request"
// @Failure		422		{object}	map[string]interface{}			"Validation failed"
// @Failure		409		{object}	map[string]interface{}			"Edit conflict"
// @Failure		500		{object}	map[string]interface{}			"Internal server error"
// @Router			/auth/password [put]
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
