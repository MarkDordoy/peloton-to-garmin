package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

// NewLogger returns a zerolog logger based on the conventions in a LoggingConfig
func NewLogger(logLevel string, pretty bool) zerolog.Logger {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	out := io.Writer(os.Stdout)
	if pretty {
		out = zerolog.ConsoleWriter{Out: out}
	}

	logger := zerolog.New(out).With().Timestamp().Logger()

	return logger.Level(level)
}
