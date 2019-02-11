package httpx_test

import (
	"net/http"
	"testing"

	"github.com/go-zoo/bone"
	"github.com/hamba/pkg/httpx"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	m := httpx.NewMux()
	m.Get("/test", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))

	srv := httpx.NewServer("127.0.0.1:65234", m)

	assert.IsType(t, &http.Server{}, srv)
	assert.Equal(t, m, srv.Handler)
}

func TestNewServer_MultipleMuxes(t *testing.T) {
	h1 := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	m1 := httpx.NewMux()
	m1.Get("/test", h1)

	h2 := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	m2 := httpx.NewMux()
	m2.Post("/foobar", h2)

	srv := httpx.NewServer("127.0.0.1:65234", m1, m2)

	assert.Len(t, srv.Handler.(*bone.Mux).Routes, 2)
}
