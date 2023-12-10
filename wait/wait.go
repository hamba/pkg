// Package wait provides functions for polling for condition changes.
package wait

import (
	"context"
	"time"
)

// ConditionFunc returns true if the condition is satisfied, or an error
// if the loop should be aborted.
type ConditionFunc func(context.Context) (done bool, err error)

// PollUntil tries a condition until stopped by the context.
func PollUntil(ctx context.Context, fn ConditionFunc, interval time.Duration) error {
	return poll(ctx, false, fn, interval)
}

// PollImmediateUntil tries a condition until stopped by the context.
func PollImmediateUntil(ctx context.Context, fn ConditionFunc, interval time.Duration) error {
	return poll(ctx, true, fn, interval)
}

func poll(ctx context.Context, immediate bool, fn ConditionFunc, interval time.Duration) error {
	if immediate {
		done, err := fn(ctx)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}

	tick := time.NewTicker(interval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
			done, err := fn(ctx)
			if err != nil {
				return err
			}
			if done {
				return nil
			}
		}
	}
}
