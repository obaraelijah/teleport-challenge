package ostest

type RemoveRecord struct {
	Path string
}

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
