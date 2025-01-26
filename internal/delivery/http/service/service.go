// Package service implements the HTTP handlers for the REST API.
package service

import (
	"github.com/tyriis/rest-go/internal/domain"
	"github.com/tyriis/rest-go/internal/usecases"
)

// WebserviceHandler handles HTTP requests for the lock management API.
type WebserviceHandler struct {
	LockUseCase *usecases.LockUseCase
	logger      domain.Logger
}

// NewWebserviceHandler creates a new WebserviceHandler with the given use case and logger.
func NewWebserviceHandler(lockUseCase *usecases.LockUseCase, logger domain.Logger) *WebserviceHandler {
	return &WebserviceHandler{
		LockUseCase: lockUseCase,
		logger:      logger,
	}
}
