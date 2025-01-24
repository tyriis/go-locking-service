// Package usecases implements the application's business logic.
package usecases

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tyriis/rest-go/internal/domain"
)

// LockUseCase handles the business logic for lock management.
type LockUseCase struct {
	lockRepo domain.LockRepository
	logger   domain.Logger
}

// NewLockUseCase creates a new LockUseCase with the given repository and logger.
func NewLockUseCase(lockRepo domain.LockRepository, logger domain.Logger) *LockUseCase {
	return &LockUseCase{
		lockRepo: lockRepo,
		logger:   logger,
	}
}

// CreateLock creates a new lock if it doesn't exist.
func (uc *LockUseCase) CreateLock(lockInput *domain.LockInput) (*domain.Lock, error) {
	uc.logger.Debug("LockUseCase.CreateLock - START")
	// Check if lock exists
	existingLock, _ := uc.lockRepo.Get(lockInput.Key)
	if existingLock != nil {
		const msg = "LockUseCase.CreateLock(%s) >"
		return nil, &domain.LockConflictError{Message: fmt.Sprintf(msg, lockInput.Key)}
	}

	// Parse duration
	duration, err := time.ParseDuration(lockInput.Duration)
	if err != nil {
		const msg = "LockUseCase.CreateLock - time.ParseDuration > %s"
		return nil, &domain.InputError{Message: fmt.Sprintf(msg, err.Error())}
	}

	// Create lock structure
	now := time.Now().UTC()
	lock := &domain.Lock{
		Key:       lockInput.Key,
		Owner:     lockInput.Owner,
		Duration:  int64(duration.Seconds()),
		CreatedAt: now,
		ExpireAt:  now.Add(duration),
	}

	// Convert to JSON
	lockValue, err := json.Marshal(lock)
	if err != nil {
		const msg = "LockUseCase.CreateLock - json.Marshal > %s"
		return nil, &domain.InternalError{Message: fmt.Sprintf(msg, err.Error())}
	}

	// Set lock
	result, err := uc.lockRepo.Set(lock.Key, string(lockValue), duration)
	if err != nil {
		const msg = "LockUseCase.CreateLock - uc.lockRepo.Set > %s"
		return nil, &domain.InternalError{Message: fmt.Sprintf(msg, err.Error())}
	}
	const msg = "LockUseCase.CreateLock - Lock created: %s"
	uc.logger.Info(fmt.Sprintf(msg, lockInput.Key))
	uc.logger.Debug("LockUseCase.CreateLock - END")
	return result, nil
}

// DeleteLock removes an existing lock.
func (uc *LockUseCase) DeleteLock(key string) error {
	uc.logger.Debug("LockUseCase.DeleteLock - START")
	if key == "" {
		const msg = "LockUseCase.DeleteLock - key is empty >"
		return &domain.InputError{Message: msg}
	}
	err := uc.lockRepo.Del(key)
	if err != nil {
		const msg = "LockUseCase.DeleteLock - uc.lockRepo.Del > %s"
		return &domain.InternalError{Message: fmt.Sprintf(msg, err.Error())}
	}
	const msg = "LockUseCase.DeleteLock - Lock deleted: %s"
	uc.logger.Info(fmt.Sprintf(msg, key))
	uc.logger.Debug("LockUseCase.DeleteLock - END")
	return nil
}

// GetLock retrieves a specific lock by key.
func (uc *LockUseCase) GetLock(key *string) (*domain.Lock, error) {
	uc.logger.Debug("LockUseCase.GetLock - START")
	lock, err := uc.lockRepo.Get(*key)
	if lock == nil {
		const msg = "LockUseCase.GetLock - uc.lockRepo.Get(%s) >"
		return nil, &domain.NotFoundError{Message: fmt.Sprintf(msg, *key)}
	}
	if err != nil {
		const msg = "LockUseCase.GetLock - uc.lockRepo.Get > %s"
		return nil, &domain.InternalError{Message: fmt.Sprintf(msg, err.Error())}
	}
	uc.logger.Debug("LockUseCase.GetLock - END")
	return lock[0], nil
}

// ListLocks retrieves all existing locks.
func (uc *LockUseCase) ListLocks() ([]*domain.Lock, error) {
	uc.logger.Debug("LockUseCase.ListLocks - START")
	locks, err := uc.lockRepo.Get("*")
	if err != nil {
		const msg = "LockUseCase.ListLocks - uc.lockRepo.Get > %s"
		return nil, &domain.InternalError{Message: fmt.Sprintf(msg, err.Error())}
	}

	uc.logger.Debug("LockUseCase.ListLocks - END")
	return locks, nil
}
