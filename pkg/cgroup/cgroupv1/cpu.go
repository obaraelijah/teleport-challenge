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

// CpuController implements cgroup v2 control using CFS Bandwidth Control.
// See doc/Documentation/scheduler/sched-bwc.txt in the kernel source tree for
// additional information.
//
// This implementation exposes that functionality in terms of how much of the
// available CPU resources a collection of processes can use (0.5 = half a CPU,
// 1.0 = 1 CPU, 1.5 = 1 and a half CPUs, ...).
// The period is defaultPeriodUs and the quota is cpus*period.
type CpuController struct {
	OsAdapter *os.Adapter
	Cpus      float64
}

func (CpuController) Name() string {
	return "cpu"
}

func (c *CpuController) Apply(path string) error {
	if c.Cpus != 0 {

		filename := fmt.Sprintf("%s/%s", path, CpuPeriodFilename)
		if err := c.OsAdapter.WriteFile(filename, defaultPeriodBytes, os.FileMode(0644)); err != nil {
			return err
		}
		filename = fmt.Sprintf("%s/%s", path, CpuQuotaFilename)
		quota := fmt.Sprintf("%d", int(c.Cpus*defaultPeriodUs))
		if err := c.OsAdapter.WriteFile(filename, []byte(quota), os.FileMode(0644)); err != nil {
			return err
		}
	}
	return nil
}
