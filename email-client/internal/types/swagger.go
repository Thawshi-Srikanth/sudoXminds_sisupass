package types

import "time"

// Swagger response models for documentation


// HealthResponse represents health check response
type HealthResponse struct {
    Status     string            `json:"status" example:"available"`
    SystemInfo map[string]string `json:"system_info" example:"{\"environment\": \"development\", \"version\": \"1.0.0\"}"`
}

// EmailStats represents email service statistics for Swagger docs
type EmailStats struct {
    TotalSent    int64     `json:"total_sent" example:"1250"`
    TotalFailed  int64     `json:"total_failed" example:"23"`
    TotalPending int64     `json:"total_pending" example:"5"`
    TotalQueued  int64     `json:"total_queued" example:"8"`
    QueueSize    int       `json:"queue_size" example:"3"`
    Workers      int       `json:"workers" example:"5"`
    StartTime    time.Time `json:"start_time" example:"2025-08-15T10:30:00Z"`
    LastActivity time.Time `json:"last_activity" example:"2025-08-15T15:45:30Z"`
}

// Add example tags to existing types for better Swagger docs
type EmailRequestSwagger struct {
    To          string                 `json:"to" example:"user@example.com" validate:"required,email"`
    Template    string                 `json:"template" example:"user_welcome" validate:"required"`
    Subject     string                 `json:"subject,omitempty" example:"Welcome to SisuPass"`
    Data        map[string]interface{} `json:"data" example:"{\"name\": \"John Doe\", \"activation_token\": \"abc123\"}"`
    Priority    string                 `json:"priority,omitempty" example:"normal" enums:"high,normal,low"`
    ScheduledAt *time.Time             `json:"scheduled_at,omitempty" example:"2025-08-15T16:00:00Z"`
}

type EmailResponseSwagger struct {
    ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
    Status    string    `json:"status" example:"queued" enums:"queued,pending,sent,failed,retrying"`
    Message   string    `json:"message" example:"Email queued for sending"`
    SentAt    time.Time `json:"sent_at" example:"2025-08-15T10:30:00Z"`
    Recipient string    `json:"recipient" example:"user@example.com"`
}

type BulkEmailRequestSwagger struct {
    Recipients []string               `json:"recipients" example:"[\"user1@example.com\", \"user2@example.com\"]" validate:"required,max=100"`
    Template   string                 `json:"template" example:"notification" validate:"required"`
    Subject    string                 `json:"subject,omitempty" example:"System Notification"`
    Data       map[string]interface{} `json:"data" example:"{\"message\": \"System maintenance scheduled\"}"`
}

type EmailStatusSwagger struct {
    ID         string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
    Status     string     `json:"status" example:"sent" enums:"pending,sent,failed,scheduled,retrying"`
    Recipient  string     `json:"recipient" example:"user@example.com"`
    Template   string     `json:"template" example:"user_welcome"`
    SentAt     *time.Time `json:"sent_at,omitempty" example:"2025-08-15T10:35:00Z"`
    FailReason string     `json:"fail_reason,omitempty" example:"SMTP connection failed"`
    CreatedAt  time.Time  `json:"created_at" example:"2025-08-15T10:30:00Z"`
}