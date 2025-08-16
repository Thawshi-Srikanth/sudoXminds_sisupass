package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

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
// @Router			/api/v1/auth/activate [put]
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
// @Router			/api/v1/auth/password [put]
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

// GoogleLogin initiates Google OAuth login
// @Summary		Initiate Google OAuth login
// @Description	Redirect user to Google OAuth consent screen
// @Tags			Authentication
// @Success		307		{string}	string	"Redirect to Google OAuth"
// @Router			/api/v1/auth/google/login [get]
func GoogleLogin(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := generateStateOauthCookie(w)

		app.Logger.PrintInfo("Setting oauthState cookie", map[string]string{
			"state": state,
		})

		url := app.Services.OAuth.InitiateGoogleOAuth(r.Context(), state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// GoogleCallback handles Google OAuth callback
// @Summary		Handle Google OAuth callback
// @Description	Process Google OAuth callback and authenticate user
// @Tags			Authentication
// @Param			state	query		string	true	"OAuth state parameter"
// @Param			code	query		string	true	"OAuth authorization code"
// @Success		307		{string}	string	"Redirect to frontend with auth token"
// @Failure		400		{object}	map[string]interface{}	"Bad request"
// @Failure		401		{object}	map[string]interface{}	"Invalid credentials"
// @Failure		500		{object}	map[string]interface{}	"Internal server error"
// @Router			/api/v1/auth/google/callback [get]
func GoogleCallback(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")

		// Verify state token
		cookie, err := r.Cookie("oauthstate")
		if err != nil {
			app.Logger.PrintInfo("Missing cookie", map[string]string{
				"received_state": state,
				"error":          err.Error(),
			})
			app.InvalidCredentialsResponse(w, r)
			return
		}

		if cookie.Value != state {
			app.Logger.PrintInfo("State token mismatch", map[string]string{
				"received_state": state,
				"cookie_state":   cookie.Value,
			})
			app.InvalidCredentialsResponse(w, r)
			return
		}

		// Process OAuth callback using service
		oauthRequest := services.GoogleOAuthRequest{
			State: state,
			Code:  code,
		}

		oauthResponse, err := app.Services.OAuth.ProcessGoogleOAuthCallback(r.Context(), oauthRequest)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{"message": "OAuth callback processing failed"})
			app.ServerErrorResponse(w, r, err)
			return
		}

		app.Logger.PrintInfo("Redirecting to frontend", map[string]string{
			"url":   oauthResponse.RedirectURL,
			"email": oauthResponse.User.Email,
		})

		http.Redirect(w, r, oauthResponse.RedirectURL, http.StatusTemporaryRedirect)
	}
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return state
}
