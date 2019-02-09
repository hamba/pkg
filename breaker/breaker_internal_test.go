package breaker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithSleep(t *testing.T) {
	b := &Breaker{}

	f := WithSleep(time.Millisecond)
	f(b)

	assert.Equal(t, time.Millisecond, b.sleep)
}

func TestWithTestRequests(t *testing.T) {
	b := &Breaker{}

	f := WithTestRequests(10)
	f(b)

	assert.Equal(t, uint64(10), b.testRequests)
}
