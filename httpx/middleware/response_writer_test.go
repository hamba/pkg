package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/hamba/pkg/httpx/middleware"
	"github.com/stretchr/testify/assert"
)

func TestResponseWriter_Status(t *testing.T) {
	rw := middleware.NewResponseWriter(httptest.NewRecorder())

	assert.Equal(t, 200, rw.Status())

	rw.WriteHeader(123)

	assert.Equal(t, 123, rw.Status())
	assert.Equal(t, int64(0), rw.BytesWritten())
}

func TestResponseWriter_WriteStatus(t *testing.T) {
	rw := middleware.NewResponseWriter(httptest.NewRecorder())

	assert.Equal(t, 200, rw.Status())

	rw.Write([]byte{0, 1, 2, 3})

	assert.Equal(t, 200, rw.Status())
	assert.Equal(t, int64(4), rw.BytesWritten())
}
