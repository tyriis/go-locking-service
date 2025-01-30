package usecases

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tyriis/go-locking-service/internal/domain"
	"github.com/tyriis/go-locking-service/internal/infrastructure"
	"github.com/tyriis/go-locking-service/internal/repositories"
)

const (
	testOwnerValue = "test-owner"
	testKeyValue   = "test-lock"
)

var testLock = &domain.Lock{
	Key:       testKeyValue,
	Owner:     testOwnerValue,
	Duration:  3600,
	CreatedAt: time.Now(),
	ExpireAt:  time.Now().Add(time.Hour),
}

func TestCreateLockSuccess(t *testing.T) {
	// Arrange
	testDuration := "1h"
	mockLogger := infrastructure.NewMockLogger()
	mockRepo := new(repositories.MockLockRepository)
	mockRepo.On("Set", testKeyValue, mock.AnythingOfType("string"), mock.Anything).Return(testLock, nil)
	mockRepo.On("Get", testKeyValue).Return(nil, nil)

	uc := NewLockUseCase(mockRepo, mockLogger)

	input := &domain.LockInput{
		Key:      testKeyValue,
		Owner:    testOwnerValue,
		Duration: testDuration,
	}

	// Act
	result, err := uc.CreateLock(input)

	// Assert
	mockRepo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testKeyValue, result.Key)
	assert.Equal(t, testOwnerValue, result.Owner)
	assert.Equal(t, int64(3600), result.Duration)
	assert.False(t, result.CreatedAt.IsZero())
	assert.False(t, result.ExpireAt.IsZero())
	assert.True(t, result.ExpireAt.After(result.CreatedAt))
}

func TestCreateLockConflict(t *testing.T) {
	// Arrange
	mockLogger := infrastructure.NewMockLogger()
	mockRepo := new(repositories.MockLockRepository)
	mockRepo.On("Get", testKeyValue).Return([]*domain.Lock{testLock}, nil)

	uc := NewLockUseCase(mockRepo, mockLogger)

	input := &domain.LockInput{
		Key:      testKeyValue,
		Owner:    testOwnerValue,
		Duration: "1h",
	}

	// Act
	result, err := uc.CreateLock(input)

	// Assert
	mockRepo.AssertExpectations(t)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, &domain.LockConflictError{}, err)

}

func TestDeleteLockSuccess(t *testing.T) {
	// Arrange
	mockLogger := infrastructure.NewMockLogger()
	mockRepo := new(repositories.MockLockRepository)
	mockRepo.On("Del", testKeyValue).Return(nil)

	uc := NewLockUseCase(mockRepo, mockLogger)

	// Act
	err := uc.DeleteLock(testKeyValue)

	// Assert
	mockRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestDeleteLockEmptyKey(t *testing.T) {
	// Arrange
	mockLogger := infrastructure.NewMockLogger()
	mockRepo := new(repositories.MockLockRepository)

	uc := NewLockUseCase(mockRepo, mockLogger)

	// Act
	err := uc.DeleteLock("")

	// Assert
	mockRepo.AssertExpectations(t)
	assert.Error(t, err)
	assert.IsType(t, &domain.InputError{}, err)

}

func TestGetLockNotFound(t *testing.T) {
	// Arrange
	mockLogger := infrastructure.NewMockLogger()
	mockRepo := new(repositories.MockLockRepository)
	mockRepo.On("Get", testKeyValue).Return([]*domain.Lock{nil}, nil)

	uc := NewLockUseCase(mockRepo, mockLogger)

	// Act
	result, err := uc.GetLock(testKeyValue)

	// Assert
	mockRepo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Nil(t, result)
}
