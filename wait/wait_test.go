package wait_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hamba/pkg/v2/wait"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPollUntil(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	called := make(chan struct{})
	go func() {
		defer close(called)

		err := wait.PollUntil(ctx, func(context.Context) (done bool, err error) {
			called <- struct{}{}
			return false, nil
		}, time.Microsecond)

		assert.ErrorIs(t, err, context.Canceled)
	}()

	// Wait for the initial condition call, and the first tick
	// condition call.
	<-called
	<-called

	// Stop waiting.
	cancel()

	// Assert that the condition is not called more than once after
	// canceling the context.
	var calledCount int
	for range called {
		calledCount++
	}
	assert.LessOrEqual(t, calledCount, 1)
}

func TestPollUntil_HandlesError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var count int
	err := wait.PollUntil(ctx, func(context.Context) (done bool, err error) {
		count++
		if count < 2 {
			return false, nil
		}
		return false, errors.New("test")
	}, time.Microsecond)

	require.Error(t, err)
}

func TestPollUntil_HandlesDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var count int
	err := wait.PollUntil(ctx, func(context.Context) (done bool, err error) {
		count++
		if count < 2 {
			return false, nil
		}
		return true, nil
	}, time.Microsecond)

	require.NoError(t, err)
}

func TestPollUntil_HandlesImmediateError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := wait.PollUntil(ctx, func(context.Context) (done bool, err error) {
		return false, errors.New("test")
	}, time.Microsecond)

	require.Error(t, err)
}

func TestPollUntil_HandlesImmediateDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := wait.PollUntil(ctx, func(context.Context) (done bool, err error) {
		return true, nil
	}, time.Microsecond)

	require.NoError(t, err)
}
