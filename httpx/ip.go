package httpx

import (
	"net/http"
	"strings"
)

// RealIP resolves the real client IP address from the request.
func RealIP(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")

	if len(hdrForwardedFor) == 0 && len(hdrRealIP) == 0 {
		idx := strings.LastIndex(r.RemoteAddr, ":")
		if idx == -1 {
			return r.RemoteAddr
		}
		return r.RemoteAddr[:idx]
	}

	if len(hdrForwardedFor) > 0 {
		pieces := strings.Split(hdrForwardedFor, ",")
		return pieces[0]
	}

	return hdrRealIP
}
