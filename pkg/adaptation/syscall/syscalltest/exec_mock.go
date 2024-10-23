package syscalltest

import "fmt"

type ExecMock struct {
	Argv0 string
	Argv  []string
	Envv  []string
	Error error
}

func (n *ExecMock) Exec(argv0 string, argv []string, envv []string) (err error) {
	n.Argv0 = argv0
	n.Argv = argv
	n.Envv = envv
	err = n.Error
	if err == nil {
		// Exec must never return a non-error value
		err = fmt.Errorf("nilexec: exec failed")
	}

	return err
}
