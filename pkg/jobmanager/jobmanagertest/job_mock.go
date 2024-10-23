package jobmanagertest

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/obaraelijah/teleport-challenge/pkg/cgroup/v1"
	"github.com/obaraelijah/teleport-challenge/pkg/io"
	"github.com/obaraelijah/teleport-challenge/pkg/jobmanager"
)

type mockJob struct {
	owner   string
	name    string
	id      uuid.UUID
	running bool
	stdout  io.OutputBuffer
	stderr  io.OutputBuffer
}

func NewMockJob(
	owner string,
	jobName string,
	controllers []cgroup.Controller,
	programPath string,
	arguments ...string,
) jobmanager.Job {
	return &mockJob{
		owner:  owner,
		name:   jobName,
		id:     uuid.New(),
		stdout: io.NewMemoryBuffer(),
		stderr: io.NewMemoryBuffer(),
	}
}

func (m *mockJob) Start() error {

	if m.running {
		return fmt.Errorf("job %s (%v) has already been started", m.name, m.id)
	}

	m.running = true
	_, _ = m.stdout.Write([]byte("this is standard output"))
	_, _ = m.stderr.Write([]byte("this is standard error"))

	return nil
}

func (m *mockJob) Stop() error {
	m.running = false
	m.stdout.Close()
	m.stderr.Close()
	return nil
}

func (m *mockJob) StdoutStream() *io.ByteStream {
	return io.NewByteStream(m.stdout)
}

func (m *mockJob) StderrStream() *io.ByteStream {
	return io.NewByteStream(m.stderr)
}

func (m *mockJob) Status() *jobmanager.JobStatus {
	exitCode := -1

	if m.running {
		exitCode = 0
	}

	return &jobmanager.JobStatus{
		Name:     m.name,
		Id:       m.id.String(),
		Running:  m.running,
		Pid:      1234,
		ExitCode: exitCode,
	}
}

func (m *mockJob) Id() uuid.UUID {
	return m.id
}

func (m *mockJob) Name() string {
	return m.name
}
