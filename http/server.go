package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/hamba/logger/v2"
	lctx "github.com/hamba/logger/v2/ctx"
	"github.com/hamba/pkg/v2/http/healthz"
	"github.com/hamba/pkg/v2/http/middleware"
	"github.com/hamba/statter/v2"
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

var testHookServerServe func(net.Listener)

// Server is a convenience wrapper around the standard
// library HTTP server.
type Server struct {
	srv *http.Server
}

// NewServer returns a server with the base context ctx.
func NewServer(ctx context.Context, addr string, h http.Handler, opts ...SrvOptFunc) *Server {
	srv := &http.Server{
		BaseContext: func(ln net.Listener) context.Context {
			if testHookServerServe != nil {
				testHookServerServe(ln)
			}
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

// Serve starts the server in a non-blocking way.
func (s *Server) Serve(errFn func(error)) {
	go func() {
		if s.srv.TLSConfig != nil {
			if err := s.srv.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errFn(err)
			}
			return
		}

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

// HealthServerConfig configures a HealthServer.
type HealthServerConfig struct {
	Addr    string
	Handler http.Handler

	ReadyzChecks []healthz.HealthChecker
	LivezChecks  []healthz.HealthChecker

	Stats *statter.Statter
	Log   *logger.Logger
}

// AddHealthzChecks adds the given checks to the config.
func (c *HealthServerConfig) AddHealthzChecks(checks ...healthz.HealthChecker) {
	c.AddReadyzChecks(checks...)
	c.AddLivezChecks(checks...)
}

// AddReadyzChecks adds the given checks to the config.
func (c *HealthServerConfig) AddReadyzChecks(checks ...healthz.HealthChecker) {
	c.ReadyzChecks = append(c.ReadyzChecks, checks...)
}

// AddLivezChecks adds the given checks to the config.
func (c *HealthServerConfig) AddLivezChecks(checks ...healthz.HealthChecker) {
	c.LivezChecks = append(c.LivezChecks, checks...)
}

// HealthServer is an HTTP server with healthz capabilities.
type HealthServer struct {
	srv *Server
	mux *http.ServeMux

	shudownCh chan struct{}

	readyzMu        sync.Mutex
	readyzInstalled bool
	readyzChecks    []healthz.HealthChecker

	livezMu        sync.Mutex
	livezInstalled bool
	livezChecks    []healthz.HealthChecker

	stats *statter.Statter
	log   *logger.Logger
}

// NewHealthServer returns an HTTP server with healthz capabilities.
func NewHealthServer(ctx context.Context, cfg HealthServerConfig, opts ...SrvOptFunc) *HealthServer {
	// Setup the mux early so H2C can attach properly.
	mux := http.NewServeMux()
	mux.Handle("/", cfg.Handler)

	srv := NewServer(ctx, cfg.Addr, mux, opts...)

	return &HealthServer{
		srv:          srv,
		mux:          mux,
		shudownCh:    make(chan struct{}),
		readyzChecks: cfg.ReadyzChecks,
		livezChecks:  cfg.LivezChecks,
		stats:        cfg.Stats,
		log:          cfg.Log,
	}
}

// AddHealthzChecks adds health checks to both readyz and livez.
func (s *HealthServer) AddHealthzChecks(checks ...healthz.HealthChecker) error {
	if err := s.AddReadyzChecks(checks...); err != nil {
		return err
	}
	return s.AddLivezChecks(checks...)
}

// AddReadyzChecks adds health checks to readyz.
func (s *HealthServer) AddReadyzChecks(checks ...healthz.HealthChecker) error {
	s.readyzMu.Lock()
	defer s.readyzMu.Unlock()
	if s.readyzInstalled {
		return errors.New("could not add checks as readyz has already been installed")
	}
	s.readyzChecks = append(s.readyzChecks, checks...)
	return nil
}

// AddLivezChecks adds health checks to livez.
func (s *HealthServer) AddLivezChecks(checks ...healthz.HealthChecker) error {
	s.livezMu.Lock()
	defer s.livezMu.Unlock()
	if s.livezInstalled {
		return errors.New("could not add checks as livez has already been installed")
	}
	s.livezChecks = append(s.livezChecks, checks...)
	return nil
}

// Serve installs the health checks and starts the server in a non-blocking way.
func (s *HealthServer) Serve(errFn func(error)) {
	s.installChecks()

	s.srv.Serve(errFn)
}

func (s *HealthServer) installChecks() {
	s.installLivezChecks(s.mux)

	// When shutdown is started, the readyz check should start failing.
	if err := s.AddReadyzChecks(shutdownCheck{ch: s.shudownCh}); err != nil {
		s.log.Error("Could not install readyz shutdown check", lctx.Err(err))
	}
	s.installReadyzChecks(s.mux)
}

func (s *HealthServer) installReadyzChecks(mux *http.ServeMux) {
	s.readyzMu.Lock()
	defer s.readyzMu.Unlock()
	s.readyzInstalled = true
	s.installCheckers(mux, "/readyz", s.readyzChecks)
}

func (s *HealthServer) installLivezChecks(mux *http.ServeMux) {
	s.livezMu.Lock()
	defer s.livezMu.Unlock()
	s.livezInstalled = true
	s.installCheckers(mux, "/livez", s.livezChecks)
}

func (s *HealthServer) installCheckers(mux *http.ServeMux, path string, checks []healthz.HealthChecker) {
	if len(checks) == 0 {
		checks = []healthz.HealthChecker{healthz.PingHealth}
	}

	s.log.Info("Installing health checkers",
		lctx.Str("path", path),
		lctx.Str("checks", strings.Join(checkNames(checks), ",")),
	)

	name := strings.TrimPrefix(path, "/")
	h := healthz.Handler(name, func(output string) {
		s.log.Info(fmt.Sprintf("%s check failed\n%s", name, output))
	}, checks...)
	mux.Handle(path, middleware.WithStats(name, s.stats, h))
}

// Shutdown attempts to close all server connections.
func (s *HealthServer) Shutdown(timeout time.Duration) error {
	close(s.shudownCh)

	return s.srv.Shutdown(timeout)
}

// Close closes the server.
func (s *HealthServer) Close() error {
	return s.srv.Close()
}

func checkNames(checks []healthz.HealthChecker) []string {
	names := make([]string, len(checks))
	for i, check := range checks {
		names[i] = check.Name()
	}
	return names
}

type shutdownCheck struct {
	ch <-chan struct{}
}

func (s shutdownCheck) Name() string { return "shutdown" }

func (s shutdownCheck) Check(*http.Request) error {
	select {
	case <-s.ch:
		return errors.New("server is shutting down")
	default:
		return nil
	}
}
