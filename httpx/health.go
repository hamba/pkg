package httpx

import (
	"net/http"
)

// DefaultHealthPath is the default HTTP path for checking health.
var DefaultHealthPath = "/health"

// Health represents an object that can check its health.
type Health interface {
	IsHealthy() error
}

// NewHealthHandler returns a handler for application health checking.
func NewHealthHandler(v ...Health) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		for _, h := range v {
			if err := h.IsHealthy(); err != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		rw.WriteHeader(http.StatusOK)
	}
}
