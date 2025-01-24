package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/tyriis/rest-go/internal/domain"
	"github.com/tyriis/rest-go/internal/usecases"
)

type WebserviceHandler struct {
	LockUseCase *usecases.LockUseCase
	logger      domain.Logger
}

func NewWebserviceHandler(lockUseCase *usecases.LockUseCase, logger domain.Logger) *WebserviceHandler {
	return &WebserviceHandler{
		LockUseCase: lockUseCase,
		logger:      logger,
	}
}

func (h WebserviceHandler) respondWithJSON(res http.ResponseWriter, data interface{}) {
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	encoder.SetIndent("", "  ")
	encoder.Encode(data)
}

/**
 * CreateLock creates a lock if not exists
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
		// TODO: distinguish between 500 and 409
		h.logger.Error(err.Error())
		http.Error(res, err.Error(), http.StatusConflict)
		return
	}

	h.respondWithJSON(res, lock)
	h.logger.Debug("WebserviceHandler.CreateLock - END")
}

/**
 * DeleteLock deletes a lock if exists
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
 * ShowOneLock returns a single lock if exists
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
 * ShowAllLocks returns all locks
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
