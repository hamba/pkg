package breaker_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/hamba/pkg/breaker"
	"github.com/stretchr/testify/assert"
)

func TestThresholdFuse(t *testing.T) {
	f := breaker.ThresholdFuse(5)

	assert.False(t, f.Trip(breaker.Counter{Failures: 5}))
	assert.True(t, f.Trip(breaker.Counter{Failures: 6}))
}

func TestConsecutiveFuse(t *testing.T) {
	f := breaker.ConsecutiveFuse(5)

	assert.False(t, f.Trip(breaker.Counter{ConsecutiveFailures: 5}))
	assert.True(t, f.Trip(breaker.Counter{ConsecutiveFailures: 6}))
}

func TestRateFuse(t *testing.T) {
	f := breaker.RateFuse(10)

	assert.False(t, f.Trip(breaker.Counter{Requests: 10, Failures: 9}))
	assert.True(t, f.Trip(breaker.Counter{Requests: 10, Failures: 10}))
}

func TestRateFusePanicsIfRateOver100(t *testing.T) {
	assert.Panics(t, func() { breaker.RateFuse(101) })
}

func TestNewBreaker(t *testing.T) {
	b := breaker.NewBreaker(breaker.ThresholdFuse(5))

	assert.IsType(t, &breaker.Breaker{}, b)
	assert.Equal(t, breaker.StateClosed, b.State())
}

func TestBreaker_Run(t *testing.T) {
	b := breaker.NewBreaker(breaker.ThresholdFuse(1), breaker.WithSleep(100*time.Millisecond))

	err := b.Run(successFunc)

	assert.NoError(t, err)
	assert.Equal(t, breaker.StateClosed, b.State())

	err = b.Run(failureFunc)

	assert.Error(t, err)
	assert.Equal(t, breaker.StateClosed, b.State())

	err = b.Run(failureFunc)

	assert.Error(t, err)
	assert.Equal(t, breaker.StateOpen, b.State())

	err = b.Run(failureFunc)

	assert.Equal(t, breaker.ErrOpenState, err)
	assert.Equal(t, breaker.StateOpen, b.State())

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, breaker.StateHalfOpen, b.State())

	err = b.Run(failureFunc)

	assert.Equal(t, errTest, err)
	assert.Equal(t, breaker.StateOpen, b.State())

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, breaker.StateHalfOpen, b.State())

	err = b.Run(successFunc)

	assert.NoError(t, err)
	assert.Equal(t, breaker.StateClosed, b.State())
}

func TestBreaker_RunOnlyAllowsXTestRequests(t *testing.T) {
	b := breaker.NewBreaker(breaker.ThresholdFuse(1), breaker.WithSleep(100*time.Millisecond), breaker.WithTestRequests(2))
	_ = b.Run(failureFunc)
	_ = b.Run(failureFunc)

	assert.Equal(t, breaker.StateOpen, b.State())

	time.Sleep(101 * time.Millisecond)

	errs := []error{}
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			err := b.Run(delayedSuccessFunc)

			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()

			wg.Done()
		}()
	}

	wg.Wait()

	tooMany := 0
	for _, err := range errs {
		if err == breaker.ErrTooManyRequests {
			tooMany++
		}
	}

	assert.Equal(t, 2, tooMany)
}

func TestBreaker_RunHandlesPanic(t *testing.T) {
	b := breaker.NewBreaker(breaker.ThresholdFuse(1), breaker.WithSleep(100*time.Millisecond))
	_ = b.Run(failureFunc)

	assert.Panics(t, func() { _ = b.Run(panicFunc) })
	assert.Equal(t, breaker.StateOpen, b.State())
}

func BenchmarkBreaker_RunSuccess(b *testing.B) {
	br := breaker.NewBreaker(breaker.ThresholdFuse(1), breaker.WithSleep(100*time.Millisecond))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = br.Run(successFunc)
	}
}

func BenchmarkBreaker_RunFailure(b *testing.B) {
	br := breaker.NewBreaker(breaker.ThresholdFuse(uint64(b.N)), breaker.WithSleep(100*time.Millisecond))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = br.Run(failureFunc)
	}
}

var errTest = errors.New("test")

func successFunc() error {
	return nil
}

func delayedSuccessFunc() error {
	time.Sleep(10 * time.Millisecond)
	return nil
}

func failureFunc() error {
	return errTest
}

func panicFunc() error {
	panic(errTest)
}
