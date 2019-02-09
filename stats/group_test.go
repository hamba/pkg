package stats_test

import (
	"context"
	"testing"
	"time"

	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/assert"
)

func TestGroup_Inc(t *testing.T) {
	m := new(MockStats)
	m.On("Inc", "prefix.test", int64(1), float32(1), []interface{}{"foo", "bar"}).Return(nil)
	ctx := stats.WithStatter(context.Background(), m)

	stats.Group(ctx, "prefix", func(s stats.Statter) {
		s.Inc("test", 1, 1.0, "foo", "bar")
	})

	m.AssertExpectations(t)
}

func TestGroup_Gauge(t *testing.T) {
	m := new(MockStats)
	m.On("Gauge", "prefix.test", float64(1), float32(1), []interface{}{"foo", "bar"}).Return(nil)
	ctx := stats.WithStatter(context.Background(), m)

	stats.Group(ctx, "prefix", func(s stats.Statter) {
		s.Gauge("test", 1.0, 1.0, "foo", "bar")
	})

	m.AssertExpectations(t)
}

func TestGroup_Timing(t *testing.T) {
	m := new(MockStats)
	m.On("Timing", "prefix.test", time.Millisecond, float32(1), []interface{}{"foo", "bar"}).Return(nil)
	ctx := stats.WithStatter(context.Background(), m)

	stats.Group(ctx, "prefix", func(s stats.Statter) {
		s.Timing("test", time.Millisecond, 1.0, "foo", "bar")
	})

	m.AssertExpectations(t)
}

func TestGroup_Close(t *testing.T) {
	m := new(MockStats)
	ctx := stats.WithStatter(context.Background(), m)

	assert.Panics(t, func() {
		stats.Group(ctx, "prefix", func(s stats.Statter) {
			_ = s.Close()
		})
	})
}
