package usecases

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tyriis/rest-go/internal/domain"
)

type LockUseCase struct {
	lockRepo domain.LockRepository
}

func NewLockUseCase(lockRepo domain.LockRepository) *LockUseCase {
	return &LockUseCase{
		lockRepo: lockRepo,
	}
}

func (uc *LockUseCase) CreateLock(lockInput *domain.LockInput) (*domain.Lock, error) {
	// Check if lock exists
	existingLock, err := uc.lockRepo.Get(lockInput.Key)
	if err != nil {
		const msg = "LockUseCase - CreateLock - uc.lockRepo.Get > %w"
		return nil, fmt.Errorf(msg, err)
	}
	if existingLock != nil {
		const msg = "LockUseCase - CreateLock - lock '%s' exists"
		return nil, fmt.Errorf(msg, lockInput.Key)
	}

	// Parse duration
	duration, err := time.ParseDuration(lockInput.Duration)
	if err != nil {
		const msg = "LockUseCase - CreateLock - time.ParseDuration > %w"
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
	return uc.lockRepo.Set(lock.Key, string(lockValue), duration)
}

func (uc *LockUseCase) DeleteLock(key string) error {
	err := uc.lockRepo.Del(key)
	if err != nil {
		const msg = "LockUseCase.DeleteLock - uc.lockRepo.Del > %w"
		return fmt.Errorf(msg, err)
	}
	return nil
}

func (uc *LockUseCase) GetLock(key *string) (*domain.Lock, error) {
	lock, err := uc.lockRepo.Get(*key)
	if err != nil {
		const msg = "LockUseCase.GetLock - uc.lockRepo.Get > %w"
		return nil, fmt.Errorf(msg, err)
	}
	return lock, nil
}

func (uc *LockUseCase) ListLocks() ([]*domain.Lock, error) {
	locks, err := uc.lockRepo.Get("*")
	if err != nil {
		const msg = "LockUseCase.ListLocks - uc.lockRepo.Get > %w"
		return nil, fmt.Errorf(msg, err)
	}
	if locks == nil {
		return []*domain.Lock{}, nil
	}
	return []*domain.Lock{locks}, nil
}
