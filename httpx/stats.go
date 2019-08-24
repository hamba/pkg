package httpx

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/hamba/pkg/stats"
)

// DefaultStatsPath is the default HTTP path for collecting metrics.
var DefaultStatsPath = "/metrics"

// Stats represents a statter that can expose stats for collection.
type StatsHandler interface {
	Handler() http.Handler
}

// NewStatsMux returns a Mux with a stats endpoint. If no statter
// implements the StatsHandler interface, and empty Mux is returned.
func NewStatsMux(s stats.Statter) *bone.Mux {
	mux := NewMux()

	var h StatsHandler
	for s != nil {
		if sh, ok := s.(StatsHandler); ok {
			h = sh
			break
		}

		s = stats.Unwrap(s)
	}

	if h != nil {
		mux.Get(DefaultStatsPath, h.Handler())
	}

	return mux
}
