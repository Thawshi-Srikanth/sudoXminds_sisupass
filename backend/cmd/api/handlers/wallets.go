package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	appPkg "sisupass.com/sisupass/cmd/api/app"
	"sisupass.com/sisupass/internal/services"
	"sisupass.com/sisupass/internal/types"
	"sisupass.com/sisupass/internal/validator"
)

// CreateWallet creates a new wallet for a user
// @Summary		Create a new wallet
// @Description	Create a new wallet for the authenticated user
// @Tags			Wallets
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		services.CreateWalletRequest	true	"Wallet creation data"
// @Success		201		{object}	map[string]interface{}			"Wallet created successfully"
// @Failure		400		{object}	map[string]interface{}			"Bad request"
// @Failure		401		{object}	map[string]interface{}			"Authentication required"
// @Failure		422		{object}	map[string]interface{}			"Validation failed"
// @Failure		500		{object}	map[string]interface{}			"Internal server error"
// @Router			/api/v1/wallets [post]
func CreateWallet(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input services.CreateWalletRequest

		err := app.ReadJSON(w, r, &input)
		if err != nil {
			app.BadRequestResponse(w, r, err)
			return
		}

		// Get authenticated user from context
		user := app.ContextGetUser(r)
		if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
			app.AuthenticationRequiredResponse(w, r)
			return
		}

		walletResponse, err := app.Services.Wallets.CreateWallet(r.Context(), user.ID, input)
		if err != nil {
			switch {
			case validator.IsValidationError(err):
				validationErr := err.(validator.ValidationError)
				app.FailedValidationResponse(w, r, validationErr.Errors)
			default:
				app.ServerErrorResponse(w, r, err)
			}
			return
		}

		response := appPkg.Envelope{
			"wallet":  walletResponse,
			"message": "Wallet created successfully",
		}

		err = app.WriteJSON(w, http.StatusCreated, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}

// GetWallet retrieves a specific wallet by ID
// @Summary		Get wallet by ID
// @Description	Retrieve a specific wallet by its ID (only if owned by authenticated user)
// @Tags			Wallets
// @Produce		json
// @Security		BearerAuth
// @Param			wallet_id	path		string	true	"Wallet ID"
// @Success		200			{object}	map[string]interface{}	"Wallet retrieved successfully"
// @Failure		400			{object}	map[string]interface{}	"Bad request"
// @Failure		401			{object}	map[string]interface{}	"Authentication required"
// @Failure		403			{object}	map[string]interface{}	"Access denied"
// @Failure		404			{object}	map[string]interface{}	"Wallet not found"
// @Failure		500			{object}	map[string]interface{}	"Internal server error"
// @Router			/api/v1/wallets/{wallet_id} [get]
func GetWallet(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user from context
		user := app.ContextGetUser(r)
		if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
			app.AuthenticationRequiredResponse(w, r)
			return
		}

		walletID := chi.URLParam(r, "wallet_id")
		if walletID == "" {
			app.BadRequestResponse(w, r, errors.New("wallet ID required"))
			return
		}

		walletResponse, err := app.Services.Wallets.GetWallet(r.Context(), walletID)
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

		// Check if the wallet belongs to the authenticated user
		if walletResponse.UserID != user.ID {
			app.NotPermittedResponse(w, r)
			return
		}

		response := appPkg.Envelope{
			"wallet": walletResponse,
		}

		err = app.WriteJSON(w, http.StatusOK, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}

// GetUserWallets retrieves all wallets for the authenticated user
// @Summary		Get user wallets
// @Description	Retrieve all wallets for the authenticated user
// @Tags			Wallets
// @Produce		json
// @Security		BearerAuth
// @Success		200		{object}	map[string]interface{}	"Wallets retrieved successfully"
// @Failure		401		{object}	map[string]interface{}	"Authentication required"
// @Failure		404		{object}	map[string]interface{}	"User not found"
// @Failure		500		{object}	map[string]interface{}	"Internal server error"
// @Router			/api/v1/wallets [get]
func GetUserWallets(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user from context
		user := app.ContextGetUser(r)
		if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
			app.AuthenticationRequiredResponse(w, r)
			return
		}

		wallets, err := app.Services.Wallets.GetUserWallets(r.Context(), user.ID.String())
		if err != nil {
			switch {
			case errors.Is(err, types.ErrRecordNotFound):
				app.NotFoundResponse(w, r)
			default:
				app.ServerErrorResponse(w, r, err)
			}
			return
		}

		response := appPkg.Envelope{
			"wallets": wallets,
			"count":   len(wallets),
		}

		err = app.WriteJSON(w, http.StatusOK, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}

