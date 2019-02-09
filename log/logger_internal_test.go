package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFatal(t *testing.T) {
	// Switch out the exit func
	calledCode := -1
	exit = func(code int) {
		calledCode = code
	}

	m := new(MockLogger)
	m.On("Error", "test log", []interface{}{"foo", "bar"})
	m.On("Close").Return(nil)
	labl := &testLoggable{l: m}

	Fatal(labl, "test log", "foo", "bar")

	m.AssertExpectations(t)
	assert.Equal(t, 1, calledCode)
}

type testLoggable struct {
	l Logger
}

func (l *testLoggable) Logger() Logger {
	return l.l
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, ctx ...interface{}) {
	m.Called(msg, ctx)
}

func (m *MockLogger) Info(msg string, ctx ...interface{}) {
	m.Called(msg, ctx)
}

func (m *MockLogger) Error(msg string, ctx ...interface{}) {
	m.Called(msg, ctx)
}

func (m *MockLogger) Close() error {
	args := m.Called()

	return args.Error(0)
}
