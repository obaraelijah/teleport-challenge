package os

import (
	goos "os"
)

type FileMode = goos.FileMode

var Args = goos.Args

// Adapter serves as a shim between between callers of standard os.* APIs
// and the functions themselves.  The default behavior is simply to dispatch
// the call to the corresponding os.* function.
//
// If the field associated with that function is non-nil, then the adapter will
// dispatch to that function instead.
type Adapter struct {
	MkdirAllFn  func(path string, perm goos.FileMode) error
	RemoveFn    func(name string) error
	WriteFileFn func(name string, data []byte, perm goos.FileMode) error
	GetpidFn    func() int
	EnvironFn   func() []string
}

func (a *Adapter) MkdirAll(path string, perm goos.FileMode) error {
	fn := goos.MkdirAll

	if a != nil && a.MkdirAllFn != nil {
		fn = a.MkdirAllFn
	}

	return fn(path, perm)
}

func (a *Adapter) Remove(name string) error {
	fn := goos.Remove

	if a != nil && a.RemoveFn != nil {
		fn = a.RemoveFn
	}

	return fn(name)
}

func (a *Adapter) WriteFile(name string, data []byte, perm goos.FileMode) error {
	fn := goos.WriteFile

	if a != nil && a.WriteFileFn != nil {
		fn = a.WriteFileFn
	}

	return fn(name, data, perm)
}

func (a *Adapter) Getpid() int {
	fn := goos.Getpid
	if a != nil && a.GetpidFn != nil {
		fn = a.GetpidFn
	}

	return fn()
}

func (a *Adapter) Environ() []string {
	fn := goos.Environ
	if a != nil && a.EnvironFn != nil {
		fn = a.EnvironFn
	}

	return fn()
}
