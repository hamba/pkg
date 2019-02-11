package httpx

import (
	"net/http"
	"time"

	"github.com/go-zoo/bone"
)

// NewServer creates a new http Server with the given Muxes.
func NewServer(addr string, mux *bone.Mux, muxes ...*bone.Mux) *http.Server {
	if len(muxes) > 0 {
		allMuxes := append([]*bone.Mux{mux}, muxes...)
		mux = CombineMuxes(allMuxes...)
	}

	return &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}
