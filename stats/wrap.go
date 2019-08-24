package stats

// Wrapper represents a statter that can return its wrapped statter.
type Wrapper interface {
	Unwrap() Statter
}

// Unwrap returns the result of calling the Unwrap method on stats.
// If stats does not have an Unwrap method, nil will be returned.
func Unwrap(stats Statter) Statter {
	u, ok := stats.(Wrapper)
	if !ok {
		return nil
	}
	return u.Unwrap()
}
