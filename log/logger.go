package log

import (
	"context"
	"fmt"
	"io"
	"os"
)

type key int

const (
	ctxKey key = iota
)

var (
	// Null is the null Logger instance.
	Null = &nullLogger{}
)

// Logger represents an abstract logging object.
type Logger interface {
	// Debug logs a debug message.
	Debug(msg string, ctx ...interface{})
	// Info logs an informational message.
	Info(msg string, ctx ...interface{})
	// Error logs an error message.
	Error(msg string, ctx ...interface{})
}

// WithLogger sets Logger in the context.
func WithLogger(ctx context.Context, logger Logger) context.Context {
	if logger == nil {
		logger = Null
	}
	return context.WithValue(ctx, ctxKey, logger)
}

// FromContext returns the instance Logger in the context.
func FromContext(ctx context.Context) (Logger, bool) {
	stats, ok := ctx.Value(ctxKey).(Logger)
	return stats, ok
}

// Debug logs a debug message.
func Debug(ctx context.Context, msg string, pairs ...interface{}) {
	withLogger(ctx, func(l Logger) {
		l.Debug(msg, pairs...)
	})
}

// Info logs an informational message.
func Info(ctx context.Context, msg string, pairs ...interface{}) {
	withLogger(ctx, func(l Logger) {
		l.Info(msg, pairs...)
	})
}

// Error logs an error message.
func Error(ctx context.Context, msg string, pairs ...interface{}) {
	withLogger(ctx, func(l Logger) {
		l.Error(msg, pairs...)
	})
}

type exitFunc func(int)

var exit exitFunc = os.Exit

// Fatal is equivalent to Error() followed by a call to os.Exit(1).
//
// Fatal will attempt to call Close() on the logger.
func Fatal(ctx context.Context, msg interface{}, pairs ...interface{}) {
	withLogger(ctx, func(l Logger) {
		l.Error(fmt.Sprintf("%+v", msg), pairs...)

		if cl, ok := l.(io.Closer); ok {
			_ = cl.Close()
		}
	})

	exit(1)
}

func withLogger(ctx context.Context, fn func(l Logger)) {
	if l, ok := FromContext(ctx); ok {
		fn(l)
		return
	}

	fn(Null)
}

type nullLogger struct{}

// Debug logs a debug message.
func (l nullLogger) Debug(msg string, ctx ...interface{}) {}

// Info logs an informational message.
func (l nullLogger) Info(msg string, ctx ...interface{}) {}

// Error logs an error message.
func (l nullLogger) Error(msg string, ctx ...interface{}) {}
