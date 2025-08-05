package common

import (
	"context"
	"errors"
	"log/slog"
)

//
//-- Types
//

type loggerKey struct{}

//
//-- Errors
//

var ErrNoLoggerInCtx = errors.New("failed to get slog instance; this context has no logger bound to it")

//
//-- With functions
//

// Adds an slog logger to a context.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

//
//-- From functions
//

// Pulls an slog logger out of a context, falling back to `slog.Default()` if none is present.
func LoggerFromCtx(ctx context.Context) *slog.Logger {
	logger, err := LoggerFromCtxE(ctx)
	if err != nil && errors.Is(err, ErrNoLoggerInCtx) {
		return slog.Default()
	}
	return logger
}

// Pulls an slog logger out of a context, erroring out if none is present.
func LoggerFromCtxE(ctx context.Context) (*slog.Logger, error) {
	logger, ok := ctx.Value(loggerKey{}).(*slog.Logger)
	if !ok || logger == nil {
		return nil, ErrNoLoggerInCtx
	}
	return logger, nil
}
