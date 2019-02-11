package httpx_test

import (
	"net/http"
	"testing"

	"github.com/go-zoo/bone"
	"github.com/hamba/pkg/httpx"
	"github.com/stretchr/testify/assert"
)

func TestNewMux(t *testing.T) {
	m := httpx.NewMux()

	assert.IsType(t, &bone.Mux{}, m)
}

func TestCombineMuxes(t *testing.T) {
	h1 := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	m1 := httpx.NewMux()
	m1.Get("/test", h1)
	v1 := &fakeValidator{}
	m1.RegisterValidator("test1", v1)

	h2 := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	m2 := httpx.NewMux()
	m2.Post("/foobar", h2)
	v2 := &fakeValidator{}
	m2.RegisterValidator("test2", v2)

	got := httpx.CombineMuxes(m1, m2)

	assert.Len(t, got.Routes, 2)
	assert.Len(t, got.Validators, 2)
	assert.Equal(t, map[string]bone.Validator{"test1": v1, "test2": v2}, got.Validators)
}

func TestCombineMuxes_ReturnsMuxIfOnlyOne(t *testing.T) {
	m := httpx.NewMux()

	got := httpx.CombineMuxes(m)

	assert.Equal(t, m, got)
}

type fakeValidator struct{}

func (fakeValidator) Validate(string) bool {
	return true
}
