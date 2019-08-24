package stats_test

import (
	"testing"
	"time"

	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/assert"
)

func TestUnwrap(t *testing.T) {
	m := new(MockStats)
	s := &testWrappedStats{stats: m}

	got := stats.Unwrap(s)

	assert.Equal(t, m, got)
}

func TestUnwrap_NotWrappedReturnsNil(t *testing.T) {
	m := new(MockStats)

	got := stats.Unwrap(m)

	assert.Nil(t, got)
}

type testWrappedStats struct {
	stats stats.Statter
}

func (t testWrappedStats) Inc(name string, value int64, rate float32, tags ...string) {}

func (t testWrappedStats) Gauge(name string, value float64, rate float32, tags ...string) {}

func (t testWrappedStats) Timing(name string, value time.Duration, rate float32, tags ...string) {}

func (t testWrappedStats) Unwrap() stats.Statter {
	return t.stats
}

func (t testWrappedStats) Close() error {
	return nil
}
