package common

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

// Return object for the error handler.
type ErrHandlerResult int8

const (
	EH_OK ErrHandlerResult = iota
	EH_TIMEOUT_EXCEEDED
	EH_OTHER_ERR
)

// Struct wrapper for error handler parameters.
type ErrorHandlerParams struct {
	Msg  string
	Step string
	Err  error

	FiredCtxErr *bool
	OutErr      *error
}

// Builds a new error handler with parameters.
func NewErrHandlerParams(m string, step string, err error, firedCtxErr *bool, outErr *error) ErrorHandlerParams {
	return ErrorHandlerParams{
		Msg:         m,
		Step:        step,
		Err:         err,
		FiredCtxErr: firedCtxErr,
		OutErr:      outErr,
	}
}

// Handler to catch and process errors, including timeout exceeded errors.
func HandleErr(l *slog.Logger, p ErrorHandlerParams) ErrHandlerResult {
	//Ignore non-errors
	if p.Err == nil {
		return EH_OK
	}

	//Handle the error correctly
	if errors.Is(p.Err, context.DeadlineExceeded) {
		//Workaround for capitalized error warning
		efmt := fmt.Sprintf(
			"Time's up! The request took too long to complete, and has been halted prematurely (step: %s); reason:",
			p.Step,
		)

		//Emit the error
		terr := fmt.Errorf("%s %w", efmt, p.Err)
		if !(*p.FiredCtxErr) {
			//Don't log this error twice and don't overwrite the original result
			l.Error(terr.Error())
			*p.OutErr = terr
		}

		*p.FiredCtxErr = true
		return EH_TIMEOUT_EXCEEDED

	} else {
		LogErr(l, p.Msg, p.Err)
		return EH_OTHER_ERR
	}
}

// Small utility to log errors to slog while preserving structure.
func LogErr(l *slog.Logger, m string, err error) {
	l.Error(m, slog.Any("error", err))
}
