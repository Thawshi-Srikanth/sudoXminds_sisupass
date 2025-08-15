package app

// Swagger documentation response types

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request data"`
}

// ValidationResponse represents validation error response
type ValidationResponse struct {
	Errors map[string]string `json:"errors" example:"{\"email\": \"must be a valid email address\", \"template\": \"must be provided\"}"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status     string            `json:"status" example:"available"`
	SystemInfo map[string]string `json:"system_info" example:"{\"environment\": \"development\", \"version\": \"1.0.0\"}"`
}
