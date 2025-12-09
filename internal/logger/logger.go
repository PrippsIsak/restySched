package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes the global logger with pretty console output for development
func Init(debug bool) {
	// Use pretty console writer for human-readable logs
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	log.Logger = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	// Set log level
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

// Get returns the global logger
func Get() *zerolog.Logger {
	return &log.Logger
}

// WithContext returns a logger with additional context
func WithContext(key, value string) zerolog.Logger {
	return log.With().Str(key, value).Logger()
}
