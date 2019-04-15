package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hamba/pkg/httpx/middleware"
	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/mock"
)

func TestWithRequestStats(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		tagFuncs []middleware.TagsFunc
		wantTags []string
	}{
		{
			name:     "With Default Tags",
			path:     "/test",
			tagFuncs: nil,
			wantTags: []string{"method", "GET", "path", "/test"},
		},
		{
			name:     "With Custom Tags",
			path:     "/test",
			tagFuncs: []middleware.TagsFunc{testTags},
			wantTags: []string{"method", "GET"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := new(MockStats)
			s.On("Inc", "request.start", int64(1), float32(1.0), tt.wantTags)
			s.On("Timing", "request.time", mock.Anything, float32(1.0), tt.wantTags)
			s.On("Inc", "request.complete", int64(1), float32(1.0), mock.Anything)

			m := middleware.WithRequestStats(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
				}),
				stats.NewMockStatable(s),
				tt.tagFuncs...,
			)

			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.path, nil)

			m.ServeHTTP(resp, req)

			s.AssertExpectations(t)
		})
	}
}

func testTags(r *http.Request) []string {
	return []string{
		"method", r.Method,
	}
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
