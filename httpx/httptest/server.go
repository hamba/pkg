/*
Package httptest provides utilities for HTTP testing.

Deprecated: Use github.com/hamba/testutils instead.

Example Server Usage:

	import (
		"net/http"
		"testing"

		"github.com/hamba/pkg/httpx/httptest"
	)

	func TestSomething(t *testing.T) {
		s := httptest.NewServer(t)
		s.On(http.MethodGet, "/your/endpoint").Times(2).ReturnsString(http.StatusOK, "some return")
		defer s.Close()

		// Call the server

		s.AssertExpectations()
	}
*/
package httptest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ryanuber/go-glob"
)

const (
	// Anything is used where the expectation should not be considered.
	Anything = "httpx/httptest.Anything"
)

// Expectation represents an http request expectation.
type Expectation struct {
	method string
	path   string

	fn http.HandlerFunc

	headers []string
	body    []byte
	status  int

	times int
}

// Times sets the number of times the request can be made.
func (e *Expectation) Times(times int) *Expectation {
	e.times = times

	return e
}

// Header sets the HTTP headers that should be returned.
func (e *Expectation) Header(k, v string) *Expectation {
	e.headers = append(e.headers, k, v)

	return e
}

// Handle sets the HTTP handler function to be run on the request.
func (e *Expectation) Handle(fn http.HandlerFunc) {
	e.fn = fn
}

// ReturnsStatus sets the HTTP stats code to return.
func (e *Expectation) ReturnsStatus(status int) {
	e.body = []byte{}
	e.status = status
}

// Returns sets the HTTP stats and body bytes to return.
func (e *Expectation) Returns(status int, body []byte) {
	e.body = body
	e.status = status
}

// ReturnsString sets the HTTP stats and body string to return.
func (e *Expectation) ReturnsString(status int, body string) {
	e.body = []byte(body)
	e.status = status
}

// Server represents a mock http server.
type Server struct {
	t   *testing.T
	srv *httptest.Server

	expect []*Expectation
}

// NewServer creates a new mock http server.
func NewServer(t *testing.T) *Server {
	srv := &Server{
		t: t,
	}
	srv.srv = httptest.NewServer(http.HandlerFunc(srv.handler))

	return srv
}

// URL returns the url of the mock server.
func (s *Server) URL() string {
	return s.srv.URL
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	for i, exp := range s.expect {
		if exp.method != method && exp.method != Anything {
			continue
		}

		if exp.path != Anything && !glob.Glob(exp.path, path) {
			continue
		}

		for i := 0; i < len(exp.headers); i += 2 {
			w.Header().Add(exp.headers[i], exp.headers[i+1])
		}

		if exp.fn != nil {
			exp.fn(w, r)
		} else {
			w.WriteHeader(exp.status)
			if len(exp.body) > 0 {
				w.Write(exp.body)
			}
		}

		exp.times--
		if exp.times == 0 {
			s.expect = append(s.expect[:i], s.expect[i+1:]...)
		}
		return
	}

	s.t.Errorf("Unexpected call to %s %s", method, path)
}

// On creates an expectation of a request on the server.
func (s *Server) On(method, path string) *Expectation {
	exp := &Expectation{
		method: method,
		path:   path,
		times:  -1,
		status: 200,
	}
	s.expect = append(s.expect, exp)

	return exp
}

// AssertExpectations asserts all expectations have been met.
func (s *Server) AssertExpectations() {
	for _, exp := range s.expect {
		if exp.times > 0 || exp.times == -1 {
			s.t.Errorf("mock: server: Expected a call to %s %s but got none", exp.method, exp.path)
		}
	}
}

// Close closes the server.
func (s *Server) Close() {
	s.srv.Close()
}
