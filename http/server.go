package http

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// SrvOptFunc represents a server option function.
type SrvOptFunc func(*http.Server)

// WithTLSConfig sets the serve tls config.
func WithTLSConfig(cfg *tls.Config) SrvOptFunc {
	return func(srv *http.Server) {
		srv.TLSConfig = cfg
	}
}

// WithReadTimeout sets the server read timeout.
func WithReadTimeout(d time.Duration) SrvOptFunc {
	return func(srv *http.Server) {
		srv.ReadTimeout = d
	}
}

// WithWriteTimeout sets the server write timeout.
func WithWriteTimeout(d time.Duration) SrvOptFunc {
	return func(srv *http.Server) {
		srv.WriteTimeout = d
	}
}

// WithH2C allows the server to handle h2c connections.
func WithH2C() SrvOptFunc {
	return func(srv *http.Server) {
		h2s := &http2.Server{
			IdleTimeout: 120 * time.Second,
		}

		srv.Handler = h2c.NewHandler(srv.Handler, h2s)
	}
}

// Server is a convenience wrapper around the standard
// library http server.
type Server struct {
	srv *http.Server
}

// NewServer returns a server.
func NewServer(ctx context.Context, addr string, h http.Handler, opts ...SrvOptFunc) *Server {
	srv := &http.Server{
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return &Server{
		srv: srv,
	}
}

// Serve starts the server in a in-blocking way.
func (s *Server) Serve(errFn func(error)) {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errFn(err)
		}
	}()
}

// Shutdown attempts to close all server connections.
func (s *Server) Shutdown(timeout time.Duration) error {
	stopCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return s.srv.Shutdown(stopCtx)
}

// Close closes the server.
func (s *Server) Close() error {
	return s.srv.Close()
}
