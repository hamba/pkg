// Package ptr implements functions to take the pointer of values.
package ptr

// Bool converts a bool into a bool pointer.
func Bool(b bool) *bool {
	return &b
}

// Float32 converts a float32 into a float32 pointer.
func Float32(f float32) *float32 {
	return &f
}

// Float64 converts a float64 into a float64 pointer.
func Float64(f float64) *float64 {
	return &f
}

// Int converts an int into an int pointer.
func Int(i int) *int {
	return &i
}

// Int64 converts an int64 into an int64 pointer.
func Int64(i int64) *int64 {
	return &i
}

// String converts a string into a string pointer.
func String(s string) *string {
	return &s
}
