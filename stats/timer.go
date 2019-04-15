package stats

import (
	"github.com/hamba/timex/mono"
)

// Timer represents a timer.
type Timer interface {
	// Start starts the timer.
	Start()
	// Done stops the timer and submits the Timing metric.
	Done()
}

type timer struct {
	sable Statable
	start int64
	name  string
	rate  float32
	tags  []string
}

// Time is a shorthand for Timing.
func Time(sable Statable, name string, rate float32, tags ...string) Timer {
	t := &timer{sable: sable, name: name, rate: rate, tags: tags}
	t.Start()
	return t
}

func (t *timer) Start() {
	t.start = mono.Now()
}

func (t *timer) Done() {
	Timing(t.sable, t.name, mono.Since(t.start), t.rate, t.tags...)
}
