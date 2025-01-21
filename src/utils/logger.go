package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	format := os.Getenv("LOG_FORMAT")
	switch format {
	case "console":
		// Console logging (default)
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Stack().Logger()
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
}
