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

	// BytesWritten returns the number of bytes written to the writer.
	BytesWritten() int64
}

// NewResponseWriter create a new ResponseWriter.
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{ResponseWriter: rw, status: http.StatusOK}
}

type responseWriter struct {
	http.ResponseWriter
	status int
	bytes  int64
}

// Write writes the data to the connection as part of an HTTP reply.
func (rw *responseWriter) Write(p []byte) (int, error) {
	rw.bytes += int64(len(p))
	return rw.ResponseWriter.Write(p)
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

// BytesWritten returns the number of bytes written to the writer.
func (rw *responseWriter) BytesWritten() int64 {
	return rw.bytes
}
