package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	appPkg "sisupass.com/sisupass/cmd/api/app"
	"sisupass.com/sisupass/internal/services"
	"sisupass.com/sisupass/internal/types"
	"sisupass.com/sisupass/internal/validator"
)

// CreateSlot creates a new slot for a user
// @Summary		Create a new slot
// @Description	Create a new slot for the authenticated user
// @Tags			Slots
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		types.Slot	true	"Slot creation data"
// @Success		201		{object}	map[string]interface{}	"Slot created successfully"
// @Failure		400		{object}	map[string]interface{}	"Bad request"
// @Failure		401		{object}	map[string]interface{}	"Authentication required"
// @Failure		422		{object}	map[string]interface{}	"Validation failed"
// @Failure		500		{object}	map[string]interface{}	"Internal server error"
// @Router			/api/v1/slots [post]
func CreateSlot(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			SlotType    string  `json:"slot_type"`
			Title       string  `json:"title"`
			Description *string `json:"description,omitempty"`
			Action      string  `json:"action"`
			Status      string  `json:"status,omitempty"`
			Fields      []byte  `json:"fields,omitempty"`
		}

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

		// Create slot object
		slot := &types.Slot{
			SlotType:    input.SlotType,
			Title:       input.Title,
			Description: input.Description,
			Action:      input.Action,
			Status:      input.Status,
			Fields:      input.Fields,
		}

		// Create slot using service
		createdSlot, err := app.Services.Slots.CreateSlot(r.Context(), user.ID, slot)
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
			"slot":    createdSlot,
			"message": "Slot created successfully",
		}

		err = app.WriteJSON(w, http.StatusCreated, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}

// GetSlot retrieves a specific slot by ID
// @Summary		Get slot by ID
// @Description	Retrieve a specific slot by its ID (only if owned by authenticated user)
// @Tags			Slots
// @Produce		json
// @Security		BearerAuth
// @Param			slot_id	path		string	true	"Slot ID"
// @Success		200		{object}	map[string]interface{}	"Slot retrieved successfully"
// @Failure		400		{object}	map[string]interface{}	"Bad request"
// @Failure		401		{object}	map[string]interface{}	"Authentication required"
// @Failure		403		{object}	map[string]interface{}	"Access denied"
// @Failure		404		{object}	map[string]interface{}	"Slot not found"
// @Failure		500		{object}	map[string]interface{}	"Internal server error"
// @Router			/api/v1/slots/{slot_id} [get]
func GetSlot(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user from context
		user := app.ContextGetUser(r)
		if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
			app.AuthenticationRequiredResponse(w, r)
			return
		}

		// Get slot ID from URL
		slotID := chi.URLParam(r, "id")
		if slotID == "" {
			app.BadRequestResponse(w, r, fmt.Errorf("slot ID is required"))
			return
		}

		// Get slot using service
		slot, err := app.Services.Slots.GetSlot(r.Context(), user.ID, slotID)
		if err != nil {
			switch err {
			case services.ErrSlotNotFound:
				app.NotFoundResponse(w, r)
			default:
				app.ServerErrorResponse(w, r, err)
			}
			return
		}

		response := appPkg.Envelope{
			"slot": slot,
		}

		err = app.WriteJSON(w, http.StatusOK, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}

// GetUserSlots retrieves all slots for the authenticated user
// @Summary		Get user slots
// @Description	Retrieve all slots for the authenticated user with optional filters
// @Tags			Slots
// @Produce		json
// @Security		BearerAuth
// @Param			slot_type	query	string	false	"Filter by slot type"
// @Param			status		query	string	false	"Filter by status"
// @Success		200			{array}	types.Slot	"Slots retrieved successfully"
// @Failure		401			{object}	map[string]interface{}	"Authentication required"
// @Failure		500			{object}	map[string]interface{}	"Internal server error"
// @Router			/api/v1/slots [get]
func GetUserSlots(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user from context
		user := app.ContextGetUser(r)
		if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
			app.AuthenticationRequiredResponse(w, r)
			return
		}

		// Get query parameters
		slotType := r.URL.Query().Get("slot_type")
		status := r.URL.Query().Get("status")

		// Get slots using service
		slots, err := app.Services.Slots.GetUserSlots(r.Context(), user.ID, slotType, status)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
			return
		}

		response := appPkg.Envelope{
			"slots": slots,
		}

		err = app.WriteJSON(w, http.StatusOK, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}

// UpdateSlot updates an existing slot
// @Summary		Update a slot
// @Description	Update an existing slot (only if owned by authenticated user)
// @Tags			Slots
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			slot_id	path		string		true	"Slot ID"
// @Param			request	body		types.Slot	true	"Updated slot data"
// @Success		200		{object}	map[string]interface{}	"Slot updated successfully"
// @Failure		400		{object}	map[string]interface{}	"Bad request"
// @Failure		401		{object}	map[string]interface{}	"Authentication required"
// @Failure		403		{object}	map[string]interface{}	"Access denied"
// @Failure		404		{object}	map[string]interface{}	"Slot not found"
// @Failure		422		{object}	map[string]interface{}	"Validation failed"
// @Failure		500		{object}	map[string]interface{}	"Internal server error"
// @Router			/api/v1/slots/{slot_id} [put]
func UpdateSlot(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user from context
		user := app.ContextGetUser(r)
		if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
			app.AuthenticationRequiredResponse(w, r)
			return
		}

		// Get slot ID from URL
		slotID := chi.URLParam(r, "id")
		if slotID == "" {
			app.BadRequestResponse(w, r, fmt.Errorf("slot ID is required"))
			return
		}

		// Parse request body
		var input struct {
			SlotType    string  `json:"slot_type,omitempty"`
			Title       string  `json:"title,omitempty"`
			Description *string `json:"description,omitempty"`
			Action      string  `json:"action,omitempty"`
			Status      string  `json:"status,omitempty"`
			Fields      []byte  `json:"fields,omitempty"`
		}

		err := app.ReadJSON(w, r, &input)
		if err != nil {
			app.BadRequestResponse(w, r, err)
			return
		}

		// Create update object
		updates := &types.Slot{
			SlotType:    input.SlotType,
			Title:       input.Title,
			Description: input.Description,
			Action:      input.Action,
			Status:      input.Status,
			Fields:      input.Fields,
		}

		// Update slot using service
		updatedSlot, err := app.Services.Slots.UpdateSlot(r.Context(), user.ID, slotID, updates)
		if err != nil {
			switch err {
			case services.ErrSlotNotFound:
				app.NotFoundResponse(w, r)
			default:
				if validator.IsValidationError(err) {
					validationErr := err.(validator.ValidationError)
					app.FailedValidationResponse(w, r, validationErr.Errors)
					return
				}
				app.ServerErrorResponse(w, r, err)
			}
			return
		}

		response := appPkg.Envelope{
			"slot":    updatedSlot,
			"message": "Slot updated successfully",
		}

		err = app.WriteJSON(w, http.StatusOK, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}

