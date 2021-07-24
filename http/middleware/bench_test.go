package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hamba/pkg/v2/http/middleware"
	"github.com/hamba/statter/v2"
)

func BenchmarkWithStats(b *testing.B) {
	s := statter.New(statter.DiscardReporter, time.Second)
	h := middleware.WithStats("test", s, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {}),
	)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.ServeHTTP(resp, req)
		}
	})
}
