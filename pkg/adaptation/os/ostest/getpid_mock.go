package ostest

type GetpidMock struct {
	Pid int
}

func (p *GetpidMock) Getpid() int {
	if p.Pid == 0 {
		return 1
	}
	return p.Pid
}
