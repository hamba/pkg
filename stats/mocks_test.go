package stats_test

import (
	"time"

	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/mock"
)

type testStatable struct {
	s stats.Statter
}

func (s *testStatable) Statter() stats.Statter {
	return s.s
}

type MockStats struct {
	mock.Mock
}

func (m *MockStats) Inc(name string, value int64, rate float32, tags ...string) {
	m.Called(name, value, rate, tags)
}

func (m *MockStats) Gauge(name string, value float64, rate float32, tags ...string) {
	m.Called(name, value, rate, tags)
}

func (m *MockStats) Timing(name string, value time.Duration, rate float32, tags ...string) {
	m.Called(name, value, rate, tags)
}

func (m *MockStats) Close() error {
	args := m.Called()

	return args.Error(0)
}
