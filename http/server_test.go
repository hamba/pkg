package http_test

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	httpx "github.com/hamba/pkg/v2/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
)

func TestServer(t *testing.T) {
	var handlerCalled bool
	h := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})

	srv := httpx.NewServer(context.Background(), ":16543", h)
	srv.Serve(func(err error) {
		require.NoError(t, err)
	})
	t.Cleanup(func() {
		_ = srv.Close()
	})

	res, err := http.DefaultClient.Get("http://localhost:16543/")
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	})

	err = srv.Shutdown(time.Second)
	require.NoError(t, err)

	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.True(t, handlerCalled)
}

func TestServer_WithH2C(t *testing.T) {
	var handlerCalled bool
	h := http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		assert.True(t, req.ProtoAtLeast(2, 0))

		handlerCalled = true
	})

	srv := httpx.NewServer(context.Background(), ":16543", h, httpx.WithH2C())
	srv.Serve(func(err error) {
		require.NoError(t, err)
	})
	t.Cleanup(func() {
		_ = srv.Close()
	})

	c := &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(_ context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
	res, err := c.Get("http://localhost:16543/")
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	})

	err = srv.Shutdown(time.Second)
	require.NoError(t, err)

	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.True(t, handlerCalled)
}
