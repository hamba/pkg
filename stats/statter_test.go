package stats_test

import (
	"context"
	"testing"
	"time"

	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/assert"
)

func TestInc(t *testing.T) {
	m := new(MockStats)
	m.On("Inc", "test", int64(1), float32(1.0), []interface{}(nil))
	ctx := stats.WithStatter(context.Background(), m)

	stats.Inc(ctx, "test", 1, 1.0)

	m.AssertExpectations(t)
}

func TestGauge(t *testing.T) {
	m := new(MockStats)
	m.On("Gauge", "test", float64(1), float32(1.0), []interface{}(nil))
	ctx := stats.WithStatter(context.Background(), m)

	stats.Gauge(ctx, "test", 1, 1.0)

	m.AssertExpectations(t)
}

func TestTiming(t *testing.T) {
	m := new(MockStats)
	m.On("Timing", "test", time.Second, float32(1.0), []interface{}(nil))
	ctx := stats.WithStatter(context.Background(), m)

	stats.Timing(ctx, "test", time.Second, 1.0)

	m.AssertExpectations(t)
}

func TestClose(t *testing.T) {
	m := new(MockStats)
	m.On("Close").Return(nil)
	ctx := stats.WithStatter(context.Background(), m)

	err := stats.Close(ctx)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestClose_NoContext(t *testing.T) {
	err := stats.Close(context.Background())

	assert.NoError(t, err)
}
