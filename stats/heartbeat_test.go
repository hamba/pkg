package stats_test

import (
	"testing"
	"time"

	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/mock"
)

func TestHeartbeat(t *testing.T) {
	m := new(MockStats)
	m.On("Inc", "heartbeat", int64(1), float32(1.0), mock.Anything).Return(nil)
	stats.DefaultHeartbeatInterval = time.Millisecond

	go stats.Heartbeat(m)

	time.Sleep(100 * time.Millisecond)

	m.AssertCalled(t, "Inc", "heartbeat", int64(1), float32(1.0), mock.Anything)
}

func TestHeartbeatFromStatable(t *testing.T) {
	m := new(MockStats)
	m.On("Inc", "heartbeat", int64(1), float32(1.0), mock.Anything).Return(nil)
	sable := stats.NewMockStatable(m)

	go stats.HeartbeatFromStatable(sable, time.Millisecond)

	time.Sleep(100 * time.Millisecond)

	m.AssertCalled(t, "Inc", "heartbeat", int64(1), float32(1.0), mock.Anything)
}
