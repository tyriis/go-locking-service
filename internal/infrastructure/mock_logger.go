// Package infrastructure provides concrete implementations of interfaces defined in domain.
package infrastructure

import (
	"github.com/tyriis/go-locking-service/internal/domain"
)

type MockLogger struct{}

// Add this at init time to ensure Logger implements domain.Logger
var _ domain.Logger = (*MockLogger)(nil)

// NewLogger creates a new Logger instance configured based on environment variables.
// It uses LOG_FORMAT to determine the output format (console or json).
// It uses LOG_LEVEL to set the minimum log level (debug, info, warn, error).
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

func (l *MockLogger) Debug(msg string) {
	// Do nothing
}

func (l *MockLogger) Info(msg string) {
	// Do nothing
}

func (l *MockLogger) Warn(msg string) {
	// Do nothing
}

func (l *MockLogger) Error(msg string) {
	// Do nothing
}
