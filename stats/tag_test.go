package stats_test

import (
	"testing"
	"time"

	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/assert"
)

func TestNewTaggedStatter(t *testing.T) {
	m := new(MockStats)

	s := stats.NewTaggedStatter(m, "global", "foobar")

	assert.Implements(t, (*stats.Statter)(nil), s)
	assert.Implements(t, (*stats.Wrapper)(nil), s)
	assert.IsType(t, &stats.TaggedStatter{}, s)
}

func TestTaggedStatter_MergesWithPreviousTaggedStater(t *testing.T) {
	m := new(MockStats)
	m.On("Inc", "test", int64(1), float32(1), []string{"foo", "bar", "test", "1234", "global", "foobar"}).Return(nil)

	s1 := stats.NewTaggedStatter(m, "test", "1234")

	s2 := stats.NewTaggedStatter(s1, "global", "foobar")

	s2.Inc("test", 1, 1.0, "foo", "bar")

	m.AssertExpectations(t)
}

func TestTaggedStatter_NormalisesTags(t *testing.T) {
	m := new(MockStats)
	m.On("Inc", "test", int64(1), float32(1), []string{"foo", "bar", "global", "", "STATTER_ERROR", "Normalised odd number of tags by adding an empty string"}).Return(nil)
	s := stats.NewTaggedStatter(m, "global")

	s.Inc("test", 1, 1.0, "foo", "bar")

	m.AssertExpectations(t)
}

func TestTaggedStatter_Inc(t *testing.T) {
	m := new(MockStats)
	m.On("Inc", "test", int64(1), float32(1), []string{"foo", "bar", "global", "foobar"}).Return(nil)
	s := stats.NewTaggedStatter(m, "global", "foobar")

	s.Inc("test", 1, 1.0, "foo", "bar")

	m.AssertExpectations(t)
}

func TestTaggedStatter_Gauge(t *testing.T) {
	m := new(MockStats)
	m.On("Gauge", "test", float64(1), float32(1), []string{"foo", "bar", "global", "foobar"}).Return(nil)
	s := stats.NewTaggedStatter(m, "global", "foobar")

	s.Gauge("test", 1.0, 1.0, "foo", "bar")

	m.AssertExpectations(t)
}

func TestTaggedStatter_Timing(t *testing.T) {
	m := new(MockStats)
	m.On("Timing", "test", time.Millisecond, float32(1), []string{"foo", "bar", "global", "foobar"}).Return(nil)
	s := stats.NewTaggedStatter(m, "global", "foobar")

	s.Timing("test", time.Millisecond, 1.0, "foo", "bar")

	m.AssertExpectations(t)
}

func TestTaggedStatter_Unwrap(t *testing.T) {
	m := new(MockStats)
	s := stats.NewTaggedStatter(m, "global", "foobar")

	got := s.Unwrap()

	assert.Equal(t, m, got)
}

func TestTaggedStatter_Close(t *testing.T) {
	m := new(MockStats)
	m.On("Close").Return(nil)
	s := stats.NewTaggedStatter(m, "global", "foobar")

	_ = s.Close()

	m.AssertExpectations(t)
}
