package jobmanager

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/obaraelijah/teleport-challenge/pkg/cgroup/cgroupv1"
	"github.com/obaraelijah/teleport-challenge/pkg/config"
	"github.com/obaraelijah/teleport-challenge/pkg/io"
)

const (
	// Superuser is the name of the user who can access any job.
	Superuser = "administrator"
)

// Job defines an interface for objects that behave like jobs.  This enables
// us to define both a production job type as well a a job type for unit
// testing.
type Job interface {
	Start() error
	Stop() error
	Status() *JobStatus
	StdoutStream() *io.ByteStream
	StderrStream() *io.ByteStream
	Name() string
	ID() uuid.UUID
}

// JobConstructor is a type that models a function for creating Jobs.
// This enables us to have a "real" job constructor function as well as
// constructor functions for mock job implementations that share the same
// signature.
type JobConstructor func(
	owner string,
	jobName string,
	controllers []cgroupv1.Controller,
	programPath string,
	arguments ...string,
) Job

// Manager maintains the set of jobs and enforces the authorization policy
type Manager struct {
	mutex               sync.RWMutex
	jobsByUserByJobID   map[string]map[string]Job // userID->jobID->job
	jobsByUserByJobName map[string]map[string]Job // userID->jobName->job
	allJobsByJobID      map[string]Job            // jobID->job
	controllers         []cgroupv1.Controller
	jobConstructor      JobConstructor
}

// NewManager creates and returns a new standard Manager.
func NewManager() *Manager {
	controllers := []cgroupv1.Controller{
		cgroupv1.NewCpuController().SetCpus(config.CgroupDefaultCpuLimit),
		cgroupv1.NewMemoryController().SetLimit(config.CgroupDefaultMemoryLimit),
		cgroupv1.NewBlockIOController().
			SetReadBpsDevice(config.CgroupDefaultBlkioReadLimit).
			SetWriteBpsDevice(config.CgroupDefaultBlkioWriteLimit),
	}

	return NewManagerDetailed(NewJob, controllers)

}

// NewManagerDetailed returns a new Manger with the given values.
// The jobConstructor is a function for creating new jobs.  In production
// this will point to NewJob.  For unit tests, this might point to a
// constructor function for a mock type.
// The given controllers is the list of cgroup controllers to manage while
// running jobs.
func NewManagerDetailed(jobConstructor JobConstructor, controllers []cgroupv1.Controller) *Manager {
	return &Manager{
		jobsByUserByJobID:   make(map[string]map[string]Job),
		jobsByUserByJobName: make(map[string]map[string]Job),
		allJobsByJobID:      make(map[string]Job),
		controllers:         controllers,
		jobConstructor:      jobConstructor,
	}
}

// Start starts a new job with the given JobName for the given userID.
// The programPath and arguments are the program the user wants to run and
// the arguments to that program.
func (m *Manager) Start(userID, jobName, programPath string, arguments []string) (Job, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.jobsByUserByJobID[userID]; !exists {
		m.jobsByUserByJobID[userID] = make(map[string]Job)
		m.jobsByUserByJobName[userID] = make(map[string]Job)
	}

	if _, exists := m.jobsByUserByJobName[userID][jobName]; exists {
		return nil, fmt.Errorf("job with name '%s' exists already", jobName)
	}

	job := m.jobConstructor(userID, jobName, m.controllers, programPath, arguments...)

	m.jobsByUserByJobID[userID][job.ID().String()] = job
	m.jobsByUserByJobName[userID][jobName] = job
	m.allJobsByJobID[job.ID().String()] = job

	return job, job.Start()
}

// Stop stops an existing job with the given jobID for the given userID.
func (m *Manager) Stop(userID, jobID string) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	job, err := m.findJobByUser(userID, jobID)
	if err != nil {
		return err
	}
	return job.Stop()
}

// List returns a list of the jobs owned by the given userID.
func (m *Manager) List(userID string) []*JobStatus {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var jobStatusList []*JobStatus

	if userID == Superuser {
		for _, job := range m.allJobsByJobID {
			jobStatusList = append(jobStatusList, job.Status())
		}
	} else {
		if l2map, exists := m.jobsByUserByJobID[userID]; exists {
			jobStatusList = make([]*JobStatus, 0, len(l2map))

			for _, job := range l2map {
				jobStatusList = append(jobStatusList, job.Status())
			}
		}
	}

	return jobStatusList
}

// Status returns the status of the job with the given JobID owned by
// the given userID.
func (m *Manager) Status(userID, jobID string) (*JobStatus, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	job, err := m.findJobByUser(userID, jobID)
	if err != nil {
		return nil, err
	}
	return job.Status(), nil
}

// StdoutStream returns an io.ByteStream for reading the standard output generated
// by the job with the given jobID own by the given userID.
func (m *Manager) StdoutStream(userID, jobID string) (*io.ByteStream, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	job, err := m.findJobByUser(userID, jobID)
	if err != nil {
		return nil, err
	}
	return job.StdoutStream(), nil

}

// Stderr returns an io.ByteStream for reading the standard error generated
// by the job with the given jobID own by the given userID.
func (m *Manager) StderrStream(userID, jobID string) (*io.ByteStream, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	job, err := m.findJobByUser(userID, jobID)
	if err != nil {
		return nil, err
	}
	return job.StderrStream(), nil
}

// findJobByUser finds a the job with the given jobID that is owned by
// the given userID.  If no such job is found, it returns an error.
// The caller must own the read lock associated with the given Manager.
func (m *Manager) findJobByUser(userID, jobID string) (Job, error) {
	if userID == Superuser {
		if job, exists := m.allJobsByJobID[jobID]; exists {
			return job, nil
		}
	} else {
		if l2map, exists := m.jobsByUserByJobID[userID]; exists {
			if job, exists := l2map[jobID]; exists {
				return job, nil
			}
		}
	}

	return nil, fmt.Errorf("job '%s' does not exist", jobID)
}
