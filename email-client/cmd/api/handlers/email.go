package handlers

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "sisupass.com/email-client/cmd/api/app"
    "sisupass.com/email-client/internal/services"
    "sisupass.com/email-client/internal/types"
    "sisupass.com/email-client/internal/validator"
)

type EmailHandler struct {
    mailerService *services.MailerService
}

func NewEmailHandler(mailerService *services.MailerService) *EmailHandler {
    return &EmailHandler{
        mailerService: mailerService,
    }
}

// SendEmail sends a single email
// @Summary		Send single email
// @Description	Send a single email using a specified template with dynamic data
// @Tags			Email Operations
// @Accept			json
// @Produce		json
// @Param			request	body		types.EmailRequest		true	"Email request with template and data"
// @Success		200		{object}	types.EmailResponse		"Email queued successfully"
// @Failure		400		{object}	app.ErrorResponse		"Bad request - invalid JSON or missing required fields"
// @Failure		422		{object}	app.ValidationResponse	"Validation failed - invalid email format or template"
// @Failure		500		{object}	app.ErrorResponse		"Internal server error - email service unavailable"
// @Router			/email/send [post]
func (h *EmailHandler) SendEmail(appPtr *app.Application) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var input types.EmailRequest

        err := appPtr.ReadJSON(w, r, &input)
        if err != nil {
            appPtr.BadRequestResponse(w, r, err)
            return
        }

        // Validate input
        v := validator.New()
        v.Check(input.To != "", "to", "must be provided")
        v.Check(validator.Matches(input.To, validator.EmailRx), "to", "must be a valid email address")
        v.Check(input.Template != "", "template", "must be provided")
        v.Check(validator.PermittedValue(input.Template,
            types.TemplateWelcome,
            types.TemplatePasswordReset,
            types.TemplateVerification,
            types.TemplateNotification,
            types.TemplateInvoice,
            types.TemplateBooking,
        ), "template", "must be a valid template")

        if input.Priority != "" {
            v.Check(validator.PermittedValue(input.Priority,
                types.PriorityHigh,
                types.PriorityNormal,
                types.PriorityLow,
            ), "priority", "must be high, normal, or low")
        } else {
            input.Priority = types.PriorityNormal
        }

        if !v.Valid() {
            appPtr.FailedValidationResponse(w, r, v.Errors)
            return
        }

        // Send email
        response, err := h.mailerService.SendEmail(&input)
        if err != nil {
            appPtr.ServerErrorResponse(w, r, err)
            return
        }

        err = appPtr.WriteJSON(w, http.StatusOK, app.Envelope{"email": response}, nil)
        if err != nil {
            appPtr.ServerErrorResponse(w, r, err)
        }
    }
}

// SendBulkEmail sends emails to multiple recipients
// @Summary		Send bulk emails
// @Description	Send the same email template to multiple recipients at once
// @Tags			Email Operations
// @Accept			json
// @Produce		json
// @Param			request	body		types.BulkEmailRequest	true	"Bulk email request with recipients list"
// @Success		200		{array}		types.EmailResponse		"Emails queued successfully"
// @Failure		400		{object}	app.ErrorResponse		"Bad request - invalid JSON"
// @Failure		422		{object}	app.ValidationResponse	"Validation failed - invalid emails or too many recipients"
// @Failure		500		{object}	app.ErrorResponse		"Internal server error"
// @Router			/email/bulk [post]
func (h *EmailHandler) SendBulkEmail(appPtr *app.Application) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var input types.BulkEmailRequest

        err := appPtr.ReadJSON(w, r, &input)
        if err != nil {
            appPtr.BadRequestResponse(w, r, err)
            return
        }

        // Validate input
        v := validator.New()
        v.Check(len(input.Recipients) > 0, "recipients", "must be provided")
        v.Check(len(input.Recipients) <= 100, "recipients", "cannot exceed 100 recipients")

        // Validate each email address
        for _, email := range input.Recipients {
            v.Check(validator.Matches(email, validator.EmailRx), "recipients",
                "all recipients must be valid email addresses")
        }

        v.Check(input.Template != "", "template", "must be provided")
        v.Check(validator.PermittedValue(input.Template,
            types.TemplateWelcome,
            types.TemplatePasswordReset,
            types.TemplateVerification,
            types.TemplateNotification,
            types.TemplateInvoice,
            types.TemplateBooking,
        ), "template", "must be a valid template")

        // Check for unique recipients
        v.Check(validator.Unique(input.Recipients), "recipients", "must not contain duplicates")

        if !v.Valid() {
            appPtr.FailedValidationResponse(w, r, v.Errors)
            return
        }

        // Send bulk emails
        responses, err := h.mailerService.SendBulkEmail(&input)
        if err != nil {
            appPtr.ServerErrorResponse(w, r, err)
            return
        }

        err = appPtr.WriteJSON(w, http.StatusOK, app.Envelope{"emails": responses}, nil)
        if err != nil {
            appPtr.ServerErrorResponse(w, r, err)
        }
    }
}

