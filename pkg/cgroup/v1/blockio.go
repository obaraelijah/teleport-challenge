package cgroup

import "github.com/obaraelijah/teleport-challenge/pkg/adaptation/os"

const (
	BlkioThrottleReadBpsDevice  = "blkio.throttle.read_bps_device"
	BlkioThrottleWriteBpsDevice = "blkio.throttle.write_bps_device"
)

type blockIo struct {
	base
	// TODO: It might be helpful to enable these to be lists so that a single
	// blockIo object can handle multiple devices.
	readBpsDevice  *string
	writeBpsDevice *string
}

func NewBlockIoController() *blockIo {
	return NewBlockIoControllerDetailed(nil)
}

func NewBlockIoControllerDetailed(osAdapter *os.Adapter) *blockIo {
	return &blockIo{
		base: newBase("blkio", osAdapter),
	}
}

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

// Fluent interface for setting read limit
func (b *blockIo) SetReadBpsDevice(value string) *blockIo {
	b.readBpsDevice = &value

	return b
}

// Fluent interface for setting write limit
func (b *blockIo) SetWriteBpsDevice(value string) *blockIo {
	b.writeBpsDevice = &value

	return b
}
