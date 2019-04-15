package stats

import (
	"io"
	"time"
)

type key int

const (
	ctxKey key = iota
)

var (
	// Null is the null Statter instance.
	Null = &nullStatter{}
)

// Statable represents an object that has a Statter.
type Statable interface {
	Statter() Statter
}

// Statter represents a stats instance.
type Statter interface {
	io.Closer

	// Inc increments a count by the value.
	Inc(name string, value int64, rate float32, tags ...string)

	// Gauge measures the value of a metric.
	Gauge(name string, value float64, rate float32, tags ...string)

	// Timing sends the value of a Duration.
	Timing(name string, value time.Duration, rate float32, tags ...string)
}

// Inc increments a count by the value.
func Inc(sable Statable, name string, value int64, rate float32, tags ...string) {
	sable.Statter().Inc(name, value, rate, tags...)
}

// Gauge measures the value of a metric.
func Gauge(sable Statable, name string, value float64, rate float32, tags ...string) {
	sable.Statter().Gauge(name, value, rate, tags...)
}

// Timing sends the value of a Duration.
func Timing(sable Statable, name string, value time.Duration, rate float32, tags ...string) {
	sable.Statter().Timing(name, value, rate, tags...)
}

// Close closes the client and flushes buffered stats, if applicable
func Close(sable Statable) error {
	return sable.Statter().Close()
}

type nullStatter struct{}

func (s nullStatter) Inc(name string, value int64, rate float32, tags ...string) {}

func (s nullStatter) Gauge(name string, value float64, rate float32, tags ...string) {}

func (s nullStatter) Timing(name string, value time.Duration, rate float32, tags ...string) {}

func (s nullStatter) Close() error {
	return nil
}
