package stats_test

import (
	"testing"
	"time"

	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/assert"
)

func TestInc(t *testing.T) {
	m := new(MockStats)
	m.On("Inc", "test", int64(1), float32(1.0), []string(nil))
	sable := &testStatable{s: m}

	stats.Inc(sable, "test", 1, 1.0)

	m.AssertExpectations(t)
}

func TestGauge(t *testing.T) {
	m := new(MockStats)
	m.On("Gauge", "test", float64(1), float32(1.0), []string(nil))
	sable := &testStatable{s: m}

	stats.Gauge(sable, "test", 1, 1.0)

	m.AssertExpectations(t)
}

func TestTiming(t *testing.T) {
	m := new(MockStats)
	m.On("Timing", "test", time.Second, float32(1.0), []string(nil))
	sable := &testStatable{s: m}

	stats.Timing(sable, "test", time.Second, 1.0)

	m.AssertExpectations(t)
}

func TestClose(t *testing.T) {
	m := new(MockStats)
	m.On("Close").Return(nil)
	sable := &testStatable{s: m}

	err := stats.Close(sable)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestNullStats(t *testing.T) {
	s := stats.Null

	s.Inc("test", 1, 1.0)
	s.Gauge("test", 1.0, 1.0)
	s.Timing("test", 0, 1.0)

	assert.NoError(t, s.Close())
}
