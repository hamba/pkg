package http

import "net/http"

// OK replies to the request with an HTTP 200 ok reply.
//
// Deprecated: Use healthz instead.
func OK(rw http.ResponseWriter, _ *http.Request) { rw.WriteHeader(http.StatusOK) }

// OKHandler returns a simple request handler
// that replies to each request with a “200 OK” reply.
//
// Deprecated: Use healthz instead.
func OKHandler() http.Handler { return http.HandlerFunc(OK) }

// DefaultHealthPath is the default HTTP path for checking health.
//
// Deprecated: Use healthz instead.
var DefaultHealthPath = "/health"

// Health represents an object that can check its health.
//
// Deprecated: Use healthz instead.
type Health interface {
	IsHealthy() error
}

// NewHealthHandler returns a handler for application health checking.
//
// Deprecated: Use healthz instead.
func NewHealthHandler(v ...Health) http.HandlerFunc {
	return func(rw http.ResponseWriter, _ *http.Request) {
		for _, h := range v {
			if err := h.IsHealthy(); err != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		rw.WriteHeader(http.StatusOK)
	}
}
