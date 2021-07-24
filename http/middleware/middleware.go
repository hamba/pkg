// Package middleware implements reusable HTTP middleware.
package middleware

import (
	"net/http"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
	"github.com/hamba/statter/v2"
	"github.com/hamba/statter/v2/tags"
	"github.com/hamba/timex/mono"
)

// WithRecovery recovers from panics and log the error.
func WithRecovery(h http.Handler, log *logger.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				log.Error("Panic while serving request",
					ctx.Interface("err", v),
					ctx.Str("method", req.Method),
					ctx.Str("url", req.URL.String()),
					ctx.Stack("stack"),
				)

				rw.WriteHeader(http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(rw, req)
	})
}

// WithStats collects statistics about HTTP requests.
func WithStats(name string, s *statter.Statter, h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t := make([]statter.Tag, 1, 2)
		if name == "" {
			name = req.URL.Path
		}
		t[0] = tags.Str("handler", name)

		s.Counter("requests", t...).Inc(1)

		wrap := newResponseWrapper(rw)

		start := mono.Now()
		h.ServeHTTP(wrap, req)
		dur := mono.Since(start)

		t = append(t, tags.StatusCode("code", wrap.Status()))
		s.Counter("responses", t...).Inc(1)
		s.Histogram("response.size", t...).Observe(float64(wrap.BytesWritten()))
		s.Timing("response.duration", t...).Observe(dur)
	})
}

type responseWrapper struct {
	http.ResponseWriter

	status int
	bytes  int64
}

func newResponseWrapper(rw http.ResponseWriter) *responseWrapper {
	return &responseWrapper{
		ResponseWriter: rw,
		status:         http.StatusOK,
	}
}

// Write writes the data to the connection as part of an HTTP reply.
func (rw *responseWrapper) Write(p []byte) (int, error) {
	rw.bytes += int64(len(p))
	return rw.ResponseWriter.Write(p)
}

// WriteHeader sends an HTTP response header with status code.
func (rw *responseWrapper) WriteHeader(s int) {
	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

// Status returns the status code of the response or 0 if the response has
// not be written.
func (rw *responseWrapper) Status() int {
	return rw.status
}

// BytesWritten returns the number of bytes written to the writer.
func (rw *responseWrapper) BytesWritten() int64 {
	return rw.bytes
}
