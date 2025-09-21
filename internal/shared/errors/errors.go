package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"error,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Predefined errors
var (
	ErrNotFound     = NewAppError(http.StatusNotFound, "Resource not found", nil)
	ErrBadRequest   = NewAppError(http.StatusBadRequest, "Bad request", nil)
	ErrUnauthorized = NewAppError(http.StatusUnauthorized, "Unauthorized", nil)
	ErrForbidden    = NewAppError(http.StatusForbidden, "Forbidden", nil)
	ErrConflict     = NewAppError(http.StatusConflict, "Conflict", nil)
	ErrInternal     = NewAppError(http.StatusInternalServerError, "Internal server error", nil)
)

// Validation errors
func NewValidationError(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, fmt.Sprintf("Validation error: %s", message), err)
}

// Database errors
func NewDatabaseError(operation string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, fmt.Sprintf("Database %s failed", operation), err)
}

// Business logic errors
func NewBusinessError(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, err)
}
