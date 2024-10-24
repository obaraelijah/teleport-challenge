package cgroupv1

import "github.com/obaraelijah/teleport-challenge/pkg/adaptation/os"

const (
	BlkioThrottleReadBpsDevice  = "blkio.throttle.read_bps_device"
	BlkioThrottleWriteBpsDevice = "blkio.throttle.write_bps_device"
)

// blockIO implements the BlockIO cgroup controller
type blockIO struct {
	base
	// TODO: It might be helpful to enable these to be lists so that a single
	// blocIO object can handle multiple devices
	readBpsDevice  *string
	writeBpsDevice *string
}

// NewBlockIOController creates and returns a new blockIO cgroup controller
func NewBlockIOController() *blockIO {
	return NewBlockIOControllerDetailed(nil)
}

func NewBlockIOControllerDetailed(osAdapter *os.Adapter) *blockIO {
	return &blockIO{
		base: newBase("blkio", osAdapter),
	}
}

// Apply applies this cgroup controller configuration to the blkio cgroup
// at the given path.
func (b *blockIO) Apply(path string) error {
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
func (b *blockIO) SetReadBpsDevice(value string) *blockIO {
	b.readBpsDevice = &value

	return b
}

// SetWriteBpsDevice Sets the write bytes per second limit for a device.
func (b *blockIO) SetWriteBpsDevice(value string) *blockIO {
	b.writeBpsDevice = &value

	return b
}
