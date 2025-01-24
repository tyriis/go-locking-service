package infrastructure

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	logger zerolog.Logger
}

func NewLogger() *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	format := os.Getenv("LOG_FORMAT")
	switch format {
	case "console":
		// Console logging (default)
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Logger()
		log.Debug().Msg("LOG_FORMAT is set to CONSOLE")
	case "json":
		// JSON logging
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		log.Debug().Msg("LOG_FORMAT is set to JSON")
	default:
		// JSON logging
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		log.Warn().Msg("LOG_FORMAT is not set, defaulting to JSON")
	}
	return &Logger{logger: log.Logger}
}

func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}
