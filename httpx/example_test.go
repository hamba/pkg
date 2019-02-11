package httpx_test

import (
	"context"

	"github.com/hamba/pkg/httpx"
)

func ExampleCombineMuxes() {
	m1 := httpx.NewMux()
	m2 := httpx.NewMux()

	// Do something with the muxes

	mux := httpx.CombineMuxes(m1, m2)

	_ = httpx.NewServer(":8080", mux).ListenAndServe()
}

func ExampleNewServer() {
	mux := httpx.NewMux()
	srv := httpx.NewServer(":8080", mux)

	if err := srv.ListenAndServe(); err != nil {
		// Handler err
	}

	_ = srv.Shutdown(context.Background()) // Server can be shutdown
}

func ExampleNewHealthMux() {
	mux := httpx.NewHealthMux() // Add your health checker

	_ = httpx.NewServer(":8080", mux).ListenAndServe()
}
