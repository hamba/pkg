package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hamba/pkg/httpx"
	"github.com/stretchr/testify/assert"
)

func TestNewStatsMux(t *testing.T) {
	mux := httpx.NewStatsMux(&statsHandler{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/metrics", nil)
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewStatsMux_NoHandler(t *testing.T) {
	mux := httpx.NewStatsMux(&testStatter{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/metrics", nil)
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

type statsHandler struct{}

func (h statsHandler) Inc(name string, value int64, rate float32, tags ...string) {}

func (h statsHandler) Gauge(name string, value float64, rate float32, tags ...string) {}

func (h statsHandler) Timing(name string, value time.Duration, rate float32, tags ...string) {}

func (h statsHandler) Handler() http.Handler {
	return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
}

func (h statsHandler) Close() error {
	return nil
}

type testStatter struct{}

func (s testStatter) Inc(name string, value int64, rate float32, tags ...string) {}

func (s testStatter) Gauge(name string, value float64, rate float32, tags ...string) {}

func (s testStatter) Timing(name string, value time.Duration, rate float32, tags ...string) {}

func (s testStatter) Close() error {
	return nil
}
