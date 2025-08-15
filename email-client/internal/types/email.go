package types

import "time"

type EmailRequest struct {
	To          string                 `json:"to"`
	Template    string                 `json:"template"`
	Subject     string                 `json:"subject,omitempty"`
	Data        map[string]interface{} `json:"data"`
	Priority    string                 `json:"priority,omitempty"` // high, normal, low
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
}

type EmailResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	SentAt    time.Time `json:"sent_at"`
	Recipient string    `json:"recipient"`
}

type BulkEmailRequest struct {
	Recipients []string               `json:"recipients"`
	Template   string                 `json:"template"`
	Subject    string                 `json:"subject,omitempty"`
	Data       map[string]interface{} `json:"data"`
}

type EmailStatus struct {
	ID         string     `json:"id"`
	Status     string     `json:"status"` // pending, sent, failed, scheduled
	Recipient  string     `json:"recipient"`
	Template   string     `json:"template"`
	SentAt     *time.Time `json:"sent_at,omitempty"`
	FailReason string     `json:"fail_reason,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

const (
	TemplateWelcome       = "user_welcome"
	TemplatePasswordReset = "password_reset"
	TemplateVerification  = "email_verification"
	TemplateNotification  = "notification"
	TemplateInvoice       = "invoice"
	TemplateBooking       = "booking_confirmation"
)

const (
	PriorityHigh   = "high"
	PriorityNormal = "normal"
	PriorityLow    = "low"
)

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request data"`
}

type ValidationResponse struct {
	Errors map[string]string `json:"errors" example:"{\"email\": \"must be a valid email address\", \"template\": \"must be provided\"}"`
}

type TemplateResponse struct {
	Templates  map[string]string `json:"templates" example:"{\"user_welcome\": \"User welcome email with activation link\"}"`
	Priorities []string          `json:"priorities" example:"[\"high\", \"normal\", \"low\"]"`
}

