package history

var _ Repository = &repositoryInMemory{}

// repositoryInMemory provides an in-memory Repository for testing.
type repositoryInMemory struct {
	entries []*Entry
}

func (r *repositoryInMemory) HistoryGetAll() ([]*Entry, error) {
	return r.entries, nil
}

func (r *repositoryInMemory) HistoryGetCurrent() (*Entry, error) {
	if len(r.entries) == 0 {
		return &Entry{}, ErrNotFound
	}

	return r.entries[len(r.entries)-1], nil
}

func (r *repositoryInMemory) HistoryCreate(entry *Entry) error {
	r.entries = append(r.entries, entry)
	return nil
}
