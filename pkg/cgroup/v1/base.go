package cgroup

import (
	"fmt"

	"github.com/obaraelijah/teleport-challenge/pkg/adaptation/os"
)

const DefaultFileMode os.FileMode = 0644

// base is the base type for all cgroup controllers
type base struct {
	osAdapter *os.Adapter
	// name is the controller name
	name string
}

func newBase(name string, osAdapter *os.Adapter) base {
	return base{
		osAdapter: osAdapter,
		name:      name,
	}
}

// Name returns the name of this cgroup controller
func (b *base) Name() string {
	return b.name
}

// write update the cgroup file constructed with the given pathFmt and args with the given value.
func (b *base) write(value []byte, pathFmt string, args ...interface{}) error {
	return b.osAdapter.WriteFile(fmt.Sprintf(pathFmt, args...), value, DefaultFileMode)
}
