// Package errors implements common error types.
package errors

// Error is an error that is able to be a constant.
//
// See https://dave.cheney.net/2016/04/07/constant-errors for details.
type Error string

// Error returns the error string.
func (e Error) Error() string {
	return string(e)
}
