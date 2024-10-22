package jobmanager

import (
	"os/exec"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"github.com/obaraelijah/teleport-challenge/pkg/cgroup/v1"
	"github.com/obaraelijah/teleport-challenge/pkg/io"
)

type JobStatus struct {
	Running   bool
	Pid       int
	ExitCode  int
	SignalNum syscall.Signal
	RunError  error
}

type job struct {
	mutex         sync.Mutex
	id            uuid.UUID
	name          string
	cgControllers []cgroup.Controller
	programName   string
	programArgs   []string
	cmd           *exec.Cmd
	stdoutBuffer  io.OutputBuffer
	stderrBuffer  io.OutputBuffer
	running       bool
	runErrors     []error
}

func NewJob(
	name string,
	cgControllers []cgroup.Controller,
	programName string,
	programArgs ...string,
) *job {
	return NewJobDetailed(
		name,
		cgControllers,
		io.NewMemoryBuffer(),
		io.NewMemoryBuffer(),
		programName,
		programArgs...,
	)
}

func NewJobDetailed(
	name string,
	cgControllers []cgroup.Controller,
	stdoutBuffer io.OutputBuffer,
	stderrBuffer io.OutputBuffer,
	programName string,
	programArgs ...string,
) *job {
	return &job{
		id:            uuid.New(),
		name:          name,
		cgControllers: cgControllers,
		programName:   programName,
		programArgs:   programArgs,
		stdoutBuffer:  stdoutBuffer,
		stderrBuffer:  stderrBuffer,
	}
}

func (j *job) Start() error {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	return nil
}
