package jobmanager

import (
	"os/exec"
	"sync"

	"github.com/google/uuid"
	"github.com/obaraelijah/teleport-challenge/pkg/cgroup/v1"
	"github.com/obaraelijah/teleport-challenge/pkg/io"
)

type job struct {
	mutex        sync.Mutex
	id           uuid.UUID
	name         string
	cgroupSet    *cgroup.Set
	cmd          exec.Cmd
	stdoutBuffer io.OutputBuffer
	stderrBuffer io.OutputBuffer
}

func NewJob(name string, cgroupSet *cgroup.Set, command string, args ...string) *job {
	return NewJobDetailed(name, cgroupSet, io.NewMemoryBuffer(), io.NewMemoryBuffer(), command, args...)
}

func NewJobDetailed(
	name string,
	cgroupSet *cgroup.Set,
	stdoutBuffer io.OutputBuffer,
	stderrBuffer io.OutputBuffer,
	command string,
	args ...string,
) *job {
	return &job{
		id:           uuid.New(),
		name:         name,
		cgroupSet:    cgroupSet,
		cmd:          *exec.Command(""),
		stdoutBuffer: stdoutBuffer,
		stderrBuffer: stderrBuffer,
	}
}
