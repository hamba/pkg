package stats_test

import (
	"testing"

	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/mock"
)

func TestTimer(t *testing.T) {
	m := new(MockStats)
	m.On("Timing", "test", mock.Anything, float32(1.0), mock.Anything).Return(nil)
	sable := &testStatable{s: m}

	ti := stats.Time(sable, "test", 1.0)
	ti.Done()

	m.AssertExpectations(t)
}
