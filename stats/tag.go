package stats

import "time"

// TaggedStatter wraps a Statter instance applying tags to all metrics.
type TaggedStatter struct {
	stats Statter
	tags  []string
}

// NewTaggedStatter creates a new TaggedStatter instance.
func NewTaggedStatter(stats Statter, tags ...string) *TaggedStatter {
	if t, ok := stats.(*TaggedStatter); ok {
		stats = t.stats
		tags = append(t.tags, tags...)
	}

	return &TaggedStatter{
		stats: stats,
		tags:  normalizeTags(tags),
	}
}

// Inc increments a count by the value.
func (s TaggedStatter) Inc(name string, value int64, rate float32, tags ...string) {
	s.stats.Inc(name, value, rate, mergeTags(tags, s.tags)...)
}

// Gauge measures the value of a metric.
func (s TaggedStatter) Gauge(name string, value float64, rate float32, tags ...string) {
	s.stats.Gauge(name, value, rate, mergeTags(tags, s.tags)...)
}

// Timing sends the value of a Duration.
func (s TaggedStatter) Timing(name string, value time.Duration, rate float32, tags ...string) {
	s.stats.Timing(name, value, rate, mergeTags(tags, s.tags)...)
}

// Unwrap returns the underlying statter.
func (s TaggedStatter) Unwrap() Statter {
	return s.stats
}

// Close closes the client and flushes buffered stats, if applicable.
func (s TaggedStatter) Close() error {
	return s.stats.Close()
}

func normalizeTags(tags []string) []string {
	// tags need to be even as they are key/value pairs
	if len(tags)%2 != 0 {
		tags = append(tags, "", "STATTER_ERROR", "Normalised odd number of tags by adding an empty string")
	}

	return tags
}

func mergeTags(prefix, suffix []string) []string {
	newTags := make([]string, len(prefix)+len(suffix))

	n := copy(newTags, prefix)
	copy(newTags[n:], suffix)

	return newTags
}
