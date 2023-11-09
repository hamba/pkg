package http

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/hamba/logger/v2"
	"github.com/hamba/pkg/v2/http/healthz"
	"github.com/hamba/statter/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
)

func TestServer(t *testing.T) {
	lnCh := make(chan net.Listener, 1)
	setTestHookServerServe(func(ln net.Listener) {
		lnCh <- ln
	})
	t.Cleanup(func() { setTestHookServerServe(nil) })

	var handlerCalled bool
	h := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})

	srv := NewServer(context.Background(), "localhost:0", h)
	srv.Serve(func(err error) {
		require.NoError(t, err)
	})
	t.Cleanup(func() {
		_ = srv.Close()
	})

	var ln net.Listener
	select {
	case <-time.After(30 * time.Second):
		require.Fail(t, "Timed out waiting for server listener")
	case ln = <-lnCh:
	}

	url := "http://" + ln.Addr().String() + "/"
	statusCode, _ := requireDoRequest(t, url)

	err := srv.Shutdown(time.Second)
	require.NoError(t, err)

	assert.Equal(t, statusCode, http.StatusOK)
	assert.True(t, handlerCalled)
}

func TestServer_WithTLSConfig(t *testing.T) {
	lnCh := make(chan net.Listener, 1)
	setTestHookServerServe(func(ln net.Listener) {
		lnCh <- ln
	})
	t.Cleanup(func() { setTestHookServerServe(nil) })

	var handlerCalled bool
	h := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})

	cert, err := tls.X509KeyPair(localhostCert, localhostKey)
	require.NoError(t, err)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	srv := NewServer(context.Background(), "localhost:0", h, WithTLSConfig(tlsConfig))
	srv.Serve(func(err error) {
		require.NoError(t, err)
	})
	t.Cleanup(func() {
		_ = srv.Close()
	})

	var ln net.Listener
	select {
	case <-time.After(30 * time.Second):
		require.Fail(t, "Timed out waiting for server listener")
	case ln = <-lnCh:
	}

	url := "https://" + ln.Addr().String() + "/"
	statusCode, _ := requireDoRequest(t, url)

	err = srv.Shutdown(time.Second)
	require.NoError(t, err)

	assert.Equal(t, statusCode, http.StatusOK)
	assert.True(t, handlerCalled)
}

func TestServer_WithH2C(t *testing.T) {
	lnCh := make(chan net.Listener, 1)
	setTestHookServerServe(func(ln net.Listener) {
		lnCh <- ln
	})
	t.Cleanup(func() { setTestHookServerServe(nil) })

	var handlerCalled bool
	h := http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		assert.True(t, req.ProtoAtLeast(2, 0))

		handlerCalled = true
	})

	srv := NewServer(context.Background(), "localhost:0", h, WithH2C())
	srv.Serve(func(err error) {
		require.NoError(t, err)
	})
	t.Cleanup(func() {
		_ = srv.Close()
	})

	var ln net.Listener
	select {
	case <-time.After(30 * time.Second):
		require.Fail(t, "Timed out waiting for server listener")
	case ln = <-lnCh:
	}

	c := &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(_ context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
	res, err := c.Get("http://" + ln.Addr().String() + "/")
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	})

	err = srv.Shutdown(time.Second)
	require.NoError(t, err)

	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.True(t, handlerCalled)
}

func TestHealthServer(t *testing.T) {
	stats := statter.New(statter.DiscardReporter, time.Minute)
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Error)

	lnCh := make(chan net.Listener, 1)
	setTestHookServerServe(func(ln net.Listener) {
		lnCh <- ln
	})
	t.Cleanup(func() { setTestHookServerServe(nil) })

	check := healthz.NamedCheck("test", func(*http.Request) error {
		return nil
	})

	srv := NewHealthServer(context.Background(), HealthServerConfig{
		Addr:    "localhost:0",
		Handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
		Stats:   stats,
		Log:     log,
	})
	err := srv.AddHealthzChecks(check)
	require.NoError(t, err)
	srv.Serve(func(err error) {
		require.NoError(t, err)
	})
	t.Cleanup(func() {
		_ = srv.Close()
	})

	var ln net.Listener
	select {
	case <-time.After(30 * time.Second):
		require.Fail(t, "Timed out waiting for server listener")
	case ln = <-lnCh:
	}

	url := "http://" + ln.Addr().String() + "/readyz?verbose=1"
	statusCode, body := requireDoRequest(t, url)
	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, "+ test ok\n+ shutdown ok\nreadyz check passed", body)

	url = "http://" + ln.Addr().String() + "/livez?verbose=1"
	statusCode, body = requireDoRequest(t, url)
	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, "+ test ok\nlivez check passed", body)
}

