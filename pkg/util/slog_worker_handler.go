package util

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

var (
	WorkerIDKey     = "worker_id"
	ThreadIDKey     = "thread_id"
	ServiceIdentKey = "srv_ident"
)

// WorkerHandler is a custom slog.Handler that formats logs as:
//
//	2025/08/06 14:59:53 LEVEL <workerID::TID> Message
type WorkerHandler struct {
	level slog.Level
	out   io.Writer

	attrs []slog.Attr // store accumulated attributes here
}

// NewSimpleHandler creates a new SimpleHandler with options.
func NewWorkerHandler(out io.Writer, opts *slog.HandlerOptions) *WorkerHandler {
	var level slog.Level = slog.LevelInfo // default level
	if opts != nil && opts.Level != nil {
		level = opts.Level.Level()
	}
	return &WorkerHandler{
		level: level,
		out:   out,
	}
}

// Enabled returns true if the record's level is at least the handler's level.
func (h *WorkerHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle formats and writes the log record.
func (h *WorkerHandler) Handle(_ context.Context, r slog.Record) error {
	//Basics
	ts := r.Time.Format("2006/01/02 15:04:05")
	lvl := strings.ToUpper(r.Level.String())

	//Write the log message
	/*
		_, err := fmt.Fprintf(
			h.out, "%s %s %s%s\n",
			ts, lvl, buildPrefix(h.attrs), r.Message,
		)
	*/
	_, err := fmt.Fprintf(
		h.out, "%s %s %s\n",
		ts, lvl, r.Message,
	)
	return err
}

// WithAttrs returns a new handler with the additional attributes.
// For simplicity, returns the same handler (ignores attrs).
func (h *WorkerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Copy existing attrs plus the new ones
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	// Return a new handler instance with combined attrs
	return &WorkerHandler{
		level: h.level,
		out:   h.out,
		attrs: newAttrs,
	}
}

// WithGroup returns a new handler with the added group prefix.
// For simplicity, returns the same handler (ignores groups).
func (h *WorkerHandler) WithGroup(_ string) slog.Handler {
	return h
}

/*
func buildPrefix(attrs []slog.Attr) string {
	var workerID, threadID string

	//Extract attributes
	for _, attr := range attrs {
		fmt.Printf("current key: %s\n", attr.Key)
		switch attr.Key {
		case WorkerIDKey:
			if attr.Value.Kind() == slog.KindInt64 {
				workerID = strconv.Itoa(int(attr.Value.Int64()))
			}
		case ThreadIDKey:
			if attr.Value.Kind() == slog.KindString {
				threadID = attr.Value.String()
			}
		}
	}

	//No workerID means no prefix
	if workerID == "" {
		return ""
	}

	prefix := "<worker ID# " + workerID
	if threadID != "" {
		prefix += " :: " + threadID
	}
	prefix += "> "

	return prefix
}
*/
