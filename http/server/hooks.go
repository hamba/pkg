package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	lctx "github.com/hamba/logger/v2/ctx"
)

// PostStartHookFunc is a function called after server start.
type PostStartHookFunc[T context.Context] func(T) error

// PreShutdownHookFunc is a function called before server shutdown.
type PreShutdownHookFunc func() error

type postStartHookEntry[T context.Context] struct {
	fn     PostStartHookFunc[T]
	doneCh chan struct{}
}

// MustAddPostStartHook adds a post-start hook, panicking if there is an error.
func (s *GenericServer[T]) MustAddPostStartHook(name string, fn PostStartHookFunc[T]) {
	if err := s.AddPostStartHook(name, fn); err != nil {
		panic(err)
	}
}

// AddPostStartHook adds a post-start hook.
func (s *GenericServer[T]) AddPostStartHook(name string, fn PostStartHookFunc[T]) error {
	if name == "" {
		return errors.New("name is required")
	}
	if fn == nil {
		return errors.New("fn is required")
	}

	s.postStartHookMu.Lock()
	defer s.postStartHookMu.Unlock()

	if s.postStartHooksCalled {
		return errors.New("hooks have already been called")
	}
	if _, exists := s.postStartHooks[name]; exists {
		return fmt.Errorf("hook %q as it is already registered", name)
	}

	if s.postStartHooks == nil {
		s.postStartHooks = map[string]postStartHookEntry[T]{}
	}

	doneCh := make(chan struct{})
	err := s.AddReadyzChecks(postStartHookHealth{
		name:   "postStartHook:" + name,
		doneCh: doneCh,
	})
	if err != nil {
		return fmt.Errorf("adding readyz check: %w", err)
	}

	s.postStartHooks[name] = postStartHookEntry[T]{
		fn:     fn,
		doneCh: doneCh,
	}
	return nil
}

// MustAddPreShutdownHook adds a pre-shutdown hook, panicking if there is an error.
func (s *GenericServer[T]) MustAddPreShutdownHook(name string, fn PreShutdownHookFunc) {
	if err := s.AddPreShutdownHook(name, fn); err != nil {
		panic(err)
	}
}

// AddPreShutdownHook adds a pre-shutdown hook.
func (s *GenericServer[T]) AddPreShutdownHook(name string, fn PreShutdownHookFunc) error {
	if name == "" {
		return errors.New("name is required")
	}
	if fn == nil {
		return errors.New("fn is required")
	}

	s.preShutdownHookMu.Lock()
	defer s.preShutdownHookMu.Unlock()

	if s.preShutdownHooksCalled {
		return errors.New("hooks have already been called")
	}
	if _, exists := s.preShutdownHooks[name]; exists {
		return fmt.Errorf("hook %q as it is already registered", name)
	}

	if s.preShutdownHooks == nil {
		s.preShutdownHooks = map[string]PreShutdownHookFunc{}
	}

	s.preShutdownHooks[name] = fn
	return nil
}

func (s *GenericServer[T]) runPostStartHooks(ctx T) {
	s.postStartHookMu.Lock()
	defer s.postStartHookMu.Unlock()

	s.postStartHooksCalled = true

	for name, entry := range s.postStartHooks {
		go s.runPostStartHook(ctx, name, entry)
	}
}

func (s *GenericServer[T]) runPostStartHook(ctx T, name string, entry postStartHookEntry[T]) {
	defer func() {
		if v := recover(); v != nil {
			s.Log.Error("Panic while running post-start hook",
				lctx.Interface("error", v),
				lctx.Stack("stack"),
			)
		}
	}()

	s.Log.Info("Running post-start hook", lctx.Str("hook", name))

	if err := entry.fn(ctx); err != nil {
		s.Log.Error("Could not run post-start hook", lctx.Str("name", name), lctx.Err(err))
	}
	close(entry.doneCh)
}

func (s *GenericServer[T]) hasPreShutdownHooks() bool {
	s.preShutdownHookMu.Lock()
	defer s.preShutdownHookMu.Unlock()

	return len(s.preShutdownHooks) > 0
}

func (s *GenericServer[T]) runPreShutdownHooks() error {
	s.preShutdownHookMu.Lock()
	defer s.preShutdownHookMu.Unlock()

	s.preShutdownHooksCalled = true

	var errs error
	for name, fn := range s.preShutdownHooks {
		if err := s.runPreShutdownHook(name, fn); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (s *GenericServer[T]) runPreShutdownHook(name string, fn PreShutdownHookFunc) error {
	defer func() {
		if v := recover(); v != nil {
			s.Log.Error("Panic while running pre-shutdown hook",
				lctx.Interface("error", v),
				lctx.Stack("stack"),
			)
		}
	}()

	s.Log.Info("Running pre-shutdown hook", lctx.Str("hook", name))

	if err := fn(); err != nil {
		return fmt.Errorf("running preshutdown hook %q: %w", name, err)
	}
	return nil
}

type postStartHookHealth struct {
	name   string
	doneCh chan struct{}
}

func (h postStartHookHealth) Name() string {
	return h.name
}

func (h postStartHookHealth) Check(*http.Request) error {
	select {
	case <-h.doneCh:
		return nil
	default:
		return errors.New("not finished")
	}
}
