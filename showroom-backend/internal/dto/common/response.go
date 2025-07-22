package common

import (
	"time"
)

// APIResponse represents the standard API response format
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Error     string      `json:"error"`
	Errors    interface{} `json:"errors,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
	HasMore    bool           `json:"has_more"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}) *APIResponse {
	return &APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, err string) *ErrorResponse {
	return &ErrorResponse{
		Success:   false,
		Message:   message,
		Error:     err,
		Timestamp: time.Now(),
	}
}

// NewValidationErrorResponse creates a new validation error response
func NewValidationErrorResponse(message string, err string, errors interface{}) *ErrorResponse {
	return &ErrorResponse{
		Success:   false,
		Message:   message,
		Error:     err,
		Errors:    errors,
		Timestamp: time.Now(),
	}
}