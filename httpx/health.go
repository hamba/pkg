package httpx

import (
	"net/http"

	"github.com/go-zoo/bone"
)

// DefaultHealthPath is the default HTTP path for checking health.
var DefaultHealthPath = "/health"

// Health represents an object that can check its health.
type Health interface {
	IsHealthy() error
}

// NewHealthMux returns a Mux with a health endpoint.
func NewHealthMux(v ...Health) *bone.Mux {
	mux := NewMux()
	mux.GetFunc(DefaultHealthPath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, h := range v {
			if err := h.IsHealthy(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}))
	return mux
}
