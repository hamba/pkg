package request

import "context"

type contextKey int

const requestID contextKey = iota + 1

// WithID returns a copy of parent in which the request ID value is set.
func WithID(parent context.Context, id string) context.Context {
	return context.WithValue(parent, requestID, id)
}

// IDFrom returns the value of the request ID on the ctx.
func IDFrom(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestID).(string)
	return id, ok
}