// DeleteSlot deletes a slot
// @Summary		Delete a slot
// @Description	Delete a slot by ID (only if owned by authenticated user)
// @Tags			Slots
// @Security		BearerAuth
// @Param			slot_id	path	string	true	"Slot ID"
// @Success		200		{object}	map[string]interface{}	"Slot deleted successfully"
// @Failure		400		{object}	map[string]interface{}	"Bad request"
// @Failure		401		{object}	map[string]interface{}	"Authentication required"
// @Failure		403		{object}	map[string]interface{}	"Access denied"
// @Failure		404		{object}	map[string]interface{}	"Slot not found"
// @Failure		500		{object}	map[string]interface{}	"Internal server error"
// @Router			/api/v1/slots/{slot_id} [delete]
func DeleteSlot(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user from context
		user := app.ContextGetUser(r)
		if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
			app.AuthenticationRequiredResponse(w, r)
			return
		}

		// Get slot ID from URL
		slotID := chi.URLParam(r, "id")
		if slotID == "" {
			app.BadRequestResponse(w, r, fmt.Errorf("slot ID is required"))
			return
		}

		// Delete slot using service
		err := app.Services.Slots.DeleteSlot(r.Context(), user.ID, slotID)
		if err != nil {
			switch err {
			case services.ErrSlotNotFound:
				app.NotFoundResponse(w, r)
			default:
				app.ServerErrorResponse(w, r, err)
			}
			return
		}

		response := appPkg.Envelope{
			"message": "Slot deleted successfully",
		}

		err = app.WriteJSON(w, http.StatusOK, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}

// SearchSlots searches for slots
// @Summary		Search slots
// @Description	Search for slots by query string for the authenticated user
// @Tags			Slots
// @Produce		json
// @Security		BearerAuth
// @Param			q	query	string	true	"Search query"
// @Success		200	{array}	types.Slot	"Search results"
// @Failure		400	{object}	map[string]interface{}	"Bad request"
// @Failure		401	{object}	map[string]interface{}	"Authentication required"
// @Failure		500	{object}	map[string]interface{}	"Internal server error"
// @Router			/api/v1/slots/search [get]
func SearchSlots(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user from context
		user := app.ContextGetUser(r)
		if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
			app.AuthenticationRequiredResponse(w, r)
			return
		}

		// Get search query
		query := r.URL.Query().Get("q")
		if query == "" {
			app.BadRequestResponse(w, r, fmt.Errorf("search query is required"))
			return
		}

		// Search slots using service
		slots, err := app.Services.Slots.SearchSlots(r.Context(), user.ID, query)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
			return
		}

		response := appPkg.Envelope{
			"slots": slots,
		}

		err = app.WriteJSON(w, http.StatusOK, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}

// ExecuteSlotAction executes a slot's action
// @Summary		Execute slot action
// @Description	Execute the action defined in a slot (only if owned by authenticated user)
// @Tags			Slots
// @Produce		json
// @Security		BearerAuth
// @Param			slot_id	path	string	true	"Slot ID"
// @Success		200		{object}	map[string]interface{}	"Action executed successfully"
// @Failure		400		{object}	map[string]interface{}	"Bad request"
// @Failure		401		{object}	map[string]interface{}	"Authentication required"
// @Failure		403		{object}	map[string]interface{}	"Access denied"
// @Failure		404		{object}	map[string]interface{}	"Slot not found"
// @Failure		500		{object}	map[string]interface{}	"Internal server error"
// @Router			/api/v1/slots/{slot_id}/execute [post]
func ExecuteSlotAction(app *appPkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user from context
		user := app.ContextGetUser(r)
		if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
			app.AuthenticationRequiredResponse(w, r)
			return
		}

		// Get slot ID from URL
		slotID := chi.URLParam(r, "id")
		if slotID == "" {
			app.BadRequestResponse(w, r, fmt.Errorf("slot ID is required"))
			return
		}

		// Execute slot action using service
		result, err := app.Services.Slots.ExecuteSlotAction(r.Context(), user.ID, slotID)
		if err != nil {
			switch err {
			case services.ErrSlotNotFound:
				app.NotFoundResponse(w, r)
			case services.ErrSlotNotActive:
				app.BadRequestResponse(w, r, err)
			default:
				app.ServerErrorResponse(w, r, err)
			}
			return
		}

		response := appPkg.Envelope{
			"result": result,
		}

		err = app.WriteJSON(w, http.StatusOK, response, nil)
		if err != nil {
			app.ServerErrorResponse(w, r, err)
		}
	}
}
