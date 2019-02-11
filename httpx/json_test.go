package httpx

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSONResponse(t *testing.T) {
	tests := []struct {
		code         int
		data         interface{}
		expectedJSON string
		expectedErr  bool
	}{
		{
			200,
			struct {
				Foo string
				Bar string
			}{"foo", "bar"},
			`{"Foo":"foo","Bar":"bar"}`,
			false,
		},
		{
			500,
			make(chan int),
			"",
			true,
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		err := WriteJSONResponse(w, test.code, test.data)

		assert.Equal(t, test.code, w.Code)
		assert.Equal(t, test.expectedJSON, string(w.Body.Bytes()))
		if test.expectedErr {
			assert.Error(t, err)
			continue

		}
		assert.NoError(t, err)
		assert.Equal(t, JSONContentType, w.Header().Get("Content-Type"))
	}
}

func TestWriteJSONResponse_WriteError(t *testing.T) {
	w := FakeResponseWriter{}

	err := WriteJSONResponse(w, 200, "test")

	assert.Error(t, err)
}

type FakeResponseWriter struct{}

func (rw FakeResponseWriter) Header() http.Header {
	return http.Header{}
}

func (rw FakeResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("test error")
}

func (rw FakeResponseWriter) WriteHeader(int) {

}
