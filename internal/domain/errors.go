package domain

import "fmt"

// InputError represents an error when the input is invalid
type InputError struct {
	Message string
}

func (e *InputError) Error() string {
	msg := fmt.Sprintf("%s invalid input!", e.Message)
	return msg
}

// NotFoundError represents an error when the resource is not found
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	msg := fmt.Sprintf("%s not found!", e.Message)
	return msg
}

// LockConflictError represents an error when a lock already exists
type LockConflictError struct {
	Message string
}

func (e *LockConflictError) Error() string {
	msg := fmt.Sprintf("%s lock already exists!", e.Message)
	return msg
}

// InternalError represents an internal server error
type InternalError struct {
	Message string
}

func (e *InternalError) Error() string {
	return "internal server error: " + e.Message
}

// APIResponse represents a standardized API response
type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

// APIError represents a standardized error structure
type APIError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// NewErrorResponse creates a new API error response
func NewErrorResponse(status int, message string) *APIResponse {
	return &APIResponse{
		Error: &APIError{
			Message: message,
			Status:  status,
		},
	}
}

// NewSuccessResponse creates a new API success response
func NewSuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Data: data,
	}
}
