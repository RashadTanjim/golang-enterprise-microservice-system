package errors

import "fmt"

// AppError represents an application error with code and message
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Common error codes
const (
	ErrCodeInternal       = "INTERNAL_ERROR"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeBadRequest     = "BAD_REQUEST"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeDatabase       = "DATABASE_ERROR"
	ErrCodeCircuitOpen    = "CIRCUIT_OPEN"
	ErrCodeRateLimit      = "RATE_LIMIT_EXCEEDED"
	ErrCodeServiceUnavail = "SERVICE_UNAVAILABLE"
)

// New creates a new AppError
func New(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewInternal creates an internal error
func NewInternal(message string, err error) *AppError {
	return &AppError{
		Code:    ErrCodeInternal,
		Message: message,
		Err:     err,
	}
}

// NewNotFound creates a not found error
func NewNotFound(resource string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

// NewBadRequest creates a bad request error
func NewBadRequest(message string) *AppError {
	return &AppError{
		Code:    ErrCodeBadRequest,
		Message: message,
	}
}

// NewValidation creates a validation error
func NewValidation(message string) *AppError {
	return &AppError{
		Code:    ErrCodeValidation,
		Message: message,
	}
}

// NewConflict creates a conflict error
func NewConflict(message string) *AppError {
	return &AppError{
		Code:    ErrCodeConflict,
		Message: message,
	}
}

// NewCircuitOpen creates a circuit breaker open error
func NewCircuitOpen(service string) *AppError {
	return &AppError{
		Code:    ErrCodeCircuitOpen,
		Message: fmt.Sprintf("circuit breaker open for service: %s", service),
	}
}

// NewRateLimit creates a rate limit error
func NewRateLimit() *AppError {
	return &AppError{
		Code:    ErrCodeRateLimit,
		Message: "rate limit exceeded, please try again later",
	}
}