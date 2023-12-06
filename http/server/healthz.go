package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	lctx "github.com/hamba/logger/v2/ctx"
	"github.com/hamba/pkg/v2/http/healthz"
	"github.com/hamba/pkg/v2/http/middleware"
)

// MustAddHealthzChecks adds health checks to both readyz and livez, panicking if there is an error.
func (s *GenericServer[T]) MustAddHealthzChecks(checks ...healthz.HealthChecker) {
	if err := s.AddHealthzChecks(checks...); err != nil {
		panic(err)
	}
}

// AddHealthzChecks adds health checks to both readyz and livez.
func (s *GenericServer[T]) AddHealthzChecks(checks ...healthz.HealthChecker) error {
	if err := s.AddReadyzChecks(checks...); err != nil {
		return err
	}
	return s.AddLivezChecks(checks...)
}

// MustAddReadyzChecks adds health checks to readyz, panicking if there is an error.
func (s *GenericServer[T]) MustAddReadyzChecks(checks ...healthz.HealthChecker) {
	if err := s.AddReadyzChecks(checks...); err != nil {
		panic(err)
	}
}

// AddReadyzChecks adds health checks to readyz.
func (s *GenericServer[T]) AddReadyzChecks(checks ...healthz.HealthChecker) error {
	s.readyzMu.Lock()
	defer s.readyzMu.Unlock()
	if s.readyzInstalled {
		return errors.New("could not add checks as readyz has already been installed")
	}
	s.readyzChecks = append(s.readyzChecks, checks...)
	return nil
}

// MustAddLivezChecks adds health checks to livez, panicking if there is an error.
func (s *GenericServer[T]) MustAddLivezChecks(checks ...healthz.HealthChecker) {
	if err := s.AddLivezChecks(checks...); err != nil {
		panic(err)
	}
}

// AddLivezChecks adds health checks to livez.
func (s *GenericServer[T]) AddLivezChecks(checks ...healthz.HealthChecker) error {
	s.livezMu.Lock()
	defer s.livezMu.Unlock()
	if s.livezInstalled {
		return errors.New("could not add checks as livez has already been installed")
	}
	s.livezChecks = append(s.livezChecks, checks...)
	return nil
}

func (s *GenericServer[T]) installChecks(h http.Handler, shutdownCh chan struct{}) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", h)
	s.installLivezChecks(mux)

	// When shutdown is started, the readyz check should start failing.
	if err := s.AddReadyzChecks(shutdownCheck{ch: shutdownCh}); err != nil {
		s.Log.Error("Could not install readyz shutdown check", lctx.Err(err))
	}
	s.installReadyzChecks(mux)

	return mux
}

func (s *GenericServer[T]) installReadyzChecks(mux *http.ServeMux) {
	s.readyzMu.Lock()
	defer s.readyzMu.Unlock()
	s.readyzInstalled = true
	s.installCheckers(mux, "/readyz", s.readyzChecks)
}

func (s *GenericServer[T]) installLivezChecks(mux *http.ServeMux) {
	s.livezMu.Lock()
	defer s.livezMu.Unlock()
	s.livezInstalled = true
	s.installCheckers(mux, "/livez", s.livezChecks)
}

func (s *GenericServer[T]) installCheckers(mux *http.ServeMux, path string, checks []healthz.HealthChecker) {
	if len(checks) == 0 {
		checks = []healthz.HealthChecker{healthz.PingHealth}
	}

	s.Log.Info("Installing health checkers",
		lctx.Str("path", path),
		lctx.Str("checks", strings.Join(checkNames(checks), ",")),
	)

	name := strings.TrimPrefix(path, "/")
	h := healthz.Handler(name, func(output string) {
		s.Log.Info(fmt.Sprintf("%s check failed\n%s", name, output))
	}, checks...)
	mux.Handle(path, middleware.WithStats(name, s.Stats, h))
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
