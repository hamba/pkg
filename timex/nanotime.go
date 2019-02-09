// Package timex implements extensions to the time package.
package timex

import (
	_ "unsafe" // Required in order to import nanotime
)

// Nanotime gets the current time in nanoseconds.
//go:linkname Nanotime runtime.nanotime
func Nanotime() int64
