package stats_test

import (
	"context"

	"github.com/hamba/pkg/stats"
)

func ExampleTaggedStatter() {
	var stat stats.Statter
	// Set your Statter implementation

	stat = stats.NewTaggedStatter(stat, "app_env", "stag")

	stat.Inc("test", 1, 1.0, "tag", "foobar")

	// Output tags: "tag", "foobar", "app_env", "stag"
}

func ExampleGroup() {
	var stat stats.Statter
	// Set your Statter implementation

	ctx := stats.WithStatter(context.Background(), stat)

	stats.Group(ctx, "prefix", func(s stats.Statter) {
		stat.Inc("test", 1, 1.0, "tag", "foobar")
	})

	// Output name: "prefix.test"
}

func ExampleTime() {
	var stat stats.Statter
	// Set your Statter implementation

	ctx := stats.WithStatter(context.Background(), stat)

	timer := stats.Time(ctx, "latency", 1.0, "tag", "foobar")

	// Do something

	timer.Done()

	// Output mertic: stats.Timing(ctx, "latency", duration, 1.0, "tag", "foobar")
}
