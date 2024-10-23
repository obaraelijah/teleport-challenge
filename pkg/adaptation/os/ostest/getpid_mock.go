package ostest

// GetpidMock is a component that provides a mock implementation of the
// os.Getpid() function.  The function returns the configured Pid.
type GetpidMock struct {
	Pid int
}

func (p *GetpidMock) Getpid() int {
	if p.Pid == 0 {
		return 1
	}
	return p.Pid
}
