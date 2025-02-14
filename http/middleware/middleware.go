// Package middleware provides reusable HTTP middleware.
package middleware

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
	"github.com/hamba/pkg/v2/http/request"
	"github.com/hamba/statter/v2"
	"github.com/hamba/statter/v2/reporter/prometheus"
	"github.com/hamba/statter/v2/tags"
	"github.com/segmentio/ksuid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

// Recovery is a wrapper for WithRecovery.
func Recovery(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return WithRecovery(next, log)
	}
}

// WithRequestID sets the request id on request context and in the response.
func WithRequestID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		id := ksuid.New().String()

		rw.Header().Set("X-Request-ID", id)
		req = req.WithContext(request.WithID(req.Context(), id))

		h.ServeHTTP(rw, req)
	})
}

// RequestID is a wrapper for WithRequestID.
func RequestID() func(http.Handler) http.Handler {
	return WithRequestID
}

// WithStats collects statistics about HTTP requests.
func WithStats(name string, s *statter.Statter, h http.Handler) http.Handler {
	prometheus.RegisterHistogram(s,
		"response.size",
		[]string{"handler", "code", "code-group"},
		[]float64{200, 500, 900, 1500, 5000, 10000},
		"The size of a response in bytes",
	)

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t := make([]statter.Tag, 1, 3)
		if name == "" {
			name = req.URL.Path
		}
		t[0] = tags.Str("handler", name)

		s.Counter("requests", t...).Inc(1)

		wrap := newResponseWrapper(rw)

		start := time.Now()
		h.ServeHTTP(wrap, req)
		dur := time.Since(start)

		t = append(t, tags.StatusCode("code-group", wrap.Status()))
		t = append(t, tags.Int("code", wrap.Status()))
		s.Counter("responses", t...).Inc(1)
		s.Histogram("response.size", t...).Observe(float64(wrap.BytesWritten()))
		s.Timing("response.duration", t...).Observe(dur)
	})
}

// Stats is a wrapper for WithStats.
func Stats(name string, s *statter.Statter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return WithStats(name, s, next)
	}
}

// Tracing collects traces on HTTP requests.
func Tracing(op string, opts ...otelhttp.Option) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, op, opts...)
	}
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

// Hijack returns a hijacked connection or an error.
//
// This is required by some websocket libraries.
func (rw *responseWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijacker not supported")
	}
	return h.Hijack()
}

// Unwrap returns the underlying response writer.
// This is used by http.ResponseController to find the first
// response writer that implements an interface.
func (rw *responseWrapper) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}
