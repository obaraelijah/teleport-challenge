package ostest

type RemoveRecord struct {
	Path string
}

// RemoveMock is a component that provides a mock implementation of the
// os.Remove() function.  The implementation records the paramters received
// and returns the configured NextError.
type RemoveMock struct {
	Events    []*RemoveRecord
	NextError error
}

func (w *RemoveMock) Remove(path string) error {
	w.Events = append(w.Events, &RemoveRecord{
		Path: path,
	})
	return w.NextError
}
