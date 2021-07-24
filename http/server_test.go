package http_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	httpx "github.com/hamba/pkg/v2/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	var handlerCalled bool
	h := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handlerCalled = true
	})

	srv := httpx.NewServer(context.Background(), ":16543", h)
	srv.Serve(func(err error) {})
	defer srv.Close()

	res, err := http.DefaultClient.Get("http://localhost:16543/")
	require.NoError(t, err)
	defer res.Body.Close()

	srv.Shutdown(time.Second)

	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.True(t, handlerCalled)
}
