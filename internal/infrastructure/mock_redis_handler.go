package infrastructure

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// MockRedisHandler mocks the RedisHandler interface.
type MockRedisHandler struct {
	mock.Mock
}

func (m *MockRedisHandler) Set(key string, value string, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *MockRedisHandler) Get(key string) ([]string, error) {
	args := m.Called(key)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisHandler) Del(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockRedisHandler) GetMultiple(keys []string) ([]string, error) {
	args := m.Called(keys)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisHandler) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisHandler) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisHandler) Count() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}
