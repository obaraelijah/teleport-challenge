package cgroupv1

import (
	"fmt"

	"github.com/obaraelijah/teleport-challenge/pkg/adaptation/os"
)

const (
	BlkioThrottleReadBpsDevice  = "blkio.throttle.read_bps_device"
	BlkioThrottleWriteBpsDevice = "blkio.throttle.write_bps_device"
)

// BlockIOController implements the BlockIOController cgroup controller
type BlockIOController struct {
	OsAdapter      *os.Adapter
	ReadBpsDevice  string
	WriteBpsDevice string
}

func (BlockIOController) Name() string {
	return "blkio"
}

// Apply applies this cgroup controller configuration to the blkio cgroup
// at the given path.
func (b *BlockIOController) Apply(path string) error {
	if b.ReadBpsDevice != "" {
		filename := fmt.Sprintf("%s/%s", path, BlkioThrottleReadBpsDevice)
		if err := b.OsAdapter.WriteFile(filename, []byte(b.ReadBpsDevice), os.FileMode(0644)); err != nil {
			return err
		}
	}

	if b.WriteBpsDevice != "" {
		filename := fmt.Sprintf("%s/%s", path, BlkioThrottleWriteBpsDevice)
		if err := b.OsAdapter.WriteFile(filename, []byte(b.WriteBpsDevice), os.FileMode(0644)); err != nil {
			return err
		}
	}

	return nil
}
