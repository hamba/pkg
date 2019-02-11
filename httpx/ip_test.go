package httpx_test

import (
	"net/http"
	"testing"

	"github.com/hamba/pkg/httpx"
	"github.com/stretchr/testify/assert"
)

func TestRealIP(t *testing.T) {
	tests := []struct {
		req *http.Request
		ip  string
	}{
		{
			&http.Request{RemoteAddr: "127.0.0.1"},
			"127.0.0.1",
		},
		{
			&http.Request{RemoteAddr: "127.0.0.1:8888"},
			"127.0.0.1",
		},
		{
			&http.Request{
				RemoteAddr: "127.0.0.1",
				Header:     http.Header{http.CanonicalHeaderKey("X-Real-Ip"): []string{"1.2.3.4"}},
			},
			"1.2.3.4",
		},
		{
			&http.Request{
				RemoteAddr: "127.0.0.1",
				Header:     http.Header{http.CanonicalHeaderKey("X-Forwarded-For"): []string{"1.2.3.4"}},
			},
			"1.2.3.4",
		},
		{
			&http.Request{
				RemoteAddr: "127.0.0.1",
				Header: http.Header{
					http.CanonicalHeaderKey("X-Forwarded-For"): []string{"1.2.3.4"},
					http.CanonicalHeaderKey("X-Real-Ip"):       []string{"5.6.7.8"}},
			},
			"1.2.3.4",
		},
	}

	for _, tt := range tests {
		got := httpx.RealIP(tt.req)
		assert.Equal(t, tt.ip, got)
	}
}
