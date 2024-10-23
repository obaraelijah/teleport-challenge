package cgroup

import "github.com/obaraelijah/teleport-challenge/pkg/adaptation/os"

const (
	MemoryLimitInBytesFilename = "memory.limit_in_bytes"
)

// memory configures the memory cgroup controller.
type memory struct {
	base
	limit *string
}

// NewMemoryController creates an returns a new memory cgroup controller.
func NewMemoryController() *memory {
	return NewMemoryControllerDetailed(nil)
}

func NewMemoryControllerDetailed(osAdapter *os.Adapter) *memory {
	return &memory{
		base: newBase("memory", osAdapter),
	}
}

// SetLimit sets the memory limit in bytes that this cgroup controller will
// enforce.
func (m *memory) SetLimit(value string) *memory {
	m.limit = &value

	return m
}

func (m *memory) Apply(path string) error {
	if m.limit != nil {
		if err := m.write([]byte(*m.limit), "%s/%s", path, MemoryLimitInBytesFilename); err != nil {
			return err
		}
	}

	return nil
}
