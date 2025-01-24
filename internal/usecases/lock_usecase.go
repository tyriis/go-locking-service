package usecases

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tyriis/rest-go/internal/domain"
)

type LockUseCase struct {
	lockRepo domain.LockRepository
	logger   domain.Logger
}

func NewLockUseCase(lockRepo domain.LockRepository, logger domain.Logger) *LockUseCase {
	return &LockUseCase{
		lockRepo: lockRepo,
		logger:   logger,
	}
}

func (uc *LockUseCase) CreateLock(lockInput *domain.LockInput) (*domain.Lock, error) {
	uc.logger.Debug("LockUseCase.CreateLock - START")
	// Check if lock exists
	existingLock, err := uc.lockRepo.Get(lockInput.Key)
	if err == nil || existingLock != nil {
		const msg = "LockUseCase.CreateLock > lock '%s' exists!"
		return nil, fmt.Errorf(msg, lockInput.Key)
	}

	// Parse duration
	duration, err := time.ParseDuration(lockInput.Duration)
	if err != nil {
		const msg = "LockUseCase.CreateLock - time.ParseDuration > %w"
		return nil, fmt.Errorf(msg, err)
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
		const msg = "LockUseCase.CreateLock - json.Marshal > %w"
		return nil, fmt.Errorf(msg, err)
	}

	// Set lock
	uc.logger.Debug("LockUseCase.CreateLock - END")
	return uc.lockRepo.Set(lock.Key, string(lockValue), duration)
}

func (uc *LockUseCase) DeleteLock(key string) error {
	uc.logger.Debug("LockUseCase.DeleteLock - START")
	if key == "" {
		const msg = "LockUseCase.DeleteLock - key is empty"
		return fmt.Errorf(msg)
	}
	err := uc.lockRepo.Del(key)
	if err != nil {
		const msg = "LockUseCase.DeleteLock - uc.lockRepo.Del > %w"
		return fmt.Errorf(msg, err)
	}
	const msg = "LockUseCase.DeleteLock - Lock deleted: %s"
	uc.logger.Info(fmt.Sprintf(msg, key))
	uc.logger.Debug("LockUseCase.DeleteLock - END")
	return nil
}

func (uc *LockUseCase) GetLock(key *string) (*domain.Lock, error) {
	uc.logger.Debug("LockUseCase.GetLock - START")
	lock, err := uc.lockRepo.Get(*key)
	if err != nil {
		const msg = "LockUseCase.GetLock - uc.lockRepo.Get > %w"
		return nil, fmt.Errorf(msg, err)
	}
	uc.logger.Debug("LockUseCase.GetLock - END")
	return lock[0], nil
}

func (uc *LockUseCase) ListLocks() ([]*domain.Lock, error) {
	uc.logger.Debug("LockUseCase.ListLocks - START")
	locks, err := uc.lockRepo.Get("*")
	if err != nil {
		const msg = "LockUseCase.ListLocks - uc.lockRepo.Get > %w"
		return nil, fmt.Errorf(msg, err)
	}

	uc.logger.Debug("LockUseCase.ListLocks - END")
	return locks, nil
}
