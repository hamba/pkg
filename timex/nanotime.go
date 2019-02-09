// Package timex implements extensions to the time package.
package timex

import (
	"time"
	_ "unsafe" // Required in order to import nanotime
)

//go:linkname getNanotime runtime.nanotime
func getNanotime() int64

// Nanotime represents an instant in time with nanosecond precision useing runtime.nanotime.
type Nanotime int64

// Now gets the current local time as a Nanotime.
func Now() Nanotime {
	return Nanotime(getNanotime())
}

// Since returns the time elapsed since t.
func Since(t Nanotime) time.Duration {
	return time.Duration(Now() - t)
}
