package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hamba/pkg/httpx/middleware"
	"github.com/hamba/pkg/log"
)

func TestWithRecovery(t *testing.T) {
	h := middleware.WithRecovery(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("panic")
		}),
		log.NewMockLoggable(log.Null),
	)

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	h.ServeHTTP(resp, req)
}

func TestWithRecovery_Error(t *testing.T) {
	h := middleware.WithRecovery(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(errors.New("panic"))
		}),
		log.NewMockLoggable(log.Null),
	)

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	defer func() {
		if err := recover(); err != nil {
			t.Fatal("Expected the panic to be handled.")
		}
	}()

	h.ServeHTTP(resp, req)
}
