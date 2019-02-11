package httpx

import "github.com/go-zoo/bone"

// NewMux wraps bones New for convenience.
func NewMux() *bone.Mux {
	return bone.New()
}

// CombineMuxes combines the given Muxes into a new Mux.
func CombineMuxes(muxes ...*bone.Mux) *bone.Mux {
	if len(muxes) == 1 {
		return muxes[0]
	}

	res := bone.New()

	for _, mux := range muxes {
		// Register the routes with the new Mux
		for method, routes := range mux.Routes {
			for _, route := range routes {
				res.Register(method, route.Path, route.Handler)
			}
		}

		// Register the validators with the new Mux
		for name, val := range mux.Validators {
			res.RegisterValidator(name, val)
		}
	}

	return res
}
