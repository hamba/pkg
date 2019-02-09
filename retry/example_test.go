package retry_test

import (
	"errors"
	"time"

	"github.com/hamba/pkg/retry"
)

func ExampleRun() {
	pol := retry.ExponentialPolicy(3, time.Millisecond)

	err := retry.Run(pol, func() error {
		// Do work

		return nil
	})
	if err != nil {
		// Handle the error
	}
}

func ExampleStop() {
	pol := retry.ExponentialPolicy(3, time.Millisecond)

	err := retry.Run(pol, func() error {
		// Do work that results in error

		return retry.Stop(errors.New("test error"))
	})
	if err != nil {
		// Handle the error
	}
}
