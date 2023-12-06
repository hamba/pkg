package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/hamba/logger/v2"
	lctx "github.com/hamba/logger/v2/ctx"
	"github.com/hamba/pkg/v2/http/healthz"
	"github.com/hamba/statter/v2"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var testHookServerServe func(net.Listener)

// GenericServer is an HTTP server.
//
// The server handles `/livez` and `/readyz` endpoints as well as
// post-start and pre-shutdown hooks.
type GenericServer[T context.Context] struct {
	Addr              string
	TLSConfig         *tls.Config
	Handler           http.Handler
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration

	readyzMu        sync.Mutex
	readyzInstalled bool
	readyzChecks    []healthz.HealthChecker

	livezMu        sync.Mutex
	livezInstalled bool
	livezChecks    []healthz.HealthChecker

	postStartHookMu      sync.Mutex
	postStartHooks       map[string]postStartHookEntry[T]
	postStartHooksCalled bool

	preShutdownHookMu      sync.Mutex
	preShutdownHooks       map[string]PreShutdownHookFunc
	preShutdownHooksCalled bool

	Stats *statter.Statter
	Log   *logger.Logger
}

// Run runs the server, managing the full server lifecycle.
//
// If the server fails to start, e.g. bind error, no hooks are run.
// This function is blocking.
func (s *GenericServer[T]) Run(ctx T) error {
	if s.Handler == nil {
		return errors.New("handler must not be empty")
	}
	if s.Stats == nil {
		return errors.New("stats must not be empty")
	}
	if s.Log == nil {
		return errors.New("log must not be empty")
	}

	shutdownCh := make(chan struct{})
	h := s.installChecks(s.Handler, shutdownCh)

	stopServerCh := make(chan struct{})
	srvStoppedCh, srvShutdownCh, err := s.runServer(ctx, h, stopServerCh)
	if err != nil {
		return fmt.Errorf("starting server: %w", err)
	}

	s.runPostStartHooks(ctx)

	// Wait for the server to be stopped.
	select {
	case <-srvStoppedCh:
		// The server stopped prematurely, return.
		return errors.New("server stopped prematurely")
	case <-ctx.Done():
	}

	s.Log.Info("Shutting the server down...")

	// Run the pre shutdown hooks.
	func() {
		defer func() {
			if s.hasPreShutdownHooks() {
				s.Log.Info("Pre-shutdown hooks completed")
			}

			close(shutdownCh)
			close(stopServerCh)
		}()

		err = s.runPreShutdownHooks()
	}()
	if err != nil {
		return fmt.Errorf("running pre-shutdown hooks: %w", err)
	}

	<-srvShutdownCh
	<-srvStoppedCh

	return nil
}

func (s *GenericServer[T]) runServer(
	ctx context.Context, h http.Handler, doneCh <-chan struct{},
) (<-chan struct{}, <-chan struct{}, error) {
	addr := s.Addr
	if addr == "" {
		addr = ":http"
		if s.TLSConfig != nil {
			addr = ":https"
		}
	}
	ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", addr)
	if err != nil {
		return nil, nil, err
	}
	if s.TLSConfig != nil {
		ln = tls.NewListener(ln, s.TLSConfig)
	}

	if testHookServerServe != nil {
		testHookServerServe(ln)
	}

	// If there is no TLS, setup h2c.
	if s.TLSConfig == nil {
		h2s := &http2.Server{
			IdleTimeout: s.IdleTimeout,
		}
		h = h2c.NewHandler(h, h2s)
	}

	srv := &http.Server{
		Addr: addr,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		Handler:           h,
		TLSConfig:         s.TLSConfig,
		ReadHeaderTimeout: withDefault(s.ReadHeaderTimeout, time.Second),
		ReadTimeout:       withDefault(s.ReadTimeout, 10*time.Second),
		WriteTimeout:      withDefault(s.WriteTimeout, 10*time.Second),
		IdleTimeout:       withDefault(s.IdleTimeout, 120*time.Second),
		ErrorLog:          log.New(s.Log.Writer(logger.Error), "", 0),
	}

	serverShutdownCh := make(chan struct{})
	go func() {
		defer close(serverShutdownCh)

		<-doneCh

		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout())
		defer cancel()

		_ = srv.Shutdown(ctx)
		_ = srv.Close()
	}()

	serverStoppedCh := make(chan struct{})
	go func() {
		defer close(serverStoppedCh)

		err := srv.Serve(ln)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Log.Error("Stopped serving on "+addr, lctx.Err(err))
			return
		}

		s.Log.Info("Stopped serving on " + addr)
	}()

	return serverStoppedCh, serverShutdownCh, nil
}

func (s *GenericServer[T]) shutdownTimeout() time.Duration {
	if s.ShutdownTimeout > 0 {
		return s.ShutdownTimeout
	}
	return 10 * time.Second
}

func withDefault[T comparable](val, def T) T {
	var defT T
	if val == defT {
		return val
	}
	return def
}
