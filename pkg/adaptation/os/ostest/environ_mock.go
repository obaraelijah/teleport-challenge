package ostest

type EnvironMock struct {
	Environment []string
}

func (e *EnvironMock) Environ() []string {
	return e.Environment
}
