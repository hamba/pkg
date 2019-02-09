package stats_test

import (
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
	var statable stats.Statable
	// Set your Statter implementation

	stats.Group(statable, "prefix", func(s stats.Statter) {
		s.Inc("test", 1, 1.0, "tag", "foobar")
	})

	// Output name: "prefix.test"
}

func ExampleTime() {
	var statable stats.Statable
	// Set your Statter implementation

	timer := stats.Time(statable, "latency", 1.0, "tag", "foobar")

	// Do something

	timer.Done()

	// Output mertic: stats.Timing(ctx, "latency", duration, 1.0, "tag", "foobar")
}
