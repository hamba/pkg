package healthz_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hamba/pkg/v2/http/healthz"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	goodCheck := healthz.NamedCheck("good", func(*http.Request) error { return nil })

	var gotOutput string
	h := healthz.Handler("readyz", func(output string) {
		gotOutput = output
	}, goodCheck)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, gotOutput)
	assert.Equal(t, `ok`, rec.Body.String())
}

func TestHandler_Verbose(t *testing.T) {
	goodCheck := healthz.NamedCheck("good", func(*http.Request) error { return nil })

	var gotOutput string
	h := healthz.Handler("readyz", func(output string) {
		gotOutput = output
	}, goodCheck)

	req := httptest.NewRequest(http.MethodGet, "/readyz?verbose=1", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, gotOutput)
	assert.Equal(t, "+ good ok\nreadyz check passed", rec.Body.String())
}

func TestHandler_WithFailingChecks(t *testing.T) {
	goodCheck := healthz.NamedCheck("good", func(*http.Request) error { return nil })
	badCheck := healthz.NamedCheck("bad", func(*http.Request) error { return errors.New("test error") })

	var gotOutput string
	h := healthz.Handler("readyz", func(output string) {
		gotOutput = output
	}, goodCheck, badCheck)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "bad failed: test error\n", gotOutput)
	assert.Equal(t, "+ good ok\n- bad failed\nreadyz check failed\n", rec.Body.String())
}
