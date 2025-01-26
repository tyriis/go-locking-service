package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tyriis/go-locking-service/internal/domain"
)

type KVStoreHandler interface {
	Get(key string) ([]string, error)
	Set(key string, value string, expiration time.Duration) error
	Del(key string) error
	Count() (int, error)
}

type LockRepository struct {
	handler KVStoreHandler
	logger  domain.Logger
}

func NewLockRepository(handler KVStoreHandler, logger domain.Logger) *LockRepository {
	return &LockRepository{
		handler: handler,
		logger:  logger,
	}
}

func (repo *LockRepository) Get(key string) ([]*domain.Lock, error) {
	repo.logger.Debug(fmt.Sprintf("LockRepository.Get(%s) - START", key))
	result, err := repo.handler.Get(key)
	if err != nil {
		const msg = "LockRepository.Get - repo.handler.Get > %w"
		return nil, fmt.Errorf(msg, err)
	}

	// not found
	if result == nil {
		return nil, nil
	}

	if len(result) != 1 && key != "*" {
		const msg = "LockRepository.Get - expected exactly one result, got %d"
		return nil, fmt.Errorf(msg, len(result))
	}

	// itterate over the result and unmarshal the locks
	var locks []*domain.Lock = make([]*domain.Lock, 0)
	for _, r := range result {
		var lock domain.Lock
		if err := json.Unmarshal([]byte(r), &lock); err != nil {
			const msg = "LockRepository.Get - json.Unmarshal(%s) > %s"
			repo.logger.Error(fmt.Sprintf(msg, r, err.Error()))
			// TODO: would be good to know what lock is invalid, so we can prompt do delete it or even auto fix it
			repo.logger.Warn("LockRepository.Get > skipping invalid lock, store contains corrupt data")
			continue
		}
		locks = append(locks, &lock)
	}

	// TODO: check if lock content is valid, warn or delete if not

	repo.logger.Debug(fmt.Sprintf("LockRepository.Get(%s) - END", key))
	return locks, nil
}

func (repo *LockRepository) Set(key string, value string, duration time.Duration) (*domain.Lock, error) {
	repo.logger.Debug(fmt.Sprintf("LockRepository.Set(%s) - START", key))
	if err := repo.handler.Set(key, value, duration); err != nil {
		const msg = "LockRepository.Set - repo.handler.Set > %w"
		return nil, fmt.Errorf(msg, err)
	}
	locks, err := repo.Get(key)
	if err != nil {
		return nil, err
	}
	repo.logger.Debug(fmt.Sprintf("LockRepository.Set(%s) - END", key))
	return locks[0], nil
}

func (repo *LockRepository) Del(key string) error {
	repo.logger.Debug(fmt.Sprintf("LockRepository.Del(%s) - START", key))
	if err := repo.handler.Del(key); err != nil {
		const msg = "LockRepository.Del - repo.handler.Del > %w"
		return fmt.Errorf(msg, err)
	}
	repo.logger.Debug(fmt.Sprintf("LockRepository.Del(%s) - END", key))
	return nil
}

func (repo *LockRepository) Count() (int, error) {
	repo.logger.Debug("LockRepository.Count - START")
	count, err := repo.handler.Count()
	if err != nil {
		return 0, err
	}
	repo.logger.Debug("LockRepository.Count - END")
	return count, nil
}
