// Package domain defines core interfaces and types for the application.
package domain

// Logger defines the standard interface for logging operations.
// This interface abstracts the logging implementation to maintain clean architecture.
type Logger interface {
	// Debug logs a message at debug level
	Debug(msg string)
	// Info logs a message at info level
	Info(msg string)
	// Warn logs a message at warn level
	Warn(msg string)
	// Error logs a message at error level
	Error(msg string)
}
