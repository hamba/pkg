package stats_test

import (
	"testing"
	"time"

	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewAggregateStatter(t *testing.T) {
	m := new(MockStats)
	s := stats.NewAggregateStatter(m, time.Second)

	assert.Implements(t, (*stats.Statter)(nil), s)
	assert.IsType(t, &stats.AggregateStatter{}, s)
}

func TestNewAggregateStatter_WithCounterAggregator(t *testing.T) {
	agg := new(MockAggregator)
	agg.On("Aggregate", mock.Anything)
	agg.On("Flush", mock.Anything)
	statter := new(MockStats)
	statter.On("Close").Return(nil)
	s := stats.NewAggregateStatter(statter, time.Second, stats.WithCounterAggregator(agg))

	s.Inc("foobar", 1, 1)
	_ = s.Close()

	statter.AssertExpectations(t)
	agg.AssertExpectations(t)
}

func TestNewAggregateStatter_WithGaugeAggregator(t *testing.T) {
	agg := new(MockAggregator)
	agg.On("Aggregate", mock.Anything)
	agg.On("Flush", mock.Anything)
	statter := new(MockStats)
	statter.On("Close").Return(nil)
	s := stats.NewAggregateStatter(statter, time.Second, stats.WithGaugeAggregator(agg))

	s.Gauge("foobar", 1, 1)
	_ = s.Close()

	statter.AssertExpectations(t)
	agg.AssertExpectations(t)
}

func TestNewAggregateStatter_WithTimingAggregator(t *testing.T) {
	agg := new(MockAggregator)
	agg.On("Aggregate", mock.Anything)
	agg.On("Flush", mock.Anything)
	statter := new(MockStats)
	statter.On("Close").Return(nil)
	s := stats.NewAggregateStatter(statter, time.Second, stats.WithTimingAggregator(agg))

	s.Timing("foobar", 1, 1)
	_ = s.Close()

	statter.AssertExpectations(t)
	agg.AssertExpectations(t)
}

func TestAggregateStatter_Inc(t *testing.T) {
	m := new(MockStats)
	m.On("Inc", "test", int64(2), float32(1), []string{"foo", "bar"}).Once()
	m.On("Inc", "test1", int64(1), float32(1), []string{"foo", "bar"}).Once()
	m.On("Inc", "test", int64(1), float32(1), []string{"foo1", "bar"}).Once()
	m.On("Inc", "rate", int64(10), float32(1), []string(nil)).Once()
	m.On("Close").Return(nil)
	s := stats.NewAggregateStatter(m, time.Millisecond)

	s.Inc("test", 1, 1.0, "foo", "bar")
	s.Inc("test", 1, 1.0, "foo", "bar")
	s.Inc("test1", 1, 1.0, "foo", "bar")
	s.Inc("test", 1, 1.0, "foo1", "bar")
	s.Inc("rate", 1, 0.1)

	time.Sleep(10 * time.Millisecond)

	_ = s.Close()

	m.AssertExpectations(t)
}

func TestAggregateStatter_Gauge(t *testing.T) {
	m := new(MockStats)
	m.On("Gauge", "test", float64(3), float32(1), []string{"foo", "bar"}).Once()
	m.On("Gauge", "test1", float64(4), float32(1), []string{"foo", "bar"}).Once()
	m.On("Gauge", "test", float64(5), float32(1), []string{"foo1", "bar"}).Once()
	m.On("Close").Return(nil)
	s := stats.NewAggregateStatter(m, time.Millisecond)

	s.Gauge("test", 1, 1.0, "foo", "bar")
	s.Gauge("test", 3, 1.0, "foo", "bar")
	s.Gauge("test1", 4, 1.0, "foo", "bar")
	s.Gauge("test", 5, 1.0, "foo1", "bar")

	time.Sleep(10 * time.Millisecond)

	_ = s.Close()

	m.AssertExpectations(t)
}

func TestAggregateStatter_Timing(t *testing.T) {
	m := new(MockStats)
	m.On("Timing", "test", 500*time.Millisecond+500*time.Microsecond, float32(1), []string{"foo", "bar"}).Once()
	m.On("Timing", "test1", time.Millisecond, float32(1), []string{"foo", "bar"}).Once()
	m.On("Timing", "test", time.Millisecond, float32(1), []string{"foo1", "bar"}).Once()
	m.On("Timing", "rate", time.Millisecond, float32(1), []string(nil)).Once()
	m.On("Close").Return(nil)
	s := stats.NewAggregateStatter(m, time.Millisecond)

	s.Timing("test", time.Millisecond, 1.0, "foo", "bar")
	s.Timing("test", time.Second, 1.0, "foo", "bar")
	s.Timing("test1", time.Millisecond, 1.0, "foo", "bar")
	s.Timing("test", time.Millisecond, 1.0, "foo1", "bar")
	s.Timing("rate", time.Millisecond, 0.1)

	time.Sleep(10 * time.Millisecond)

	_ = s.Close()

	m.AssertExpectations(t)
}

func TestAggregateStatter_Close(t *testing.T) {
	m := new(MockStats)
	m.On("Close").Return(nil)
	s := stats.NewAggregateStatter(m, time.Second)

	err := s.Close()

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

type MockAggregator struct {
	mock.Mock
}

func (m *MockAggregator) Aggregate(metric stats.Metric) {
	m.Called(metric)
}

func (m *MockAggregator) Flush(statter stats.Statter) {
	m.Called(statter)
}
