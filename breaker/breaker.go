// Package breaker implements the Circuit Breaker design patter.
package breaker

import (
	"errors"
	"sync"
	"time"
)

// State represents a Breakers state.
type State int8

// State constants for the Breaker.
const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

var (
	// ErrOpenState is the error returned when the Breaker is open.
	ErrOpenState = errors.New("breaker: circuit breaker is open")

	// ErrTooManyRequests is the error returned when too many test requests are
	// made when the Breaker is half-open.
	ErrTooManyRequests = errors.New("breaker: too many requests")
)

// Counter holds the number number of requests, successes and failures of a breaker.
//
// Counter is reset from time to time, and totals should not be used as a full totals.
type Counter struct {
	// Requests is the total number of requests made.
	Requests uint64
	// Successes is the total number of successes returned.
	Successes uint64
	// Failures is the total number of failures returned.
	Failures uint64
	// ConsecutiveSuccesses is the number of consecutive successes returned.
	ConsecutiveSuccesses uint64
	// ConsecutiveFailures is the number of consecutive failures returned.
	ConsecutiveFailures uint64
}

func (c *Counter) reset() {
	c.Requests = 0
	c.Successes = 0
	c.Failures = 0
	c.ConsecutiveSuccesses = 0
	c.ConsecutiveFailures = 0
}

func (c *Counter) request() {
	c.Requests++
}

func (c *Counter) success() {
	c.Successes++
	c.ConsecutiveSuccesses++
	c.ConsecutiveFailures = 0
}

func (c *Counter) failure() {
	c.Failures++
	c.ConsecutiveFailures++
	c.ConsecutiveSuccesses = 0
}

// Fuse represents a Breaker fuse used to trip the breaker.
type Fuse interface {
	// Trip determines if the Breaker should be tripped.
	Trip(Counter) bool
}

// FuseFunc is an adapter allowing to use a function as a Fuse.
type FuseFunc func(Counter) bool

// Trip determines if the Breaker should be tripped.
func (f FuseFunc) Trip(c Counter) bool {
	return f(c)
}

// ThresholdFuse trips the Breaker when the total number of failures exceeds the given count.
func ThresholdFuse(count uint64) FuseFunc {
	return FuseFunc(func(c Counter) bool {
		return c.Failures > count
	})
}

// ConsecutiveFuse trips the Breaker when the consecutive number of failures exceeds the given count.
func ConsecutiveFuse(count uint64) FuseFunc {
	return FuseFunc(func(c Counter) bool {
		return c.ConsecutiveFailures > count
	})
}

// RateFuse trips the Breaker when the percentage of failures exceeds the given rate.
//
// rate should be between 0 and 100.
func RateFuse(rate uint64) FuseFunc {
	if rate >= 100 {
		panic(errors.New("breaker: rate should be less than 100 percent"))
	}

	return FuseFunc(func(c Counter) bool {
		return c.Failures/c.Requests*100 > rate
	})
}

// OptFunc represents a configuration function for Breaker.
type OptFunc func(b *Breaker)

// WithSleep sets the time the Breaker stays open for.
func WithSleep(d time.Duration) OptFunc {
	return OptFunc(func(b *Breaker) {
		b.sleep = d
	})
}

// WithTestRequests sets the number of test requests allowed when the Breaker is half-open.
func WithTestRequests(c uint64) OptFunc {
	return OptFunc(func(b *Breaker) {
		b.testRequests = c
	})
}

// Breaker is a circuit breaker.
type Breaker struct {
	fuse         Fuse
	sleep        time.Duration
	testRequests uint64

	mu        sync.Mutex
	state     State
	counter   Counter
	openUntil time.Time
}

// NewBreaker creates a new Breaker.
func NewBreaker(f Fuse, opts ...OptFunc) *Breaker {
	b := &Breaker{
		fuse:         f,
		sleep:        10 * time.Second,
		testRequests: 1,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

// Run runs the given request if the Breaker allows it.
//
// Run returns an error immediately if the Breaker rejects the request,
// otherwise it returns the result of the request.
func (b *Breaker) Run(req func() error) error {
	if err := b.canRun(); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			b.handleResult(false)
			panic(r)
		}
	}()

	err := req()

	b.handleResult(err == nil)

	return err
}

func (b *Breaker) canRun() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	state := b.getState()

	if state == StateOpen {
		return ErrOpenState
	} else if state == StateHalfOpen && b.counter.Requests >= b.testRequests {
		return ErrTooManyRequests
	}

	b.counter.request()

	return nil
}

func (b *Breaker) handleResult(success bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	state := b.getState()

	if !success {
		b.handleFailure(state)
		return
	}

	b.handleSuccess(state)
}

func (b *Breaker) handleSuccess(state State) {
	switch state {
	case StateClosed:
		b.counter.success()

	case StateHalfOpen:
		b.counter.success()
		if state == StateHalfOpen && b.counter.ConsecutiveSuccesses >= b.testRequests {
			b.setState(StateClosed)
		}
	}
}

func (b *Breaker) handleFailure(state State) {
	switch state {
	case StateClosed:
		b.counter.failure()
		if b.fuse.Trip(b.counter) {
			b.setState(StateOpen)
		}

	case StateHalfOpen:
		b.setState(StateOpen)
	}
}

// State returns the state of the Breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.getState()
}

func (b *Breaker) getState() State {
	now := time.Now()

	switch b.state {
	case StateClosed:
		if !b.openUntil.IsZero() && b.openUntil.Before(now) {
			b.openUntil = time.Time{}
		}

	case StateOpen:
		if b.openUntil.Before(now) {
			b.setState(StateHalfOpen)
		}
	}
	return b.state
}

func (b *Breaker) setState(state State) {
	b.state = state

	if b.state == StateOpen {
		b.openUntil = time.Now().Add(b.sleep)
	}

	b.counter.reset()
}
