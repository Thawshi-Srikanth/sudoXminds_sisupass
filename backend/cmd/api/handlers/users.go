package handlers

import (
	"errors"
	"net/http"

	appPkg "sisupass.com/sisupass/cmd/api/app"
	"sisupass.com/sisupass/internal/services"
	"sisupass.com/sisupass/internal/types"
	"sisupass.com/sisupass/internal/validator"
)

// RegisterUser registers a new user account
// @Summary		Register new user
// @Description	Register a new user account and send activation email
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			request	body		services.CreateUserRequest	true	"User registration data"
// @Success		201		{object}	map[string]interface{}		"User created successfully"
// @Failure		400		{object}	map[string]interface{}		"Bad request"
// @Failure		422		{object}	map[string]interface{}		"Validation failed"
// @Failure		500		{object}	map[string]interface{}		"Internal server error"
// @Router			/api/v1/users/register [post]
func RegisterUser(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input services.CreateUserRequest

		err := app.ReadJSON(w, r, &input)
		if err != nil {
			app.BadRequestResponse(w, r, err)
			return
		}

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

		activationToken, err := app.Services.Tokens.CreateActivationToken(r.Context(), userResponse.ID)
		if err != nil {
			app.LogError(r, err)
		} else {
			app.WG.Add(1)
			go func() {
				defer app.WG.Done()
				emailErr := app.Services.Mail.SendWelcomeEmail(r.Context(), userResponse.Email, userResponse.Username, activationToken.Plaintext)
				if emailErr != nil {
					app.Logger.PrintError(emailErr, map[string]string{
						"operation": "send_welcome_email",
						"user_id":   userResponse.ID.String(),
						"email":     userResponse.Email,
					})
				}
			}()
		}

		response := appPkg.Envelope{
			"user":    userResponse,
			"message": "User created successfully. Please check your email for activation instructions.",
		}

		err = app.WriteJSON(w, http.StatusCreated, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}

// GetUser retrieves a user by email
// @Summary		Get user by email
// @Description	Get user information by email address
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			email	query		string					true	"User email address"
// @Success		200		{object}	map[string]interface{}	"User retrieved successfully"
// @Failure		400		{object}	map[string]interface{}	"Bad request"
// @Failure		404		{object}	map[string]interface{}	"User not found"
// @Failure		422		{object}	map[string]interface{}	"Validation failed"
// @Failure		500		{object}	map[string]interface{}	"Internal server error"
// @Security		Bearer
// @Router			/api/v1/users [get]
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

// SendPasswordResetEmail sends a password reset email
// @Summary		Send password reset email
// @Description	Send a password reset email to the specified email address
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			request	body		map[string]string		true	"Email address"	example({"email": "user@example.com"})
// @Success		202		{object}	map[string]interface{}	"Password reset email sent"
// @Failure		400		{object}	map[string]interface{}	"Bad request"
// @Failure		404		{object}	map[string]interface{}	"User not found"
// @Failure		422		{object}	map[string]interface{}	"Validation failed"
// @Failure		500		{object}	map[string]interface{}	"Internal server error"
// @Router			/api/v1/auth/password-reset [post]
func SendPasswordResetEmail(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Email string `json:"email"`
		}

		err := app.ReadJSON(w, r, &input)
		if err != nil {
			app.BadRequestResponse(w, r, err)
			return
		}

		v := validator.New()
		validator.ValidateEmail(v, input.Email)
		if !v.Valid() {
			app.FailedValidationResponse(w, r, v.Errors)
			return
		}

		userResponse, err := app.Services.Users.GetUserByEmail(r.Context(), input.Email)
		if err != nil {
			switch {
			case errors.Is(err, types.ErrRecordNotFound):
				app.WriteJSON(w, http.StatusOK, appPkg.Envelope{
					"message": "If that email address is in our database, you will receive a password reset email shortly.",
				}, nil)
				return
			default:
				app.ServerErrorResponse(w, r, err)
				return
			}
		}

		resetToken, err := app.Services.Tokens.CreatePasswordResetToken(r.Context(), userResponse.ID)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
			return
		}

		app.WG.Add(1)
		go func() {
			defer app.WG.Done()
			emailErr := app.Services.Mail.SendPasswordResetEmail(r.Context(), userResponse.Email, userResponse.Username, resetToken.Plaintext)
			if emailErr != nil {
				app.Logger.PrintError(emailErr, map[string]string{
					"operation": "send_password_reset_email",
					"user_id":   userResponse.ID.String(),
					"email":     userResponse.Email,
				})
			}
		}()

		err = app.WriteJSON(w, http.StatusOK, appPkg.Envelope{
			"message": "If that email address is in our database, you will receive a password reset email shortly.",
		}, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}
