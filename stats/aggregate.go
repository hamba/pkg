package stats

import (
	"sort"
	"strings"
	"sync"
	"time"
)

// Metric types.
const (
	IncType int = iota
	GaugeType
	TimingType
)

// Metric represents a stats metric.
type Metric struct {
	Type     int
	Hash     string
	Name     string
	IntVal   int64
	FloatVal float64
	DurVal   time.Duration
	Rate     float32
	Tags     []string
}

// Aggregator represents a metric aggregator.
type Aggregator interface {
	// Aggregate aggregates the given metric.
	Aggregate(Metric)

	// Flush flushes the aggregated metrics to the given Statter.
	Flush(Statter)
}

type counterAggregator struct {
	agg map[string]Metric
}

func (a *counterAggregator) Aggregate(metric Metric) {
	if metric.Rate < 1 {
		metric.IntVal = int64(float32(metric.IntVal) / metric.Rate)
	}

	cached, ok := a.agg[metric.Hash]
	if !ok {
		a.agg[metric.Hash] = metric
		return
	}

	cached.IntVal += metric.IntVal
	a.agg[metric.Hash] = cached
}

func (a *counterAggregator) Flush(s Statter) {
	for _, metric := range a.agg {
		s.Inc(metric.Name, metric.IntVal, 1, metric.Tags...)
	}

	a.agg = map[string]Metric{}
}

type gaugeAggregator struct {
	agg map[string]Metric
}

func (a *gaugeAggregator) Aggregate(metric Metric) {
	a.agg[metric.Hash] = metric
}

func (a *gaugeAggregator) Flush(s Statter) {
	for _, metric := range a.agg {
		s.Gauge(metric.Name, metric.FloatVal, 1, metric.Tags...)
	}

	a.agg = map[string]Metric{}
}

type timingAggregator struct {
	agg map[string]Metric
}

func (a *timingAggregator) Aggregate(metric Metric) {
	metric.IntVal = int64(1 / metric.Rate)
	metric.FloatVal = float64(float32(metric.DurVal) / metric.Rate)

	cached, ok := a.agg[metric.Hash]
	if !ok {
		a.agg[metric.Hash] = metric
		return
	}

	cached.IntVal += metric.IntVal
	cached.FloatVal += metric.FloatVal
	a.agg[metric.Hash] = cached
}

func (a *timingAggregator) Flush(s Statter) {
	for _, metric := range a.agg {
		dur := time.Duration(metric.FloatVal / float64(metric.IntVal))
		s.Timing(metric.Name, dur, 1, metric.Tags...)
	}

	a.agg = map[string]Metric{}
}

// AggregateFunc represents a configuration function for an AggregateStatter.
type AggregateFunc func(*AggregateStatter)

// WithCounterAggregator sets the aggregator to use with counters.
func WithCounterAggregator(a Aggregator) AggregateFunc {
	return func(s *AggregateStatter) {
		s.counterAgg = a
	}
}

// WithGaugeAggregator sets the aggregator to use with gauges.
func WithGaugeAggregator(a Aggregator) AggregateFunc {
	return func(s *AggregateStatter) {
		s.gaugeAgg = a
	}
}

// WithTimingAggregator sets the aggregator to use with timings.
func WithTimingAggregator(a Aggregator) AggregateFunc {
	return func(s *AggregateStatter) {
		s.timingAgg = a
	}
}

// AggregateStatter aggregates the incoming stats for a given interval.
//
// By default counters will be summed, gauges will take the last value
// and timings will be averaged.
type AggregateStatter struct {
	stats Statter

	wg sync.WaitGroup
	ch chan Metric

	counterAgg Aggregator
	gaugeAgg   Aggregator
	timingAgg  Aggregator
}

// NewAggregateStatter creates a new aggregate statter, with the given interval.
func NewAggregateStatter(stats Statter, interval time.Duration, opts ...AggregateFunc) *AggregateStatter {
	s := &AggregateStatter{
		stats:      stats,
		ch:         make(chan Metric, 1000),
		counterAgg: &counterAggregator{agg: map[string]Metric{}},
		gaugeAgg:   &gaugeAggregator{agg: map[string]Metric{}},
		timingAgg:  &timingAggregator{agg: map[string]Metric{}},
	}

	for _, opt := range opts {
		opt(s)
	}

	s.wg.Add(1)
	go s.runAggregation(interval)

	return s
}

func (s *AggregateStatter) runAggregation(interval time.Duration) {
	defer s.wg.Done()

	timer := time.NewTicker(interval)
	defer timer.Stop()

	for {
		select {
		case metric, ok := <-s.ch:
			if !ok {
				s.flush()
				return
			}

			switch metric.Type {
			case IncType:
				s.counterAgg.Aggregate(metric)

			case GaugeType:
				s.gaugeAgg.Aggregate(metric)

			case TimingType:
				s.timingAgg.Aggregate(metric)
			}

		case <-timer.C:
			s.flush()
		}
	}
}

func (s *AggregateStatter) flush() {
	s.counterAgg.Flush(s.stats)
	s.gaugeAgg.Flush(s.stats)
	s.timingAgg.Flush(s.stats)
}

// Inc increments a count by the value.
func (s *AggregateStatter) Inc(name string, value int64, rate float32, tags ...string) {
	s.ch <- Metric{
		Type:   IncType,
		Hash:   s.hash(name, tags),
		Name:   name,
		IntVal: value,
		Rate:   rate,
		Tags:   tags,
	}
}

// Gauge measures the value of a Metric.
func (s *AggregateStatter) Gauge(name string, value float64, rate float32, tags ...string) {
	s.ch <- Metric{
		Type:     GaugeType,
		Hash:     s.hash(name, tags),
		Name:     name,
		FloatVal: value,
		Rate:     rate,
		Tags:     tags,
	}
}

// Timing sends the value of a Duration.
func (s *AggregateStatter) Timing(name string, value time.Duration, rate float32, tags ...string) {
	s.ch <- Metric{
		Type:   TimingType,
		Hash:   s.hash(name, tags),
		Name:   name,
		DurVal: value,
		Rate:   rate,
		Tags:   tags,
	}
}

func (s *AggregateStatter) hash(name string, tags []string) string {
	tg := make([]string, len(tags), len(tags)+1)
	copy(tg, tags)
	sort.Strings(tg)
	tg = append(tg, name)
	return strings.Join(tg, "")
}

// Close closes the client and flushes aggregated stats.
func (s *AggregateStatter) Close() error {
	close(s.ch)

	s.wg.Wait()

	return s.stats.Close()
}
