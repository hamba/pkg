package breaker_test

import (
	"time"

	"github.com/hamba/pkg/breaker"
)

func ExampleBreaker() {
	b := breaker.NewBreaker(breaker.ThresholdFuse(1), breaker.WithSleep(100*time.Millisecond))

	err := b.Run(func() error {
		// Your code protected by the circuit breaker...

		return nil // Return any errors
	})
	if err != nil {
		// Handle the error
	}

}
