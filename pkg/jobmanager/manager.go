package jobmanager

import (
	"fmt"
	"sync"

	"github.com/obaraelijah/teleport-challenge/pkg/cgroup/v1"
)

type JobConstructor func(
	jobName string,
	controllers []cgroup.Controller,
	programPath string,
	arguments ...string,
) Job

type Manager struct {
	mutex          sync.RWMutex
	jobsByUser     map[string]map[string]Job // userId->jobId->job
	allJobs        map[string]Job            // jobId->job
	controllers    []cgroup.Controller
	jobConstructor JobConstructor
}

func NewManager() *Manager {
	// TODO:
	readLimit := fmt.Sprintf("8:16 %d", 1024*1024*20)
	writeLimit := fmt.Sprintf("8:16 %d", 1024*1024*40)

	controllers := []cgroup.Controller{
		cgroup.NewCpuController().SetCpus(0.5),
		cgroup.NewMemoryController().SetLimit("2M"),
		cgroup.NewBlockIoController().
			SetReadBpsDevice(readLimit).
			SetWriteBpsDevice(writeLimit),
	}

	return NewManagerDetailed(NewJob, controllers)
}

func NewManagerDetailed(jobConstructor JobConstructor, controllers []cgroup.Controller) *Manager {
	return &Manager{
		jobsByUser:     make(map[string]map[string]Job),
		allJobs:        make(map[string]Job),
		controllers:    controllers,
		jobConstructor: jobConstructor,
	}
}
