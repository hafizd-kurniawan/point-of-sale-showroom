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
	Status    string      `json:"status"`
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Error     string      `json:"error"`
	ErrorCode string      `json:"error_code,omitempty"`
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

// SuccessResponse represents a success response
type SuccessResponse struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	HasMore    bool `json:"has_more"`
}

// PaginationResponse represents a paginated response
type PaginationResponse struct {
	Status    string         `json:"status"`
	Message   string         `json:"message"`
	Data      interface{}    `json:"data"`
	Meta      PaginationMeta `json:"meta"`
	Timestamp time.Time      `json:"timestamp"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Status:    "success",
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, err string, errorCode ...string) *ErrorResponse {
	response := &ErrorResponse{
		Status:    "error",
		Success:   false,
		Message:   message,
		Error:     err,
		Timestamp: time.Now(),
	}
	if len(errorCode) > 0 {
		response.ErrorCode = errorCode[0]
	}
	return response
}

// NewValidationErrorResponse creates a new validation error response
func NewValidationErrorResponse(message string, err string, errors interface{}, errorCode ...string) *ErrorResponse {
	response := &ErrorResponse{
		Status:    "error",
		Success:   false,
		Message:   message,
		Error:     err,
		Errors:    errors,
		Timestamp: time.Now(),
	}
	if len(errorCode) > 0 {
		response.ErrorCode = errorCode[0]
	}
	return response
}

// NewPaginationResponse creates a new pagination response
func NewPaginationResponse(message string, data interface{}, meta PaginationMeta) *PaginationResponse {
	return &PaginationResponse{
		Status:    "success",
		Message:   message,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now(),
	}
}