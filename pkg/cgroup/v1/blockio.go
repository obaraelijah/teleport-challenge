package cgroup

import "github.com/obaraelijah/teleport-challenge/pkg/adaptation/os"

const (
	BlkioThrottleReadBpsDevice  = "blkio.throttle.read_bps_device"
	BlkioThrottleWriteBpsDevice = "blkio.throttle.write_bps_device"
)

// blockIo implements the BlockIO cgroup controller
type blockIo struct {
	base
	// TODO: It might be helpful to enable these to be lists so that a single
	// blockIo object can handle multiple devices.
	readBpsDevice  *string
	writeBpsDevice *string
}

// NewBlockIoController creates and returns a new blockIo cgroup controller
func NewBlockIoController() *blockIo {
	return NewBlockIoControllerDetailed(nil)
}

func NewBlockIoControllerDetailed(osAdapter *os.Adapter) *blockIo {
	return &blockIo{
		base: newBase("blkio", osAdapter),
	}
}

// Apply applies this cgroup controller configuration to the blkio cgroup
// at the given path.
func (b *blockIo) Apply(path string) error {
	if b.readBpsDevice != nil {
		if err := b.write([]byte(*b.readBpsDevice), "%s/%s", path, BlkioThrottleReadBpsDevice); err != nil {
			return err
		}
	}

	if b.writeBpsDevice != nil {
		if err := b.write([]byte(*b.writeBpsDevice), "%s/%s", path, BlkioThrottleWriteBpsDevice); err != nil {
			return err
		}
	}

	return nil
}

// SetReadBpsDevice Sets the read bytes per second limit for a device.
func (b *blockIo) SetReadBpsDevice(value string) *blockIo {
	b.readBpsDevice = &value

	return b
}

// SetWriteBpsDevice Sets the write bytes per second limit for a device.
func (b *blockIo) SetWriteBpsDevice(value string) *blockIo {
	b.writeBpsDevice = &value

	return b
}
