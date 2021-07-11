package httpx_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hamba/pkg/v2/httpx"
	"github.com/stretchr/testify/assert"
)

func TestNewHealthHandler(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "healthy",
			wantCode: http.StatusOK,
		},
		{
			name:     "unhealthy",
			err:      errors.New("test"),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			m := &testHealth{err: test.err}
			h := httpx.NewHealthHandler(m)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/health", nil)
			h.ServeHTTP(w, req)

			assert.Equal(t, test.wantCode, w.Code)
		})
	}
}

type testHealth struct {
	err error
}

func (h *testHealth) IsHealthy() error {
	return h.err
}
