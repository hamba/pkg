package timex_test

import (
	"testing"
	"time"

	"github.com/hamba/pkg/timex"
	"github.com/stretchr/testify/assert"
)

func TestNanotime_Since(t *testing.T) {
	then := timex.Nanotime(timex.Now() - 1000)

	d := timex.Since(then)

	assert.True(t, time.Duration(1000) <= d)
}
