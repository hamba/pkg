package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hamba/pkg/v2/http/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONInternalServerError(t *testing.T) {
	rec := httptest.NewRecorder()

	render.JSONInternalServerError(rec)

	want := `{"code":500,"error":"internal server error"}`
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, want, rec.Body.String())
}

func TestJSONErrorf(t *testing.T) {
	rec := httptest.NewRecorder()

	render.JSONErrorf(rec, http.StatusBadRequest, "test %s", "message")

	want := `{"code":400,"error":"test message"}`
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, want, rec.Body.String())
}

func TestJSONError(t *testing.T) {
	rec := httptest.NewRecorder()

	render.JSONError(rec, http.StatusBadRequest, "test message")

	want := `{"code":400,"error":"test message"}`
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, want, rec.Body.String())
}

func TestJSON(t *testing.T) {
	rec := httptest.NewRecorder()

	msg := testMessage{
		A: 123,
		B: "test message",
	}

	err := render.JSON(rec, http.StatusBadRequest, msg)
	require.NoError(t, err)

	want := `{"a":123,"b":"test message"}`
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, want, rec.Body.String())
}

type testMessage struct {
	A int    `json:"a"`
	B string `json:"b"`
}
