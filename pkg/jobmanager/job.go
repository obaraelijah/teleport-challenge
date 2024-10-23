package jobmanager

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"github.com/obaraelijah/teleport-challenge/pkg/cgroup/v1"
	"github.com/obaraelijah/teleport-challenge/pkg/io"
)

type JobStatus struct {
	Name      string
	Id        string
	Running   bool
	Pid       int
	ExitCode  int
	SignalNum syscall.Signal
	RunError  error
}

type Job interface {
	Start() error
	Stop() error
	StdoutStream() *io.ByteStream
	StderrStream() *io.ByteStream
	Status() *JobStatus
	Id() uuid.UUID
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
) Job {

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
) Job {

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

	if j.cmd != nil {
		return fmt.Errorf("job %s (%v) has already been started", j.name, j.id)
	}

	cgroupSet := cgroup.NewSet(j.id, j.cgControllers...)
	if err := cgroupSet.Create(); err != nil {
		return err
	}

	args := cgroupSet.TaskFiles()
	args = append(args, "--")
	args = append(args, j.programName)
	args = append(args, j.programArgs...)

	//j.cmd = exec.Command("build/cgexec", args...)
	j.cmd = exec.Command("/tmp/cgexec", args...)
	j.cmd.Stdout = j.stdoutBuffer
	j.cmd.Stderr = j.stderrBuffer
	j.cmd.Env = make([]string, 0) // Do not pass along our environment
	//j.cmd.Dir = "/" // If we were to chroot
	j.cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: "", // This would be non-empty to actually do a chroot
		Cloneflags: syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET,
	}

	j.running = true

	go func() {
		defer j.lockedOperation(func() {
			if err := cgroupSet.Destroy(); err != nil {
				j.runErrors = append(j.runErrors, err)
			}
			j.running = false
		})

		// Run blocks until the newly-created process terminates.  It calls
		// Wait internally
		err := j.cmd.Run()

		// Once Wait returns, all output has been written to Stdout and Stderr
		j.lockedOperation(func() {
			if err != nil {
				j.runErrors = append(j.runErrors, err)
			}

			if err := j.stdoutBuffer.Close(); err != nil {
				j.runErrors = append(j.runErrors, err)
			}

			if err := j.stderrBuffer.Close(); err != nil {
				j.runErrors = append(j.runErrors, err)
			}
		})
	}()

	return nil
}

func (j *job) Stop() error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if !j.running {
		// If the job isn't running, it is stopped already
		return nil
	}

	if err := j.cmd.Process.Kill(); err != nil && err != os.ErrProcessDone {
		return err
	}

	return nil
}

func (j *job) StdoutStream() *io.ByteStream {
	// Unlocked read of j.stdoutBuffer should be ok since it's not modified once created
	return io.NewByteStream(j.stdoutBuffer)
}

func (j *job) StderrStream() *io.ByteStream {
	// Unlocked read of j.stderrBuffer should be ok since it's not modified once created
	return io.NewByteStream(j.stderrBuffer)
}

func (j *job) Status() *JobStatus {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	status := &JobStatus{
		Name:      j.name,
		Id:        j.id.String(),
		Running:   j.running,
		Pid:       -1,
		SignalNum: syscall.Signal(-1),
	}

	if j.runErrors != nil {
		status.RunError = fmt.Errorf("run error: %v", j.runErrors)
	}

	if state := j.cmd.ProcessState; state == nil {
		if j.cmd.Process != nil {
			status.Pid = j.cmd.Process.Pid
		}
	} else {
		if sys := state.Sys(); sys != nil {
			status.SignalNum = sys.(syscall.WaitStatus).Signal()
		}
	}

	return status
}

func (j *job) Id() uuid.UUID {
	return j.id
}

// lockedOperation is a simple runs the functor with the job lock held.
// The caller must not hold the lock.
func (j *job) lockedOperation(fn func()) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	fn()
}
