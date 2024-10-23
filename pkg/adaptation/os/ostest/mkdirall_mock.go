package ostest

import "github.com/obaraelijah/teleport-challenge/pkg/adaptation/os"

type MkdirAllRecord struct {
	Path string
	Perm os.FileMode
}

// MkdirAllMock is a component that provides a mock implementation of the
// os.MkdirAll() function.  The implementation records the paramters received
// and returns the configured NextError.
type MkdirAllMock struct {
	Events    []*MkdirAllRecord
	NextError error
}

func (w *MkdirAllMock) MkdirAll(path string, perm os.FileMode) error {
	w.Events = append(w.Events, &MkdirAllRecord{
		Path: path,
		Perm: perm,
	})

	return w.NextError
}
