package repositories

import (
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tyriis/go-locking-service/internal/domain"
)

// MockRedisHandler mocks the RedisHandler interface.
type MockLockRepository struct {
	mock.Mock
}

func (m *MockLockRepository) Get(key string) ([]*domain.Lock, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Lock), args.Error(1)
}

func (m *MockLockRepository) Set(key string, value string, duration time.Duration) (*domain.Lock, error) {
	args := m.Called(key, value, duration)
	return args.Get(0).(*domain.Lock), args.Error(1)
}

func (m *MockLockRepository) Del(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockLockRepository) Count() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}
