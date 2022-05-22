package middleware_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hamba/logger/v2"
	"github.com/hamba/pkg/v2/http/middleware"
	"github.com/hamba/statter/v2"
	"github.com/hamba/statter/v2/reporter/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWithRecovery(t *testing.T) {
	tests := []struct {
		name    string
		val     interface{}
		wantLog string
	}{
		{
			name:    "with string",
			val:     "panic text",
			wantLog: `lvl=eror msg="Panic while serving request" err="panic text" method=GET url=/ stack=`,
		},
		{
			name:    "with error",
			val:     errors.New("test error"),
			wantLog: `lvl=eror msg="Panic while serving request" err="test error" method=GET url=/ stack=`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			log := logger.New(&buf, logger.LogfmtFormat(), logger.Info)

			h := middleware.WithRecovery(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic(test.val)
				}),
				log,
			)

			req, _ := http.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()

			defer func() {
				if err := recover(); err != nil {
					t.Fatal("Expected the panic to be handled.")
				}
			}()

			h.ServeHTTP(resp, req)

			assert.Contains(t, buf.String(), test.wantLog)
		})
	}

}

func TestWithStats(t *testing.T) {
	tests := []struct {
		name        string
		handlerName string
		path        string
		wantTags    [][2]string
	}{
		{
			name:        "with handler name",
			handlerName: "my-handler",
			path:        "/test",
			wantTags:    [][2]string{{"handler", "my-handler"}},
		},
		{
			name:     "without handler name",
			path:     "/test",
			wantTags: [][2]string{{"handler", "/test"}},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			m := &mockReporter{}
			m.On("Counter", "requests", int64(1), test.wantTags)
			wantTags := append([][2]string{{"code", "3xx"}}, test.wantTags...)
			m.On("Counter", "responses", int64(1), wantTags)
			m.On("Histogram", "response.size", wantTags).Return(func(_ float64) {})
			m.On("Timing", "response.duration", wantTags).Return(func(_ time.Duration) {})

			s := statter.New(m, time.Second)

			h := middleware.WithStats(test.handlerName, s, http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(305)
				}),
			)

			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", test.path, nil)

			h.ServeHTTP(resp, req)

			err := s.Close()
			require.NoError(t, err)

			m.AssertExpectations(t)
		})
	}
}

func TestWithStats_Prometheus(t *testing.T) {
	reporter := prometheus.New("test")
	s := statter.New(reporter, time.Second)

	h := middleware.WithStats("test-handler", s, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(305)
		}),
	)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	h.ServeHTTP(resp, req)

	err := s.Close()
	require.NoError(t, err)
}

type mockReporter struct {
	mock.Mock
}

func (r *mockReporter) Counter(name string, v int64, tags [][2]string) {
	_ = r.Called(name, v, tags)
}

func (r *mockReporter) Gauge(name string, v float64, tags [][2]string) {
	_ = r.Called(name, v, tags)
}

func (r *mockReporter) Histogram(name string, tags [][2]string) func(v float64) {
	args := r.Called(name, tags)
	fn := args.Get(0)
	if fn == nil {
		return nil
	}
	return fn.(func(float64))
}

func (r *mockReporter) Timing(name string, tags [][2]string) func(v time.Duration) {
	args := r.Called(name, tags)
	fn := args.Get(0)
	if fn == nil {
		return nil
	}
	return fn.(func(time.Duration))
}
