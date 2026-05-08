// Package ptr implements functions to take the pointer of values.
package ptr

// Of returns a pointer to v.
//
// Deprecated: Use the built-in new operator instead of this function. This function is not needed
// and will be removed in a future release.
func Of[T any](v T) *T {
	return &v
}
