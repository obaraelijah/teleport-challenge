package cgroupv1

import (
	"fmt"

	"github.com/obaraelijah/teleport-challenge/pkg/adaptation/os"
)

const (
	MemoryLimitInBytesFilename = "memory.limit_in_bytes"
)

// MemoryController configures the MemoryController cgroup controller.
type MemoryController struct {
	OsAdapter *os.Adapter
	Limit     string
}

func (MemoryController) Name() string {
	return "memory"
}

func (m *MemoryController) Apply(path string) error {
	if m.Limit != "" {
		filename := fmt.Sprintf("%s/%s", path, MemoryLimitInBytesFilename)
		if err := m.OsAdapter.WriteFile(filename, []byte(m.Limit), os.FileMode(0644)); err != nil {
			return err
		}
	}

	return nil
}
