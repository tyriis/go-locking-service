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
