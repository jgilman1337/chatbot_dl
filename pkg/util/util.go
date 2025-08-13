package util

import (
	"fmt"
	"time"
)

// Mimics a ternary operator (from: https://stackoverflow.com/a/59375088).
func If[T any](cond bool, vtrue T, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

// A wrapper for slog that allows for string formatting (via sprintf) without using the named parameters feature of slog.
func LogFmt(lfunc func(msg string, args ...any), msg string, args ...any) {
	lfunc(
		fmt.Sprintf(msg, args...),
	)
}

// Returns a timestamp according to the format: `yyyyMMdd_HHmmss`, eg: `20250723_133723`.
func Timestamp() string {
	layout := "20060102_150405"
	return time.Now().Format(layout)
}
