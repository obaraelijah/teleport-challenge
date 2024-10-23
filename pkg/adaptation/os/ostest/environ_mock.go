package ostest

// EnvironMock is a component that provides a mock implementation of the
// os.Environ() function.  The function returns the configured Environment.
type EnvironMock struct {
	Environment []string
}

func (e *EnvironMock) Environ() []string {
	return e.Environment
}
