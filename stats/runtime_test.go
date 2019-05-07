package stats_test

import (
	"testing"
	"time"

	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/mock"
)

func TestRuntime(t *testing.T) {
	m := new(MockStats)
	m.On("Gauge", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("Timing", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	stats.DefaultRuntimeInterval = time.Microsecond

	go stats.Runtime(m)

	time.Sleep(100 * time.Millisecond)

	m.AssertCalled(t, "Gauge", "runtime.cpu.goroutines", mock.Anything, mock.Anything, mock.Anything)
}

func TestRuntimeFromStatable(t *testing.T) {
	m := new(MockStats)
	m.On("Gauge", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("Timing", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sable := stats.NewMockStatable(m)

	go stats.RuntimeFromStatable(sable, time.Millisecond)

	time.Sleep(100 * time.Millisecond)

	m.AssertCalled(t, "Gauge", "runtime.cpu.goroutines", mock.Anything, mock.Anything, mock.Anything)
}
