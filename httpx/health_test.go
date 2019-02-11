package httpx_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hamba/pkg/httpx"
	"github.com/stretchr/testify/assert"
)

func TestNewHealthMux_Healthy(t *testing.T) {
	mux := httpx.NewHealthMux(&healthy{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewHealthMux_Unhealthy(t *testing.T) {
	mux := httpx.NewHealthMux(&unhealthy{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

type healthy struct{}

func (h *healthy) IsHealthy() error {
	return nil
}

type unhealthy struct{}

func (h *unhealthy) IsHealthy() error {
	return errors.New("test error")
}
