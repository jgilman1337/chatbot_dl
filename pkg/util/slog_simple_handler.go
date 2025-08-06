package util

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// SimpleHandler is a custom slog.Handler that formats logs as:
//
//	2025/08/06 14:59:53 LEVEL Message
type SimpleHandler struct {
	level slog.Level
	out   *os.File
}

// NewSimpleHandler creates a new SimpleHandler with options.
func NewSimpleHandler(out *os.File, opts *slog.HandlerOptions) *SimpleHandler {
	var level slog.Level = slog.LevelInfo // default level
	if opts != nil && opts.Level != nil {
		level = opts.Level.Level()
	}
	return &SimpleHandler{
		level: level,
		out:   out,
	}
}

// Enabled returns true if the record's level is at least the handler's level.
func (h *SimpleHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle formats and writes the log record.
func (h *SimpleHandler) Handle(_ context.Context, r slog.Record) error {
	ts := r.Time.Format("2006/01/02 15:04:05")
	lvl := strings.ToUpper(r.Level.String())
	_, err := fmt.Fprintf(h.out, "%s %s %s\n", ts, lvl, r.Message)
	return err
}

// WithAttrs returns a new handler with the additional attributes.
// For simplicity, returns the same handler (ignores attrs).
func (h *SimpleHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns a new handler with the added group prefix.
// For simplicity, returns the same handler (ignores groups).
func (h *SimpleHandler) WithGroup(_ string) slog.Handler {
	return h
}
