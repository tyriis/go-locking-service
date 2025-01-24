// Package delivery implements the HTTP handlers for the REST API.
package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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

// respondWithJSON writes a JSON response with proper indentation.
func (h WebserviceHandler) respondWithJSON(res http.ResponseWriter, status int, payload interface{}) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	encoder := json.NewEncoder(res)
	encoder.SetIndent("", "  ")
	encoder.Encode(payload)
}

// respondWithError writes a JSON error response
func (h WebserviceHandler) respondWithError(res http.ResponseWriter, status int, message string) {
	h.respondWithJSON(res, status, domain.NewErrorResponse(status, message).Error)
}

/**
 * CreateLock handles POST requests to create a new lock.
 */
func (h WebserviceHandler) CreateLock(res http.ResponseWriter, req *http.Request) {
	h.logger.Debug("WebserviceHandler.CreateLock - START")
	var input domain.LockInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		h.handleError(res, err)
		return
	}

	if err := domain.ValidateLockInput(&input); err != nil {
		h.handleError(res, err)
		return
	}

	lock, err := h.LockUseCase.CreateLock(&input)
	if err != nil {
		h.handleError(res, err)
		return
	}

	h.respondWithJSON(res, http.StatusCreated, domain.NewSuccessResponse(lock).Data)
	h.logger.Debug("WebserviceHandler.CreateLock - END")
}

/**
 * DeleteLock handles DELETE requests to remove an existing lock.
 */
func (h WebserviceHandler) DeleteLock(res http.ResponseWriter, req *http.Request) {
	h.logger.Debug("WebserviceHandler.DeleteLock - START")
	vars := mux.Vars(req)
	key := vars["key"]

	if err := h.LockUseCase.DeleteLock(key); err != nil {
		h.handleError(res, err)
		return
	}

	h.respondWithJSON(res, http.StatusOK, domain.NewSuccessResponse(nil).Data)
	h.logger.Debug("WebserviceHandler.DeleteLock - END")
}

/**
 * ShowOneLock handles GET requests to retrieve a specific lock.
 */
func (h WebserviceHandler) ShowOneLock(res http.ResponseWriter, req *http.Request) {
	h.logger.Debug("WebserviceHandler.ShowOneLock - START")
	vars := mux.Vars(req)
	key := vars["key"]
	lock, err := h.LockUseCase.GetLock(&key)
	if err != nil {
		h.handleError(res, err)
		return
	}

	h.respondWithJSON(res, http.StatusOK, domain.NewSuccessResponse(lock).Data)
	h.logger.Debug("WebserviceHandler.ShowOneLock - END")
}

/**
 * ShowAllLocks handles GET requests to retrieve all locks.
 */
func (h WebserviceHandler) ShowAllLocks(res http.ResponseWriter, req *http.Request) {
	h.logger.Debug("WebserviceHandler.ShowAllLocks - START")
	locks, err := h.LockUseCase.ListLocks()
	if err != nil {
		h.handleError(res, err)
		return
	}

	h.respondWithJSON(res, http.StatusOK, domain.NewSuccessResponse(locks).Data)
	h.logger.Debug("WebserviceHandler.ShowAllLocks - END")
}

func (h WebserviceHandler) handleError(res http.ResponseWriter, err error) {
	h.logger.Error(err.Error())
	switch e := err.(type) {
	case *domain.LockConflictError:
		h.respondWithError(res, http.StatusConflict, "lock already exists!")
	case *domain.NotFoundError:
		h.respondWithError(res, http.StatusNotFound, "not found")
	case *domain.InputError:
		h.respondWithError(res, http.StatusBadRequest, e.Error())
	default:
		h.respondWithError(res, http.StatusInternalServerError, "An unexpected error occurred")
	}
}
