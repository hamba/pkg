package log_test

import (
	"context"
	"testing"

	"github.com/hamba/pkg/log"
	"github.com/stretchr/testify/mock"
)

func TestDebug(t *testing.T) {
	m := new(MockLogger)
	m.On("Debug", "test log", []interface{}{"foo", "bar"})
	ctx := log.WithLogger(context.Background(), m)

	log.Debug(ctx, "test log", "foo", "bar")

	m.AssertExpectations(t)
}

func TestInfo(t *testing.T) {
	m := new(MockLogger)
	m.On("Info", "test log", []interface{}{"foo", "bar"})
	ctx := log.WithLogger(context.Background(), m)

	log.Info(ctx, "test log", "foo", "bar")

	m.AssertExpectations(t)
}

func TestError(t *testing.T) {
	m := new(MockLogger)
	m.On("Error", "test log", []interface{}{"foo", "bar"})
	ctx := log.WithLogger(context.Background(), m)

	log.Error(ctx, "test log", "foo", "bar")

	m.AssertExpectations(t)
}

func TestNullLogger_Debug(t *testing.T) {
	log.Null.Debug("test log", "foo", "bar")
}

func TestNullLogger_Info(t *testing.T) {
	log.Null.Info("test log", "foo", "bar")
}

func TestNullLogger_Error(t *testing.T) {
	log.Null.Error("test log", "foo", "bar")
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
