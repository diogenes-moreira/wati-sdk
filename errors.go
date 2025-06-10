package wati

import (
	"fmt"
	"net/http"
)

// WATIError representa un error específico de la API de WATI
type WATIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// Error implementa la interfaz error
func (e *WATIError) Error() string {
	return fmt.Sprintf("WATI API Error %d: %s", e.Code, e.Message)
}

// IsRetryable indica si el error es reintentable
func (e *WATIError) IsRetryable() bool {
	return e.Code >= 500 || e.Code == 429
}

// IsAuthenticationError indica si es un error de autenticación
func (e *WATIError) IsAuthenticationError() bool {
	return e.Code == 401
}

// IsAuthorizationError indica si es un error de autorización
func (e *WATIError) IsAuthorizationError() bool {
	return e.Code == 403
}

// IsNotFoundError indica si es un error de recurso no encontrado
func (e *WATIError) IsNotFoundError() bool {
	return e.Code == 404
}

// IsRateLimitError indica si es un error de límite de velocidad
func (e *WATIError) IsRateLimitError() bool {
	return e.Code == 429
}

// IsServerError indica si es un error del servidor
func (e *WATIError) IsServerError() bool {
	return e.Code >= 500
}

// Errores predefinidos comunes
var (
	ErrInvalidToken = &WATIError{
		Code:    401,
		Message: "Invalid API token",
		Type:    "authentication",
	}
	
	ErrInsufficientPermissions = &WATIError{
		Code:    403,
		Message: "Insufficient permissions",
		Type:    "authorization",
	}
	
	ErrResourceNotFound = &WATIError{
		Code:    404,
		Message: "Resource not found",
		Type:    "not_found",
	}
	
	ErrMethodNotAllowed = &WATIError{
		Code:    405,
		Message: "Method not allowed - check if you're using the correct HTTP method (POST vs GET)",
		Type:    "method_not_allowed",
	}
	
	ErrRateLimitExceeded = &WATIError{
		Code:    429,
		Message: "Rate limit exceeded",
		Type:    "rate_limit",
	}
	
	ErrServerError = &WATIError{
		Code:    500,
		Message: "Internal server error",
		Type:    "server_error",
	}
	
	ErrInvalidRequest = &WATIError{
		Code:    400,
		Message: "Invalid request",
		Type:    "bad_request",
	}
	
	ErrInvalidPhoneNumber = &WATIError{
		Code:    400,
		Message: "Invalid WhatsApp phone number",
		Type:    "validation",
	}
	
	ErrTemplateNotFound = &WATIError{
		Code:    404,
		Message: "Template not found",
		Type:    "template_error",
	}
	
	ErrContactNotFound = &WATIError{
		Code:    404,
		Message: "Contact not found",
		Type:    "contact_error",
	}
)

// NewWATIError crea un nuevo error de WATI basado en el código de estado HTTP
func NewWATIError(statusCode int, message string) *WATIError {
	errorType := "unknown"
	
	switch statusCode {
	case http.StatusBadRequest:
		errorType = "bad_request"
	case http.StatusUnauthorized:
		errorType = "authentication"
	case http.StatusForbidden:
		errorType = "authorization"
	case http.StatusNotFound:
		errorType = "not_found"
	case http.StatusMethodNotAllowed:
		errorType = "method_not_allowed"
	case http.StatusTooManyRequests:
		errorType = "rate_limit"
	case http.StatusInternalServerError:
		errorType = "server_error"
	default:
		if statusCode >= 500 {
			errorType = "server_error"
		} else if statusCode >= 400 {
			errorType = "client_error"
		}
	}
	
	return &WATIError{
		Code:    statusCode,
		Message: message,
		Type:    errorType,
	}
}

// ValidationError representa un error de validación
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implementa la interfaz error
func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error for field '%s': %s", e.Field, e.Message)
}

// MultiValidationError representa múltiples errores de validación
type MultiValidationError struct {
	Errors []ValidationError `json:"errors"`
}

// Error implementa la interfaz error
func (e *MultiValidationError) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	return fmt.Sprintf("Multiple validation errors: %d errors", len(e.Errors))
}

// Add agrega un error de validación
func (e *MultiValidationError) Add(field, message string) {
	e.Errors = append(e.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors indica si hay errores de validación
func (e *MultiValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

// NetworkError representa un error de red
type NetworkError struct {
	Operation string
	Err       error
}

// Error implementa la interfaz error
func (e *NetworkError) Error() string {
	return fmt.Sprintf("Network error during %s: %v", e.Operation, e.Err)
}

// Unwrap retorna el error subyacente
func (e *NetworkError) Unwrap() error {
	return e.Err
}

// IsRetryable indica si el error de red es reintentable
func (e *NetworkError) IsRetryable() bool {
	// La mayoría de errores de red son reintentables
	return true
}

