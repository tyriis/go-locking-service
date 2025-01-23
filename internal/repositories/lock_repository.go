package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tyriis/rest-go/internal/domain"
)

type KVStoreHandler interface {
	Get(key string) (string, error)
	Set(key string, value string, expiration time.Duration) error
	Del(key string) error
}

type LockRepository struct {
	handler KVStoreHandler
}

func NewLockRepository(handler KVStoreHandler) *LockRepository {
	return &LockRepository{
		handler: handler,
	}
}

func (repo *LockRepository) Get(key string) (*domain.Lock, error) {
	result, err := repo.handler.Get(key)
	if err != nil {
		const msg = "LockRepository.Get - repo.handler.Get > %w"
		return nil, fmt.Errorf(msg, err)
	}

	// check if additional case with multi items is possible

	var lock domain.Lock
	if err := json.Unmarshal([]byte(result), &lock); err != nil {
		const msg = "LockRepository.Get - repo.handler.Get > %w"
		return nil, fmt.Errorf(msg, err)
	}

	return &lock, nil
}

func (repo *LockRepository) Set(key string, value string, duration time.Duration) (*domain.Lock, error) {
	if err := repo.handler.Set(key, value, duration); err != nil {
		const msg = "LockRepository.Set - repo.handler.Set > %w"
		return nil, fmt.Errorf(msg, err)
	}
	return repo.Get(key)
}

func (repo *LockRepository) Del(key string) error {
	if err := repo.handler.Del(key); err != nil {
		const msg = "LockRepository.Del - repo.handler.Del > %w"
		return fmt.Errorf(msg, err)
	}
	return nil
}
