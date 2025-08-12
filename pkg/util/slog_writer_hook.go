package util

import (
	"io"
	"log/slog"
)

// Enforces compliance with the Writer interface.
var _ io.Writer = (*SlogWriterHook)(nil)

// An adapter to slog that allows applications outputting to `io.Writer` to use slog.
type SlogWriterHook struct {
	logger  *slog.Logger
	logFunc func(msg string, args ...any)
}

// Creates a new SlogWriterHook instance
func NewSlogWriterHook(logger *slog.Logger, level slog.Level) *SlogWriterHook {
	//Pick the appropriate logging function
	var lfunc func(msg string, args ...any)
	switch level {
	case slog.LevelDebug:
		lfunc = logger.Debug
	case slog.LevelWarn:
		lfunc = logger.Warn
	case slog.LevelError:
		lfunc = logger.Error
	case slog.LevelInfo:
	default:
		lfunc = logger.Info
	}

	return &SlogWriterHook{
		logger:  logger,
		logFunc: lfunc,
	}
}

// Write implements io.Writer.
func (s *SlogWriterHook) Write(p []byte) (n int, err error) {
	s.logFunc(string(p))
	return len(p), nil
}
