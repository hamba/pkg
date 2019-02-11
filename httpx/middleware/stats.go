package middleware

import (
	"net/http"

	"github.com/hamba/pkg/stats"
)

// TagsFunc returns a set of tags from a request
type TagsFunc func(*http.Request) []interface{}

// DefaultTags extracts the method and path from the request.
func DefaultTags(r *http.Request) []interface{} {
	return []interface{}{
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

		var tags []interface{}
		for _, fn := range fns {
			tags = append(tags, fn(r)...)
		}

		rw := NewResponseWriter(w)

		stats.Inc(sable, "request.start", 1, 1.0, tags...)

		t := stats.Time(sable, "request.time", 1.0, tags...)

		h.ServeHTTP(rw, r)

		t.Done()

		tags = append(tags, "status", rw.Status())
		stats.Inc(sable, "request.complete", 1, 1.0, tags...)
	})
}
