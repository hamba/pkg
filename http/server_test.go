package http_test

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/hamba/logger/v2"
	httpx "github.com/hamba/pkg/v2/http"
	"github.com/hamba/pkg/v2/http/healthz"
	"github.com/hamba/statter/v2"
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

func TestHealthServer(t *testing.T) {
	stats := statter.New(statter.DiscardReporter, time.Minute)
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Error)

	check := healthz.NamedCheck("test", func(*http.Request) error {
		return nil
	})

	srv := httpx.NewHealthServer(context.Background(), httpx.HealthServerConfig{
		Addr:    ":16543",
		Handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
		Stats:   stats,
		Log:     log,
	})
	err := srv.AddHealthzChecks(check)
	require.NoError(t, err)
	srv.Serve(func(err error) {
		require.NoError(t, err)
	})
	t.Cleanup(func() {
		_ = srv.Close()
	})

	statusCode, body := requireDoRequest(t, "http://localhost:16543/readyz?verbose=1")
	assert.Equal(t, statusCode, http.StatusOK)
	want := `+ test ok
+ shutdown ok
readyz check passed`
	assert.Equal(t, want, body)

	statusCode, body = requireDoRequest(t, "http://localhost:16543/livez?verbose=1")
	assert.Equal(t, statusCode, http.StatusOK)
	want = `+ test ok
livez check passed`
	assert.Equal(t, want, body)
}

func TestHealthServer_ShutdownCausesReadyzCheckToFail(t *testing.T) {
	stats := statter.New(statter.DiscardReporter, time.Minute)
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Error)

	calledCh := make(chan struct{})
	check := healthz.NamedCheck("test", func(*http.Request) error {
		close(calledCh)
		time.Sleep(time.Millisecond)
		return nil
	})

	srv := httpx.NewHealthServer(context.Background(), httpx.HealthServerConfig{
		Addr:    ":16543",
		Handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
		Stats:   stats,
		Log:     log,
	})
	err := srv.AddHealthzChecks(check)
	require.NoError(t, err)
	srv.Serve(func(err error) {
		require.NoError(t, err)
	})
	t.Cleanup(func() {
		_ = srv.Close()
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		<-calledCh

		err = srv.Shutdown(time.Second)
		require.NoError(t, err)
	}()

	statusCode, body := requireDoRequest(t, "http://localhost:16543/readyz?verbose=1")
	assert.Equal(t, statusCode, http.StatusInternalServerError)
	want := `+ test ok
- shutdown failed
readyz check failed
`
	assert.Equal(t, want, body)

	wg.Wait()
}

func requireDoRequest(t *testing.T, path string) (int, string) {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, string(b)
}
