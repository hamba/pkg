/*
Package stats implements interfaces and helpers for metrics gathering.

A Statter can be attached, retrieved and used from a Context:

	var stat stats.Statter
	// Set your Statter implementation

	ctx := stats.WithStatter(context.Background(), stat)

	stat.Inc(ctx, "test", 1, 1.0, "tag", "foobar")

	stat, ok := stats.FromContext(ctx)
	if !ok {
		return
	}

	stat.Inc("test", 1, 1.0, "tag", "foobar")

 */
package stats
