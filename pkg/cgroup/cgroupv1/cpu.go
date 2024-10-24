package cgroupv1

import (
	"fmt"

	"github.com/obaraelijah/teleport-challenge/pkg/adaptation/os"
)

const (
	CpuPeriodFilename = "cpu.cfs_period_us"
	CpuQuotaFilename  = "cpu.cfs_quota_us"

	defaultPeriodUs = 100000
)

var defaultPeriodBytes = []byte(fmt.Sprintf("%d", defaultPeriodUs))

// cpu implements cgroup v2 control using CFS Bandwidth Control.
// See doc/Documentation/scheduler/sched-bwc.txt in the kernel source tree for
// additional information.
//
// This implementation exposes that functionality in terms of how much of the
// available CPU resources a collection of processes can use (0.5 = half a CPU,
// 1.0 = 1 CPU, 1.5 = 1 and a half CPUs, ...).
// The period is defaultPeriodUs and the quota is cpus*period.
type cpu struct {
	base
	cpus *float64
}

// NewCpuController creates and returns a new cpu cgroup controller.
func NewCpuController() *cpu {
	return NewCpuControllerDetailed(nil)
}

func NewCpuControllerDetailed(osAdapter *os.Adapter) *cpu {
	return &cpu{
		base: newBase("cpu", osAdapter),
	}
}

// SetCpus sets the CPU limit in terms of number of CPUs.
func (c *cpu) SetCpus(value float64) *cpu {
	c.cpus = &value
	return c
}

func (c *cpu) Apply(path string) error {
	if c.cpus != nil {
		if err := c.write(defaultPeriodBytes, "%s/%s", path, CpuPeriodFilename); err != nil {
			return err
		}
		quota := fmt.Sprintf("%d", int(*c.cpus*defaultPeriodUs))
		if err := c.write([]byte(quota), "%s/%s", path, CpuQuotaFilename); err != nil {
			return err
		}
	}
	return nil
}
