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

func NewMemoryDetailed(osAdapter *os.Adapter) *memory {
	return &memory{
		base: newBase("memory", osAdapter),
	}
}

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