func TestHealthServer_ShutdownCausesReadyzCheckToFail(t *testing.T) {
	stats := statter.New(statter.DiscardReporter, time.Minute)
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Error)

	lnCh := make(chan net.Listener, 1)
	setTestHookServerServe(func(ln net.Listener) {
		lnCh <- ln
	})
	t.Cleanup(func() { setTestHookServerServe(nil) })

	calledCh := make(chan struct{})
	check := healthz.NamedCheck("test", func(*http.Request) error {
		close(calledCh)
		time.Sleep(time.Millisecond)
		return nil
	})

	srv := NewHealthServer(context.Background(), HealthServerConfig{
		Addr:    "localhost:0",
		Handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
		Stats:   stats,
		Log:     log,
	})
	err := srv.AddHealthzChecks(check)
	require.NoError(t, err)
	srv.Serve(func(err error) {
		require.NoError(t, err)
	})
	t.Cleanup(func() {
		_ = srv.Close()
	})

	var ln net.Listener
	select {
	case <-time.After(30 * time.Second):
		require.Fail(t, "Timed out waiting for server listener")
	case ln = <-lnCh:
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		<-calledCh

		err = srv.Shutdown(time.Second)
		require.NoError(t, err)
	}()

	url := "http://" + ln.Addr().String() + "/readyz?verbose=1"
	statusCode, body := requireDoRequest(t, url)
	assert.Equal(t, statusCode, http.StatusInternalServerError)
	assert.Equal(t, "+ test ok\n- shutdown failed\nreadyz check failed\n", body)

	wg.Wait()
}

func requireDoRequest(t *testing.T, path string) (int, string) {
	t.Helper()

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, string(b)
}

func setTestHookServerServe(fn func(net.Listener)) {
	testHookServerServe = fn
}

