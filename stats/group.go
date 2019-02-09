package stats

import (
	"time"
)

type group struct {
	s      Statter
	prefix string
}

// Inc increments a count by the value.
func (g group) Inc(name string, value int64, rate float32, tags ...interface{}) {
	g.s.Inc(g.prefix+name, value, rate, tags...)
}

// Gauge measures the value of a metric.
func (g group) Gauge(name string, value float64, rate float32, tags ...interface{}) {
	g.s.Gauge(g.prefix+name, value, rate, tags...)
}

// Timing sends the value of a Duration.
func (g group) Timing(name string, value time.Duration, rate float32, tags ...interface{}) {
	g.s.Timing(g.prefix+name, value, rate, tags...)
}

// Close panics, as groups cannot be closed.
func (g group) Close() error {
	panic("stats: cannot close a group")
}

// Group adds a common prefix to a set of stats.
func Group(sable Statable, prefix string, fn func(s Statter)) {
	if prefix != "" {
		prefix += "."
	}

	grp := group{
		s:      sable.Statter(),
		prefix: prefix,
	}

	fn(grp)
}
