package domain

import (
	"fmt"
	"time"
)

type Config struct {
	Redis struct {
		Host   string `yaml:"host"`
		Port   string `yaml:"path"`
		Prefix string `yaml:"keyPrefix"`
	} `yaml:"redis"`
	Api struct {
		Port string `yaml:"port"`
	} `yaml:"api"`
}

type ConfigRepository interface {
	Load() (*Config, error)
}

type ConfigValidator interface {
	Validate(data interface{}) error
}

type Lock struct {
	Key       string    `json:"key"`
	Owner     string    `json:"owner"`
	Duration  int64     `json:"duration"`
	ExpireAt  time.Time `json:"expireAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type LockError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type LockRepository interface {
	Get(key string) ([]*Lock, error)
	Set(key string, value string, ttl time.Duration) (*Lock, error)
	Del(key string) error
}

type ValidationError struct {
	Operation string
	Message   string
}

func (e *ValidationError) Error() string {
	return e.Operation + ": " + e.Message
}

func NewValidationError(operation string, message string) *ValidationError {
	return &ValidationError{
		Operation: operation,
		Message:   message,
	}
}

type LockInput struct {
	Key      string `json:"key"`
	Owner    string `json:"owner"`
	Duration string `json:"duration"`
}

func ValidateLockInput(input *LockInput) error {
	// input.Key need to be a string, not empty and minimum 3 character
	if err := ValidateLockKeyInput(&input.Key); err != nil {
		return err
	}
	// input.Owner need to be a string, not empty
	if input.Owner == "" {
		return NewValidationError("LOCK_REQUIRES_OWNER", "owner is required")
	}
	// input.Duration need to be a string and valid duration as timestring f.e. 1h20m
	if _, err := time.ParseDuration(input.Duration); err != nil {
		return NewValidationError("LOCK_INVALID_DURATION", fmt.Sprintf("duration '%s' is invalid", input.Duration))
	}
	return nil
}

func ValidateLockKeyInput(key *string) error {
	// input.Key need to be a string, not empty and minimum 3 character
	if len(*key) < 3 {
		return NewValidationError("LOCK_INVALID_KEY", fmt.Sprintf("key '%s' is invalid", *key))
	}
	return nil
}
