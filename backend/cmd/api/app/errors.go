package app

import (
	"fmt"
	"net/http"
)

// ValidationError represents validation errors for request data
type ValidationError struct {
	Errors map[string]string
}

// Error implements the error interface for ValidationError
func (ve ValidationError) Error() string {
	var errStr string
	for field, msg := range ve.Errors {
		errStr += fmt.Sprintf("%s: %s\n", field, msg)
	}
	return errStr
}

// LogError logs an error with request context
func (app *Application) LogError(r *http.Request, err error) {
	app.Logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

// ErrorResponse sends a JSON error response
func (app *Application) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Envelope{"error": message}

	err := app.WriteJSON(w, status, env, nil)
	if err != nil {
		app.LogError(r, err)
		w.WriteHeader(500)
	}
}

// ServerErrorResponse sends a 500 Internal Server Error response
func (app *Application) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.LogError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

// NotFoundResponse sends a 404 Not Found response
func (app *Application) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.ErrorResponse(w, r, http.StatusNotFound, message)
}

// MethodNotAllowedResponse sends a 405 Method Not Allowed response
func (app *Application) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// BadRequestResponse sends a 400 Bad Request response
func (app *Application) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

// FailedValidationResponse sends a 422 Unprocessable Entity response
func (app *Application) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// EditConflictResponse sends a 409 Conflict response
func (app *Application) EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.ErrorResponse(w, r, http.StatusConflict, message)
}

// RateLimitExceededResponse sends a 429 Too Many Requests response
func (app *Application) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.ErrorResponse(w, r, http.StatusTooManyRequests, message)
}

// InvalidCredentialsResponse sends a 401 Unauthorized response for invalid credentials
func (app *Application) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

// InvalidAuthenticationTokenResponse sends a 401 Unauthorized response for invalid tokens
func (app *Application) InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	app.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

// AuthenticationRequiredResponse sends a 401 Unauthorized response when authentication is required
func (app *Application) AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

// InactiveAccountResponse sends a 403 Forbidden response for inactive accounts
func (app *Application) InactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated before you can access this resource"
	app.ErrorResponse(w, r, http.StatusForbidden, message)
}

// NotPermittedResponse sends a 403 Forbidden response for insufficient permissions
func (app *Application) NotPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account does not have the necessary permissions to access this resource"
	app.ErrorResponse(w, r, http.StatusForbidden, message)
}
