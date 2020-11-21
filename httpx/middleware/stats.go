package middleware

import (
	"net/http"
	"strconv"

	"github.com/hamba/pkg/stats"
	"github.com/hamba/timex/mono"
)

// TagsFunc returns a set of tags from a request.
type TagsFunc func(*http.Request) []string

// DefaultTags extracts the method and path from the request.
func DefaultTags(r *http.Request) []string {
	return []string{
		"method", r.Method,
		"path", r.URL.Path,
	}
}

// WithRequestStats collects statistics about the request.
func WithRequestStats(h http.Handler, sable stats.Statable, fns ...TagsFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(fns) == 0 {
			fns = []TagsFunc{DefaultTags}
		}

		var tags []string
		for _, fn := range fns {
			tags = append(tags, fn(r)...)
		}

		stats.Inc(sable, "request.start", 1, 1.0, tags...)

		rw := NewResponseWriter(w)

		start := mono.Now()
		h.ServeHTTP(rw, r)
		dur := mono.Since(start)

		status := strconv.Itoa(rw.Status())
		tags = append(tags, "status", status, "status-group", string(status[0])+"xx")
		stats.Timing(sable, "request.time", dur, 1.0, tags...)
		stats.Inc(sable, "request.complete", 1, 1.0, tags...)
		stats.Inc(sable, "request.size", rw.BytesWritten(), 1.0, tags...)
	})
}
