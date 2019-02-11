package middleware

import (
	"net/http"
)

// ResponseWriter is a wrapper around http.ResponseWriter that provides extra
// information about the response.
type ResponseWriter interface {
	http.ResponseWriter

	// Status returns the status code of the response or 0 if the response has
	// not be written.
	Status() int
}

// NewResponseWriter create a new ResponseWriter.
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{rw, 0}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

// Write writes the data to the connection as part of an HTTP reply.
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}

	return rw.ResponseWriter.Write(b)
}

// WriteHeader sends an HTTP response header with status code.
func (rw *responseWriter) WriteHeader(s int) {
	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

// Status returns the status code of the response or 0 if the response has
// not be written.
func (rw *responseWriter) Status() int {
	return rw.status
}
