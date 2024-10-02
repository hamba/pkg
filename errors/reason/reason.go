// Package reason contains types and functions commonly used to handle reason detection and extraction.
package reason

import (
	"errors"
	"fmt"
)

// InternalReason is the reason set when the system is experiencing
// an error that the user cannot resolve.
const InternalReason = "The system has an internal error"

// Error is a reason error that is detectable.
//
// This should not be used as a "normal" error,
// instead extract the reasons from the chains
// using Extract.
type Error struct {
	// Msg is the reason message.
	Msg string
}

// Errorf returns a formatted reason error.
func Errorf(format string, a ...any) Error {
	return Error{Msg: fmt.Sprintf(format, a...)}
}

// Error return the reason as if it were a message.
// This is used to conform with the error type.
func (e Error) Error() string {
	return e.Msg
}

// Extract removes all reason errors from the error
// chains, returning all other errors and the reason
// messages.
func Extract(err error) ([]string, error) {
	//nolint:errorlint // This is the only way to check for the interface.
	switch x := err.(type) {
	case interface{ Unwrap() []error }:
		var (
			reasons []string
			errs    []error
		)
		for _, err = range x.Unwrap() {
			r, e := Extract(err)
			reasons = append(reasons, r...)
			if e != nil {
				errs = append(errs, e)
			}
		}

		switch len(errs) {
		case 0:
			return reasons, nil
		case 1:
			return reasons, errs[0]
		default:
			return reasons, errors.Join(errs...)
		}
	case interface{ Unwrap() error }:
		return Extract(x.Unwrap())
	default:
		var r Error
		if errors.As(err, &r) {
			return []string{r.Msg}, nil
		}
		return nil, err
	}
}
