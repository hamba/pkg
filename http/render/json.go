// Package render provides HTTP output rendering helper functions.
package render

import (
	"fmt"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// JSONContentType represents MIME type for JSON content.
const JSONContentType = "application/json"

// APIError contains error information that is rendered by JSONError.
type APIError struct {
	// Code is the http status code.
	Code int `json:"code"`

	// Error is the reason for the error.
	Error string `json:"error"`
}

// JSONInternalServerError writes a JSON internal server error.
func JSONInternalServerError(rw http.ResponseWriter) {
	JSONError(rw, http.StatusInternalServerError, strings.ToLower(http.StatusText(http.StatusInternalServerError)))
}

// JSONErrorf writes a JSON error message.
func JSONErrorf(rw http.ResponseWriter, code int, format string, args ...any) {
	JSONError(rw, code, fmt.Sprintf(format, args...))
}

// JSONError writes a JSON error message.
func JSONError(rw http.ResponseWriter, code int, reason string) {
	rw.Header().Set("Content-Type", JSONContentType)
	rw.WriteHeader(code)

	apiErr := APIError{Code: code, Error: reason}
	b, err := jsoniter.Marshal(apiErr)
	if err != nil {
		_, _ = rw.Write([]byte(`{"code":500,"error":"internal server error"}`))
	}

	_, _ = rw.Write(b)
}

// JSON writes a JSON response.
func JSON(rw http.ResponseWriter, code int, v any) error {
	b, err := jsoniter.Marshal(v)
	if err != nil {
		return err
	}

	rw.Header().Set("Content-Type", JSONContentType)
	rw.WriteHeader(code)
	_, _ = rw.Write(b)
	return nil
}
