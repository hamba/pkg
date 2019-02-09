package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWithLogger(t *testing.T) {
	ctx := WithLogger(context.Background(), Null)

	got := ctx.Value(ctxKey)

	assert.Equal(t, Null, got)
}

func TestWithLogger_NilLogger(t *testing.T) {
	ctx := WithLogger(context.Background(), nil)

	got := ctx.Value(ctxKey)

	assert.Equal(t, Null, got)
}

func TestFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKey, Null)

	got, ok := FromContext(ctx)

	assert.True(t, ok)
	assert.Equal(t, Null, got)
}

func TestFromContext_NotSet(t *testing.T) {
	ctx := context.Background()

	got, ok := FromContext(ctx)

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestFatal(t *testing.T) {
	// Switch out the exit func
	calledCode := -1
	exit = func(code int) {
		calledCode = code
	}

	m := new(MockLogger)
	m.On("Error", "test log", []interface{}{"foo", "bar"})
	m.On("Close", ).Return(nil)
	ctx := WithLogger(context.Background(), m)

	Fatal(ctx, "test log", "foo", "bar")

	m.AssertExpectations(t)
	assert.Equal(t, 1, calledCode)
}

func TestWithLoggerFunc(t *testing.T) {
	tests := []struct {
		ctx    context.Context
		expect Logger
	}{
		{context.Background(), Null},
	}

	for _, tt := range tests {
		withLogger(tt.ctx, func(l Logger) {
			assert.Equal(t, tt.expect, l)
		})
	}
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
