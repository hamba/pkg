package middleware

import (
	"fmt"
	"net/http"

	"github.com/hamba/pkg/log"
)

// Recovery is a middleware that will recover from panics and logs the error.
type Recovery struct {
	handler http.Handler
	l       log.Logger
}

// WithRecovery recovers from panics and log the error.
func WithRecovery(h http.Handler, lable log.Loggable) http.Handler {
	return &Recovery{
		handler: h,
		l:       lable.Logger(),
	}
}

// ServeHTTP serves the request.
func (m Recovery) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if v := recover(); v != nil {
			m.l.Error(fmt.Sprintf("%+v", v))
			w.WriteHeader(500)
		}
	}()

	m.handler.ServeHTTP(w, r)
}
