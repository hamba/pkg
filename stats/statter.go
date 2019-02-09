package stats

import (
	"context"
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

// Statter represents a stats instance.
type Statter interface {
	io.Closer

	// Inc increments a count by the value.
	Inc(name string, value int64, rate float32, tags ...interface{})

	// Gauge measures the value of a metric.
	Gauge(name string, value float64, rate float32, tags ...interface{})

	// Timing sends the value of a Duration.
	Timing(name string, value time.Duration, rate float32, tags ...interface{})
}

// WithStatter sets Statter in the context.
func WithStatter(ctx context.Context, s Statter) context.Context {
	if s == nil {
		s = Null
	}

	return context.WithValue(ctx, ctxKey, s)
}

// FromContext returns the instance of Statter in the context.
func FromContext(ctx context.Context) (Statter, bool) {
	stats, ok := ctx.Value(ctxKey).(Statter)
	return stats, ok
}

// Inc increments a count by the value.
func Inc(ctx context.Context, name string, value int64, rate float32, tags ...interface{}) {
	withStatter(ctx, func(s Statter) {
		s.Inc(name, value, rate, tags...)
	})
}

// Gauge measures the value of a metric.
func Gauge(ctx context.Context, name string, value float64, rate float32, tags ...interface{}) {
	withStatter(ctx, func(s Statter) {
		s.Gauge(name, value, rate, tags...)
	})
}

// Timing sends the value of a Duration.
func Timing(ctx context.Context, name string, value time.Duration, rate float32, tags ...interface{}) {
	withStatter(ctx, func(s Statter) {
		s.Timing(name, value, rate, tags...)
	})
}

// Close closes the client and flushes buffered stats, if applicable
func Close(ctx context.Context) error {
	if s, ok := FromContext(ctx); ok {
		return s.Close()
	}

	return nil
}

func withStatter(ctx context.Context, fn func(s Statter)) {
	if s, ok := FromContext(ctx); ok {
		fn(s)
		return
	}

	fn(Null)
}

type nullStatter struct{}

func (s nullStatter) Inc(name string, value int64, rate float32, tags ...interface{}) {}

func (s nullStatter) Gauge(name string, value float64, rate float32, tags ...interface{}) {}

func (s nullStatter) Timing(name string, value time.Duration, rate float32, tags ...interface{}) {}

func (s nullStatter) Close() error {
	return nil
}
