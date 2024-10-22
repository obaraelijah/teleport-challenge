package ostest

type PidGenerator struct {
	Pid int
}

func (p *PidGenerator) Getpid() int {
	if p.Pid == 0 {
		return 1
	}
	return p.Pid
}
