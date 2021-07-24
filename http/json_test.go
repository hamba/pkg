package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	httpx "github.com/hamba/pkg/v2/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		data     interface{}
		wantJSON string
		wantErr  require.ErrorAssertionFunc
	}{
		{
			name: "encodes JSON",
			code: http.StatusOK,
			data: struct {
				Foo string
				Bar string
			}{"foo", "bar"},
			wantJSON: `{"Foo":"foo","Bar":"bar"}`,
			wantErr:  require.NoError,
		},
		{
			name:    "handles bad data",
			code:    http.StatusInternalServerError,
			data:    make(chan int),
			wantErr: require.Error,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			err := httpx.JSON(w, test.code, test.data)

			test.wantErr(t, err)
			assert.Equal(t, test.code, w.Code)
			assert.Equal(t, test.wantJSON, string(w.Body.Bytes()))

			if test.code/100 == 2 {
				assert.Equal(t, httpx.JSONContentType, w.Header().Get("Content-Type"))
			}
		})
	}
}

func TestJSON_WriteError(t *testing.T) {
	w := testResponseWriter{}

	err := httpx.JSON(w, 200, "test")

	assert.Error(t, err)
}

type testResponseWriter struct{}

func (rw testResponseWriter) Header() http.Header       { return http.Header{} }
func (rw testResponseWriter) Write([]byte) (int, error) { return 0, errors.New("test error") }
func (rw testResponseWriter) WriteHeader(int)           {}
