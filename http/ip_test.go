package http_test

import (
	"net/http"
	"testing"

	httpx "github.com/hamba/pkg/v2/http"
	"github.com/stretchr/testify/assert"
)

func TestRealIP(t *testing.T) {
	tests := []struct {
		name string
		req  *http.Request
		want string
	}{
		{
			name: "remote addr without port",
			req:  &http.Request{RemoteAddr: "127.0.0.1"},
			want: "127.0.0.1",
		},
		{
			name: "remote addr with port",
			req:  &http.Request{RemoteAddr: "127.0.0.1:8888"},
			want: "127.0.0.1",
		},
		{
			name: "real-ip",
			req: &http.Request{
				RemoteAddr: "127.0.0.1",
				Header:     http.Header{http.CanonicalHeaderKey("X-Real-Ip"): []string{"1.2.3.4"}},
			},
			want: "1.2.3.4",
		},
		{
			name: "forwarded for",
			req: &http.Request{
				RemoteAddr: "127.0.0.1",
				Header:     http.Header{http.CanonicalHeaderKey("X-Forwarded-For"): []string{"1.2.3.4", "127.0.0.1"}},
			},
			want: "1.2.3.4",
		},
		{
			name: "forwarded for over real ip",
			req: &http.Request{
				RemoteAddr: "127.0.0.1",
				Header: http.Header{
					http.CanonicalHeaderKey("X-Forwarded-For"): []string{"1.2.3.4", "11.0.0.1"},
					http.CanonicalHeaderKey("X-Real-Ip"):       []string{"5.6.7.8"}},
			},
			want: "1.2.3.4",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := httpx.RealIP(test.req)

			assert.Equal(t, test.want, got)
		})
	}
}
