package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/tyriis/rest-go/internal/domain"
	"github.com/tyriis/rest-go/internal/infrastructure"
	"github.com/tyriis/rest-go/internal/usecases"
)

type WebserviceHandler struct {
	LockUseCase *usecases.LockUseCase
	logger      *infrastructure.Logger
}

func NewWebserviceHandler(lockUseCase *usecases.LockUseCase, logger *infrastructure.Logger) *WebserviceHandler {
	return &WebserviceHandler{
		LockUseCase: lockUseCase,
		logger:      logger,
	}
}

/**
 * CreateLock creates a lock if not exists
 */
func (h WebserviceHandler) CreateLock(res http.ResponseWriter, req *http.Request) {
	// Check if the request body is a valid JSON structure
	var input domain.LockInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		log.Error().Err(err).Msg(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the input
	if err := domain.ValidateLockInput(&input); err != nil {
		log.Error().Err(err).Msg(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	lock, err := h.LockUseCase.CreateLock(&input)
	if err != nil {
		// TODO: distinguish between 500 and 409
		log.Error().Err(err).Msg(err.Error())
		http.Error(res, err.Error(), http.StatusConflict)
		return
	}

	json.NewEncoder(res).Encode(lock)
}

/**
 * DeleteLock deletes a lock if exists
 */
func (h WebserviceHandler) DeleteLock(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	key := vars["key"]

	if err := h.LockUseCase.DeleteLock(key); err != nil {
		log.Error().Err(err).Msg(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

/**
 * ShowOneLock returns a single lock if exists
 */
func (h WebserviceHandler) ShowOneLock(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	key := vars["key"]
	lock, err := h.LockUseCase.GetLock(&key)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(res).Encode(lock)
}

/**
 * ShowAllLocks returns all locks
 */
func (h WebserviceHandler) ShowAllLocks(res http.ResponseWriter, req *http.Request) {
	locks, err := h.LockUseCase.ListLocks()
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(res).Encode(locks)
}
