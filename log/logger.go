package log

import (
	"fmt"
	"io"
	"os"
)

var (
	// Null is the null Logger instance.
	Null = &nullLogger{}
)

// Loggable represents an object that has a Logger.
type Loggable interface {
	Logger() Logger
}

// Logger represents an abstract logging object.
type Logger interface {
	// Debug logs a debug message.
	Debug(msg string, ctx ...interface{})
	// Info logs an informational message.
	Info(msg string, ctx ...interface{})
	// Error logs an error message.
	Error(msg string, ctx ...interface{})
}

// Debug logs a debug message.
func Debug(lable Loggable, msg string, pairs ...interface{}) {
	lable.Logger().Debug(msg, pairs...)
}

// Info logs an informational message.
func Info(lable Loggable, msg string, pairs ...interface{}) {
	lable.Logger().Info(msg, pairs...)
}

// Error logs an error message.
func Error(lable Loggable, msg string, pairs ...interface{}) {
	lable.Logger().Error(msg, pairs...)
}

type exitFunc func(int)

var exit exitFunc = os.Exit

// Fatal is equivalent to Error() followed by a call to os.Exit(1).
//
// Fatal will attempt to call Close() on the logger.
func Fatal(lable Loggable, msg interface{}, pairs ...interface{}) {
	l := lable.Logger()
	l.Error(fmt.Sprintf("%+v", msg), pairs...)
	if cl, ok := l.(io.Closer); ok {
		_ = cl.Close()
	}

	exit(1)
}

type nullLogger struct{}

// Debug logs a debug message.
func (l nullLogger) Debug(msg string, ctx ...interface{}) {}

// Info logs an informational message.
func (l nullLogger) Info(msg string, ctx ...interface{}) {}

// Error logs an error message.
func (l nullLogger) Error(msg string, ctx ...interface{}) {}
