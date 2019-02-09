package stats

import (
	"context"

	"github.com/hamba/pkg/timex"
)

// Timer represents a timer.
type Timer interface {
	// Start starts the timer.
	Start()
	// Done stops the timer and submits the Timing metric.
	Done()
}

type timer struct {
	start timex.Nanotime
	ctx   context.Context
	name  string
	rate  float32
	tags  []interface{}
}

// Time is a shorthand for Timing.
func Time(ctx context.Context, name string, rate float32, tags ...interface{}) Timer {
	t := &timer{ctx: ctx, name: name, rate: rate, tags: tags}
	t.Start()
	return t
}

func (t *timer) Start() {
	t.start = timex.Now()
}

func (t *timer) Done() {
	Timing(t.ctx, t.name, timex.Since(t.start), t.rate, t.tags...)
}
