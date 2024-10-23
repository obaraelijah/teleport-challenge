package ostest

import "github.com/obaraelijah/teleport-challenge/pkg/adaptation/os"

type WriteFileRecord struct {
	Name string
	Data []byte
	Perm os.FileMode
}

// WriteFileMock is a component that provides a mock implementation of the
// os.WriteFile() function.  The implementation records the paramters received
// and returns the configured NextError.
type WriteFileMock struct {
	Events    []*WriteFileRecord
	NextError error
}

func (w *WriteFileMock) WriteFile(name string, data []byte, perm os.FileMode) error {
	w.Events = append(w.Events, &WriteFileRecord{
		Name: name,
		Data: data,
		Perm: perm,
	})
	return w.NextError
}
