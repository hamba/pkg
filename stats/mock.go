package stats

// MockStatable implements the Statable interface.
type MockStatable struct {
	s Statter
}

// NewMockStatable creates a new MockLoggable.
func NewMockStatable(s Statter) *MockStatable {
	return &MockStatable{
		s: s,
	}
}

// Statter implements the Statable interface.
func (m *MockStatable) Statter() Statter {
	return m.s
}
