package ostest

type EnvironGenerator struct {
	Environment []string
}

func (e *EnvironGenerator) Environ() []string {
	return e.Environment
}
