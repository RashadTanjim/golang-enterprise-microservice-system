package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	apperrors "enterprise-microservice-system/common/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorInfo represents error information in response
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Meta represents metadata for paginated responses
type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
	TotalCount int64 `json:"total_count"`
}

// Success sends a successful response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMeta sends a successful response with pagination metadata
func SuccessWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Created sends a created response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// Error sends an error response based on AppError
func Error(c *gin.Context, err error) {
	var appErr *apperrors.AppError

	switch {
	case errors.As(err, &appErr):
		// Use provided AppError
	case errors.As(err, new(validator.ValidationErrors)):
		appErr = apperrors.NewValidation(formatValidationError(err))
	case errors.As(err, new(*json.SyntaxError)):
		appErr = apperrors.NewBadRequest("invalid JSON payload")
	case errors.As(err, new(*json.UnmarshalTypeError)):
		appErr = apperrors.NewBadRequest(formatTypeError(err))
	case errors.As(err, new(*strconv.NumError)):
		appErr = apperrors.NewBadRequest("invalid numeric parameter")
	default:
		appErr = apperrors.NewInternal("internal server error", err)
	}

	statusCode := getStatusCode(appErr.Code)
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}

// getStatusCode maps error codes to HTTP status codes
func getStatusCode(code string) int {
	switch code {
	case apperrors.ErrCodeNotFound:
		return http.StatusNotFound
	case apperrors.ErrCodeBadRequest, apperrors.ErrCodeValidation:
		return http.StatusBadRequest
	case apperrors.ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case apperrors.ErrCodeForbidden:
		return http.StatusForbidden
	case apperrors.ErrCodeConflict:
		return http.StatusConflict
	case apperrors.ErrCodeRateLimit:
		return http.StatusTooManyRequests
	case apperrors.ErrCodeCircuitOpen, apperrors.ErrCodeServiceUnavail:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func formatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return "validation error"
	}

	messages := make([]string, 0, len(validationErrors))
	for _, fieldErr := range validationErrors {
		field := strings.ToLower(fieldErr.Field())
		switch fieldErr.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", field))
		case "min":
			messages = append(messages, fmt.Sprintf("%s must be at least %s", field, fieldErr.Param()))
		case "max":
			messages = append(messages, fmt.Sprintf("%s must be at most %s", field, fieldErr.Param()))
		case "email":
			messages = append(messages, fmt.Sprintf("%s must be a valid email", field))
		case "oneof":
			messages = append(messages, fmt.Sprintf("%s must be one of [%s]", field, fieldErr.Param()))
		case "url":
			messages = append(messages, fmt.Sprintf("%s must be a valid URL", field))
		default:
			messages = append(messages, fmt.Sprintf("%s is invalid", field))
		}
	}

	if len(messages) == 0 {
		return "validation error"
	}

	return strings.Join(messages, "; ")
}

func formatTypeError(err error) string {
	var typeErr *json.UnmarshalTypeError
	if !errors.As(err, &typeErr) {
		return "invalid request payload"
	}

	field := strings.ToLower(typeErr.Field)
	if field == "" {
		return "invalid request payload"
	}

	return fmt.Sprintf("%s has an invalid type", field)
}