// UpdateWalletBalance updates the balance of a wallet
// @Summary		Update wallet balance
// @Description	Add or subtract amount from wallet balance (only if owned by authenticated user)
// @Tags			Wallets
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			wallet_id	path		string							true	"Wallet ID"
// @Param			request		body		services.UpdateBalanceRequest	true	"Balance update data"
// @Success		200			{object}	map[string]interface{}			"Balance updated successfully"
// @Failure		400			{object}	map[string]interface{}			"Bad request"
// @Failure		401			{object}	map[string]interface{}			"Authentication required"
// @Failure		403			{object}	map[string]interface{}			"Access denied"
// @Failure		404			{object}	map[string]interface{}			"Wallet not found"
// @Failure		422			{object}	map[string]interface{}			"Validation failed"
// @Failure		500			{object}	map[string]interface{}			"Internal server error"
// @Router			/api/v1/wallets/{wallet_id}/balance [patch]
func UpdateWalletBalance(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user from context
		user := app.ContextGetUser(r)
		if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
			app.AuthenticationRequiredResponse(w, r)
			return
		}

		walletID := chi.URLParam(r, "wallet_id")
		if walletID == "" {
			app.BadRequestResponse(w, r, errors.New("wallet ID required"))
			return
		}

		// First, verify that the wallet belongs to the authenticated user
		walletResponse, err := app.Services.Wallets.GetWallet(r.Context(), walletID)
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

		// Check if the wallet belongs to the authenticated user
		if walletResponse.UserID != user.ID {
			app.NotPermittedResponse(w, r)
			return
		}

		var input services.UpdateBalanceRequest

		err = app.ReadJSON(w, r, &input)
		if err != nil {
			app.BadRequestResponse(w, r, err)
			return
		}

		updatedWalletResponse, err := app.Services.Wallets.UpdateBalance(r.Context(), walletID, input)
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

		response := appPkg.Envelope{
			"wallet":  updatedWalletResponse,
			"message": "Balance updated successfully",
		}

		err = app.WriteJSON(w, http.StatusOK, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}

// DeleteWallet deletes a wallet by ID
// @Summary		Delete wallet
// @Description	Delete a wallet by its ID (only if owned by authenticated user)
// @Tags			Wallets
// @Produce		json
// @Security		BearerAuth
// @Param			wallet_id	path		string	true	"Wallet ID"
// @Success		200			{object}	map[string]interface{}	"Wallet deleted successfully"
// @Failure		400			{object}	map[string]interface{}	"Bad request"
// @Failure		401			{object}	map[string]interface{}	"Authentication required"
// @Failure		403			{object}	map[string]interface{}	"Access denied"
// @Failure		404			{object}	map[string]interface{}	"Wallet not found"
// @Failure		500			{object}	map[string]interface{}	"Internal server error"
// @Router			/api/v1/wallets/{wallet_id} [delete]
func DeleteWallet(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user from context
		user := app.ContextGetUser(r)
		if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
			app.AuthenticationRequiredResponse(w, r)
			return
		}

		walletID := chi.URLParam(r, "wallet_id")
		if walletID == "" {
			app.BadRequestResponse(w, r, errors.New("wallet ID required"))
			return
		}

		// First, verify that the wallet belongs to the authenticated user
		walletResponse, err := app.Services.Wallets.GetWallet(r.Context(), walletID)
		if err != nil {
			switch {
			case errors.Is(err, types.ErrRecordNotFound):
				app.NotFoundResponse(w, r)
			default:
				app.ServerErrorResponse(w, r, err)
			}
			return
		}

		// Check if the wallet belongs to the authenticated user
		if walletResponse.UserID != user.ID {
			app.NotPermittedResponse(w, r)
			return
		}

		err = app.Services.Wallets.DeleteWallet(r.Context(), walletID)
		if err != nil {
			switch {
			case errors.Is(err, types.ErrRecordNotFound):
				app.NotFoundResponse(w, r)
			default:
				app.ServerErrorResponse(w, r, err)
			}
			return
		}

		response := appPkg.Envelope{
			"message": "Wallet deleted successfully",
		}

		err = app.WriteJSON(w, http.StatusOK, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}
