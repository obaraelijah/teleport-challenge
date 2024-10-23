package cgrouptest

// ControllerMock provides a mock implementation of the v1.Controller interface
// for use in unit test.  This implementation does not modify any actual
// cgroup.
type ControllerMock struct {
	ControllerName   string
	ApplyReturnValue error
}

func (d *ControllerMock) Name() string {
	return d.ControllerName
}

func (d *ControllerMock) Apply(string) error {
	return d.ApplyReturnValue
}
