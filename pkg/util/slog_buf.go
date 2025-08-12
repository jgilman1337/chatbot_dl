package util

import (
	"bytes"
	"log/slog"
	"sync"
)

// Implements a thread-safe buffer for use with `slog`.
type SafeBuffer struct {
	b bytes.Buffer
	m sync.Mutex
}

// Returns the contents of the buffer as a []byte.
func (buf *SafeBuffer) Bytes() []byte {
	buf.m.Lock()
	defer buf.m.Unlock()
	return buf.b.Bytes()
}

// Flushes the contents of the buffer via `buffer.Reset()`.
func (buf *SafeBuffer) Flush() {
	buf.m.Lock()
	defer buf.m.Unlock()
	buf.b.Reset()
}

// Returns the size of the buffer.
func (buf *SafeBuffer) Size() int {
	buf.m.Lock()
	defer buf.m.Unlock()
	return buf.b.Len()
}

// Returns the contents of the buffer as a string.
func (buf *SafeBuffer) String() string {
	buf.m.Lock()
	defer buf.m.Unlock()
	return buf.b.String()
}

// Writes to the buffer, returning the number of bytes written.
func (buf *SafeBuffer) Write(p []byte) (n int, err error) {
	buf.m.Lock()
	defer buf.m.Unlock()
	return buf.b.Write(p)
}

// Creates a new slog logger with a backing thread-safe buffer.
func NewBufSlogTH(l *slog.Level) (*slog.Logger, *SafeBuffer) {
	var lev slog.Level
	if l == nil {
		lev = slog.LevelInfo
	} else {
		lev = *l
	}

	buffer := &SafeBuffer{}
	handler := slog.NewTextHandler(buffer, &slog.HandlerOptions{
		Level: lev,
	})
	logger := slog.New(handler)
	return logger, buffer
}

// Creates a new slog logger with a simple handler and a backing thread-safe buffer.
func NewBufSlogSH(l *slog.Level) (*slog.Logger, *SafeBuffer) {
	var lev slog.Level
	if l == nil {
		lev = slog.LevelInfo
	} else {
		lev = *l
	}

	buffer := &SafeBuffer{}
	handler := NewWorkerHandler(buffer, &slog.HandlerOptions{
		Level: lev,
	})
	logger := slog.New(handler)
	return logger, buffer
}
