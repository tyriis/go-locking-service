// Package delivery implements the HTTP handlers for the REST API.
package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
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
func (h WebserviceHandler) respondWithJSON(res http.ResponseWriter, data interface{}) {
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	encoder.SetIndent("", "  ")
	encoder.Encode(data)
}

/**
 * CreateLock handles POST requests to create a new lock.
 */
func (h WebserviceHandler) CreateLock(res http.ResponseWriter, req *http.Request) {
	// Check if the request body is a valid JSON structure
	h.logger.Debug("WebserviceHandler.CreateLock - START")
	var input domain.LockInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		log.Error().Err(err).Msg(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the input
	if err := domain.ValidateLockInput(&input); err != nil {
		h.logger.Error(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	lock, err := h.LockUseCase.CreateLock(&input)
	if err != nil {
		switch e := err.(type) {
		case *domain.LockConflictError:
			h.logger.Error(e.Error())
			http.Error(res, e.Error(), http.StatusConflict)
		case *domain.InternalError:
			h.logger.Error(e.Error())
			http.Error(res, e.Error(), http.StatusInternalServerError)
		case *domain.InputError:
			h.logger.Error(e.Error())
			http.Error(res, e.Error(), http.StatusBadRequest)
		default:
			h.logger.Error(err.Error())
			http.Error(res, "unexpected error", http.StatusInternalServerError)
		}
		return
	}

	h.respondWithJSON(res, lock)
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
		h.logger.Error(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
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
		h.logger.Error(err.Error())
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	h.respondWithJSON(res, lock)
	h.logger.Debug("WebserviceHandler.ShowOneLock - END")
}

/**
 * ShowAllLocks handles GET requests to retrieve all locks.
 */
func (h WebserviceHandler) ShowAllLocks(res http.ResponseWriter, req *http.Request) {
	h.logger.Debug("WebserviceHandler.ShowAllLocks - START")
	locks, err := h.LockUseCase.ListLocks()
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	h.respondWithJSON(res, locks)
	h.logger.Debug("WebserviceHandler.ShowAllLocks - END")
}
