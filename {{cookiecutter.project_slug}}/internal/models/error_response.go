package models

// ErrorResponse represents the standard error response format.
type ErrorResponse struct {
    Error string `json:"error" example:"Error message"`
}