var (
	localhostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDOTCCAiGgAwIBAgIQSRJrEpBGFc7tNb1fb5pKFzANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMCAXDTcwMDEwMTAwMDAwMFoYDzIwODQwMTI5MTYw
MDAwWjASMRAwDgYDVQQKEwdBY21lIENvMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEA6Gba5tHV1dAKouAaXO3/ebDUU4rvwCUg/CNaJ2PT5xLD4N1Vcb8r
bFSW2HXKq+MPfVdwIKR/1DczEoAGf/JWQTW7EgzlXrCd3rlajEX2D73faWJekD0U
aUgz5vtrTXZ90BQL7WvRICd7FlEZ6FPOcPlumiyNmzUqtwGhO+9ad1W5BqJaRI6P
YfouNkwR6Na4TzSj5BrqUfP0FwDizKSJ0XXmh8g8G9mtwxOSN3Ru1QFc61Xyeluk
POGKBV/q6RBNklTNe0gI8usUMlYyoC7ytppNMW7X2vodAelSu25jgx2anj9fDVZu
h7AXF5+4nJS4AAt0n1lNY7nGSsdZas8PbQIDAQABo4GIMIGFMA4GA1UdDwEB/wQE
AwICpDATBgNVHSUEDDAKBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MB0GA1Ud
DgQWBBStsdjh3/JCXXYlQryOrL4Sh7BW5TAuBgNVHREEJzAlggtleGFtcGxlLmNv
bYcEfwAAAYcQAAAAAAAAAAAAAAAAAAAAATANBgkqhkiG9w0BAQsFAAOCAQEAxWGI
5NhpF3nwwy/4yB4i/CwwSpLrWUa70NyhvprUBC50PxiXav1TeDzwzLx/o5HyNwsv
cxv3HdkLW59i/0SlJSrNnWdfZ19oTcS+6PtLoVyISgtyN6DpkKpdG1cOkW3Cy2P2
+tK/tKHRP1Y/Ra0RiDpOAmqn0gCOFGz8+lqDIor/T7MTpibL3IxqWfPrvfVRHL3B
grw/ZQTTIVjjh4JBSW3WyWgNo/ikC1lrVxzl4iPUGptxT36Cr7Zk2Bsg0XqwbOvK
5d+NTDREkSnUbie4GeutujmX3Dsx88UiV6UY/4lHJa6I5leHUNOHahRbpbWeOfs/
WkBKOclmOV2xlTVuPw==
-----END CERTIFICATE-----`)

	localhostKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDoZtrm0dXV0Aqi
4Bpc7f95sNRTiu/AJSD8I1onY9PnEsPg3VVxvytsVJbYdcqr4w99V3AgpH/UNzMS
gAZ/8lZBNbsSDOVesJ3euVqMRfYPvd9pYl6QPRRpSDPm+2tNdn3QFAvta9EgJ3sW
URnoU85w+W6aLI2bNSq3AaE771p3VbkGolpEjo9h+i42TBHo1rhPNKPkGupR8/QX
AOLMpInRdeaHyDwb2a3DE5I3dG7VAVzrVfJ6W6Q84YoFX+rpEE2SVM17SAjy6xQy
VjKgLvK2mk0xbtfa+h0B6VK7bmODHZqeP18NVm6HsBcXn7iclLgAC3SfWU1jucZK
x1lqzw9tAgMBAAECggEABWzxS1Y2wckblnXY57Z+sl6YdmLV+gxj2r8Qib7g4ZIk
lIlWR1OJNfw7kU4eryib4fc6nOh6O4AWZyYqAK6tqNQSS/eVG0LQTLTTEldHyVJL
dvBe+MsUQOj4nTndZW+QvFzbcm2D8lY5n2nBSxU5ypVoKZ1EqQzytFcLZpTN7d89
EPj0qDyrV4NZlWAwL1AygCwnlwhMQjXEalVF1ylXwU3QzyZ/6MgvF6d3SSUlh+sq
XefuyigXw484cQQgbzopv6niMOmGP3of+yV4JQqUSb3IDmmT68XjGd2Dkxl4iPki
6ZwXf3CCi+c+i/zVEcufgZ3SLf8D99kUGE7v7fZ6AQKBgQD1ZX3RAla9hIhxCf+O
3D+I1j2LMrdjAh0ZKKqwMR4JnHX3mjQI6LwqIctPWTU8wYFECSh9klEclSdCa64s
uI/GNpcqPXejd0cAAdqHEEeG5sHMDt0oFSurL4lyud0GtZvwlzLuwEweuDtvT9cJ
Wfvl86uyO36IW8JdvUprYDctrQKBgQDycZ697qutBieZlGkHpnYWUAeImVA878sJ
w44NuXHvMxBPz+lbJGAg8Cn8fcxNAPqHIraK+kx3po8cZGQywKHUWsxi23ozHoxo
+bGqeQb9U661TnfdDspIXia+xilZt3mm5BPzOUuRqlh4Y9SOBpSWRmEhyw76w4ZP
OPxjWYAgwQKBgA/FehSYxeJgRjSdo+MWnK66tjHgDJE8bYpUZsP0JC4R9DL5oiaA
brd2fI6Y+SbyeNBallObt8LSgzdtnEAbjIH8uDJqyOmknNePRvAvR6mP4xyuR+Bv
m+Lgp0DMWTw5J9CKpydZDItc49T/mJ5tPhdFVd+am0NAQnmr1MCZ6nHxAoGABS3Y
LkaC9FdFUUqSU8+Chkd/YbOkuyiENdkvl6t2e52jo5DVc1T7mLiIrRQi4SI8N9bN
/3oJWCT+uaSLX2ouCtNFunblzWHBrhxnZzTeqVq4SLc8aESAnbslKL4i8/+vYZlN
s8xtiNcSvL+lMsOBORSXzpj/4Ot8WwTkn1qyGgECgYBKNTypzAHeLE6yVadFp3nQ
Ckq9yzvP/ib05rvgbvrne00YeOxqJ9gtTrzgh7koqJyX1L4NwdkEza4ilDWpucn0
xiUZS4SoaJq6ZvcBYS62Yr1t8n09iG47YL8ibgtmH3L+svaotvpVxVK+d7BLevA/
ZboOWVe3icTy64BT3OQhmg==
-----END RSA PRIVATE KEY-----`)
)
