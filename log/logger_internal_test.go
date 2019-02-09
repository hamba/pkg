package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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
