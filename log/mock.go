package log

// MockLoggable implements the Loggable interface.
type MockLoggable struct {
	l Logger
}

// NewMockLoggable creates a new MockLoggable.
func NewMockLoggable(l Logger) *MockLoggable {
	return &MockLoggable{
		l: l,
	}
}

// Logger implements the Loggable interface.
func (m *MockLoggable) Logger() Logger {
	return m.l
}
