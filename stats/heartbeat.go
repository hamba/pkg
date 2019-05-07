package stats

import (
	"time"
)

// DefaultHeartbeatInterval is the default heartbeat ticker interval.
var DefaultHeartbeatInterval = time.Second

// Heartbeat enters a loop, reporting a heartbeat counter periodically.
func Heartbeat(stats Statter) {
	HeartbeatEvery(stats, DefaultHeartbeatInterval)
}

// HeartbeatEvery enters a loop, reporting a heartbeat counter at the specified interval.
func HeartbeatEvery(stats Statter, t time.Duration) {
	c := time.Tick(t)
	for range c {
		stats.Inc("heartbeat", 1, 1.0)
	}
}

// HeartbeatFromStatable is the same as HeartbeatEvery but from context.
func HeartbeatFromStatable(sable Statable, t time.Duration) {
	HeartbeatEvery(sable.Statter(), t)
}
