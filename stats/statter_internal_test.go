package stats

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithStatter(t *testing.T) {
	ctx := WithStatter(context.Background(), Null)

	got := ctx.Value(ctxKey)

	assert.Equal(t, Null, got)
}

func TestWithStatter_NilStats(t *testing.T) {
	ctx := WithStatter(context.Background(), nil)

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

func TestWithStatterFunc(t *testing.T) {
	tests := []struct {
		ctx    context.Context
		expect Statter
	}{
		{context.Background(), Null},
	}

	for _, tt := range tests {
		withStatter(tt.ctx, func(s Statter) {
			assert.Equal(t, tt.expect, s)
		})
	}
}

func TestNullStats(t *testing.T) {
	s := Null

	s.Inc("test", 1, 1.0)
	s.Gauge("test", 1.0, 1.0)
	s.Timing("test", 0, 1.0)

	assert.NoError(t, s.Close())
}

func BenchmarkTaggedStats_MergeTags(b *testing.B) {
	tags := []interface{}{
		"test1", "test",
		"test2", "test",
		"test3", "test",
		"test4", "test",
		"test5", "test",
	}
	addedTags := []interface{}{
		"k1", "v",
		"k2", "v",
		"k3", "v",
		"k4", "v",
		"k5", "v",
	}

	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		mergeTags(tags, addedTags)
	}
}
