// Package ptr implements functions to take the pointer of values.
package ptr

// Of returns a pointer to v.
func Of[T any](v T) *T {
	return &v
}
