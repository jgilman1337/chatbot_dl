package common

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var (
	CtxExceededErr = errors.New("the request took too long to complete, and has been halted prematurely")
)

// Struct wrapper for error handler parameters.
type ErrorHandlerParams struct {
	Msg    string
	Step   string
	logger *slog.Logger
}

// Builds a new error handler with parameters.
func NewErrHandlerParams(m string, step string, l *slog.Logger) ErrorHandlerParams {
	return ErrorHandlerParams{
		Msg:    m,
		Step:   step,
		logger: l,
	}
}

// Handler to catch and process errors, including timeout exceeded errors.
func LogErr(err error, p ErrorHandlerParams) error {
	//Ignore non-errors
	if err == nil {
		return err
	}

	//Format the error
	var efmt string
	if errors.Is(err, context.DeadlineExceeded) {
		efmt = fmt.Sprintf(
			"%s (step: %s)",
			CtxExceededErr, p.Step,
		)
	} else {
		efmt = fmt.Sprintf(
			"%s (step: %s)",
			err.Error(), p.Step,
		)
	}

	p.logger.Error(p.Msg + ": " + efmt)
	return err
}
