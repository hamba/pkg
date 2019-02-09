/*
Package log implements interfaces and helpers for logging.

A Logger can be attached, retrieved and used from a Context:

	var l log.Logger
	// Set your Logger implementation

	ctx := log.WithLogger(context.Background(), l)

	log.Info(ctx, "message", "context", "info")

	l, ok := log.FromContext(ctx)
	if !ok {
		return
	}

	log.Info("message", "context", "info")

*/
package log
