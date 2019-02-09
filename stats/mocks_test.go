package stats_test

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockStats struct {
	mock.Mock
}

func (m *MockStats) Inc(name string, value int64, rate float32, tags ...interface{}) {
	m.Called(name, value, rate, tags)
}

func (m *MockStats) Gauge(name string, value float64, rate float32, tags ...interface{}) {
	m.Called(name, value, rate, tags)
}

func (m *MockStats) Timing(name string, value time.Duration, rate float32, tags ...interface{}) {
	m.Called(name, value, rate, tags)
}

func (m *MockStats) Close() error {
	args := m.Called()

	return args.Error(0)
}
