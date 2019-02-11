package stats

import (
	"runtime"
	"time"
)

// DefaultRuntimeInterval is the default runtime ticker interval.
var DefaultRuntimeInterval = 30 * time.Second

// Runtime enters a loop, reporting runtime stats periodically.
func Runtime(stats Statter) {
	RuntimeEvery(stats, DefaultRuntimeInterval)
}

// RuntimeEvery enters a loop, reporting runtime stats at the specified interval.
func RuntimeEvery(s Statter, t time.Duration) {
	c := time.Tick(t)
	for range c {
		r := newRuntimeStats()
		r.send(s)
	}
}

// RuntimeFromStatable is the same as RuntimeEvery but from a Statable.
func RuntimeFromStatable(sable Statable, t time.Duration) {
	RuntimeEvery(sable.Statter(), t)
}

type runtimeStats struct {
	*runtime.MemStats

	goroutines int
}

func newRuntimeStats() *runtimeStats {
	r := &runtimeStats{MemStats: &runtime.MemStats{}}
	runtime.ReadMemStats(r.MemStats)
	r.goroutines = runtime.NumGoroutine()

	return r
}

func (r *runtimeStats) send(s Statter) {
	// CPU s
	s.Gauge("runtime.cpu.goroutines", float64(r.goroutines), 1.0)

	// Memory s
	// General
	s.Gauge("runtime.memory.alloc", float64(r.MemStats.Alloc), 1.0)
	s.Gauge("runtime.memory.total", float64(r.MemStats.TotalAlloc), 1.0)
	s.Gauge("runtime.memory.sys", float64(r.MemStats.Sys), 1.0)
	s.Gauge("runtime.memory.lookups", float64(r.MemStats.Lookups), 1.0)
	s.Gauge("runtime.memory.mallocs", float64(r.MemStats.Mallocs), 1.0)
	s.Gauge("runtime.memory.frees", float64(r.MemStats.Frees), 1.0)

	// Heap
	s.Gauge("runtime.memory.heap.alloc", float64(r.MemStats.HeapAlloc), 1.0)
	s.Gauge("runtime.memory.heap.sys", float64(r.MemStats.HeapSys), 1.0)
	s.Gauge("runtime.memory.heap.idle", float64(r.MemStats.HeapIdle), 1.0)
	s.Gauge("runtime.memory.heap.inuse", float64(r.MemStats.HeapInuse), 1.0)
	s.Gauge("runtime.memory.heap.objects", float64(r.MemStats.HeapObjects), 1.0)
	s.Gauge("runtime.memory.heap.released", float64(r.MemStats.HeapReleased), 1.0)

	// Stack
	s.Gauge("runtime.memory.stack.inuse", float64(r.MemStats.StackInuse), 1.0)
	s.Gauge("runtime.memory.stack.sys", float64(r.MemStats.StackSys), 1.0)
	s.Gauge("runtime.memory.stack.mcache_inuse", float64(r.MemStats.MCacheInuse), 1.0)
	s.Gauge("runtime.memory.stack.mcache_sys", float64(r.MemStats.MCacheSys), 1.0)
	s.Gauge("runtime.memory.stack.mspan_inuse", float64(r.MemStats.MSpanInuse), 1.0)
	s.Gauge("runtime.memory.stack.mspan_sys", float64(r.MemStats.MSpanSys), 1.0)

	// GC
	s.Gauge("runtime.memory.gc.last", float64(r.MemStats.LastGC), 1.0)
	s.Gauge("runtime.memory.gc.next", float64(r.MemStats.NextGC), 1.0)
	s.Gauge("runtime.memory.gc.count", float64(r.MemStats.NumGC), 1.0)
	s.Timing("runtime.memory.gc.pause", time.Duration(r.MemStats.PauseTotalNs), 1.0)
}
