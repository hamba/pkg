package retry_test

import (
	"errors"
	"testing"
	"time"

	"github.com/hamba/pkg/retry"
	"github.com/stretchr/testify/assert"
)

func TestExponentialPolicy(t *testing.T) {
	p := retry.ExponentialPolicy(3, time.Millisecond)

	sleep, ok := p.Next()
	assert.True(t, ok)
	assert.Equal(t, time.Millisecond, sleep)

	sleep, ok = p.Next()
	assert.True(t, ok)
	assert.Equal(t, 2*time.Millisecond, sleep)

	sleep, ok = p.Next()
	assert.False(t, ok)
	assert.Zero(t, sleep)

	sleep, ok = p.Next()
	assert.False(t, ok)
	assert.Zero(t, sleep)
}

func TestRun(t *testing.T) {
	var i int

	err := retry.Run(retry.ExponentialPolicy(3, time.Nanosecond), func() error {
		i++
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, i)
}

func TestRun_NilPolicy(t *testing.T) {
	err := retry.Run(nil, func() error {
		return nil
	})

	assert.Error(t, err)
}

func TestRun_MaxAttempts(t *testing.T) {
	var i int

	err := retry.Run(retry.ExponentialPolicy(3, time.Nanosecond), func() error {
		i++
		return errors.New("test error")
	})

	assert.Error(t, err)
	assert.Equal(t, 3, i)
}

func TestStop(t *testing.T) {
	var i int

	err := retry.Run(retry.ExponentialPolicy(3, time.Nanosecond), func() error {
		i++
		return retry.Stop(errors.New("test error"))
	})

	assert.Error(t, err)
	assert.Equal(t, 1, i)
}
