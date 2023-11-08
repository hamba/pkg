// Package healthz provides HTTP healthz handling.
package healthz

import (
	"bytes"
	"fmt"
	"net/http"
)

// HealthChecker represents a named health checker.
type HealthChecker interface {
	Name() string
	Check(*http.Request) error
}

type healthCheck struct {
	name  string
	check func(*http.Request) error
}

// NamedCheck returns a named health check.
func NamedCheck(name string, check func(*http.Request) error) HealthChecker {
	return &healthCheck{
		name:  name,
		check: check,
	}
}

func (c healthCheck) Name() string { return c.name }

func (c healthCheck) Check(req *http.Request) error { return c.check(req) }

// PingHealth returns true when called.
var PingHealth HealthChecker = ping{}

type ping struct{}

func (c ping) Name() string { return "ping" }

func (c ping) Check(_ *http.Request) error { return nil }

// Handler returns an HTTP check handler.
func Handler(name string, errFn func(string), checks ...HealthChecker) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var (
			checkOutput     bytes.Buffer
			failedChecks    []string
			failedLogOutput bytes.Buffer
		)
		for _, check := range checks {
			if err := check.Check(req); err != nil {
				_, _ = fmt.Fprintf(&checkOutput, "- %s failed\n", check.Name())
				_, _ = fmt.Fprintf(&failedLogOutput, "%s failed: %v\n", check.Name(), err)
				failedChecks = append(failedChecks, check.Name())
				continue
			}

			_, _ = fmt.Fprintf(&checkOutput, "+ %s ok\n", check.Name())
		}

		if len(failedChecks) > 0 {
			errFn(failedLogOutput.String())
			http.Error(rw,
				fmt.Sprintf("%s%s check failed", checkOutput.String(), name),
				http.StatusInternalServerError,
			)
			return
		}

		if _, found := req.URL.Query()["verbose"]; !found {
			_, _ = fmt.Fprint(rw, "ok")
			return
		}

		_, _ = checkOutput.WriteTo(rw)
		_, _ = fmt.Fprintf(rw, "%s check passed", name)
	})
}
