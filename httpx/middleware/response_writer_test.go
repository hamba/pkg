package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/msales/pkg/v3/httpx/middleware"
	"github.com/stretchr/testify/assert"
)

func TestResponseWriter_Status(t *testing.T) {
	rw := middleware.NewResponseWriter(httptest.NewRecorder())

	assert.Equal(t, 0, rw.Status())

	rw.WriteHeader(123)

	assert.Equal(t, 123, rw.Status())
}

func TestResponseWriter_WriteStatus(t *testing.T) {
	rw := middleware.NewResponseWriter(httptest.NewRecorder())

	assert.Equal(t, 0, rw.Status())

	rw.Write([]byte{})

	assert.Equal(t, 200, rw.Status())
}