// GetEmailStatus gets the status of a sent email
// @Summary		Get email delivery status
// @Description	Retrieve the current delivery status and details of a specific email
// @Tags			Email Tracking
// @Produce		json
// @Param			id	path		string				true	"Email ID returned from send operation"
// @Success		200	{object}	types.EmailStatus	"Email status and delivery information"
// @Failure		400	{object}	app.ErrorResponse	"Bad request - missing or invalid email ID"
// @Failure		404	{object}	app.ErrorResponse	"Email not found - invalid ID or expired"
// @Failure		500	{object}	app.ErrorResponse	"Internal server error"
// @Router			/email/{id}/status [get]
func (h *EmailHandler) GetEmailStatus(appPtr *app.Application) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")
        if id == "" {
            appPtr.BadRequestResponse(w, r,
                &validator.ValidationError{Errors: map[string]string{"id": "must be provided"}})
            return
        }

        status, err := h.mailerService.GetEmailStatus(id)
        if err != nil {
            appPtr.NotFoundResponse(w, r)
            return
        }

        err = appPtr.WriteJSON(w, http.StatusOK, app.Envelope{"status": status}, nil)
        if err != nil {
            appPtr.ServerErrorResponse(w, r, err)
        }
    }
}

// GetEmailTemplates lists available email templates
// @Summary		List available email templates
// @Description	Get all available email templates with descriptions and supported priorities
// @Tags			Email Templates
// @Produce		json
// @Success		200	{object}	types.TemplateResponse	"Available templates and priorities"
// @Router			/email/templates [get]
func (h *EmailHandler) GetEmailTemplates(appPtr *app.Application) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        templates := map[string]string{
            types.TemplateWelcome:       "User welcome email with activation link",
            types.TemplatePasswordReset: "Password reset email with reset token",
            types.TemplateVerification:  "Email verification for account activation",
            types.TemplateNotification:  "General notification email",
            types.TemplateInvoice:       "Invoice and billing email",
            types.TemplateBooking:       "Booking confirmation email",
        }

        priorities := []string{
            types.PriorityHigh,
            types.PriorityNormal,
            types.PriorityLow,
        }

        err := appPtr.WriteJSON(w, http.StatusOK, app.Envelope{
            "templates":  templates,
            "priorities": priorities,
        }, nil)
        if err != nil {
            appPtr.ServerErrorResponse(w, r, err)
        }
    }
}

// GetEmailStats gets email service statistics
// @Summary		Get email service statistics
// @Description	Retrieve comprehensive statistics about email service usage and performance
// @Tags			Email Analytics
// @Produce		json
// @Success		200	{object}	types.EmailStats	"Email service statistics and metrics"
// @Router			/email/stats [get]
func (h *EmailHandler) GetEmailStats(appPtr *app.Application) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        stats := h.mailerService.GetStats()

        err := appPtr.WriteJSON(w, http.StatusOK, app.Envelope{"stats": stats}, nil)
        if err != nil {
            appPtr.ServerErrorResponse(w, r, err)
        }
    }
}