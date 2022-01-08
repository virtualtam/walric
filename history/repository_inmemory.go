package history

var _ Repository = &repositoryInMemory{}

type repositoryInMemory struct {
	entries []*Entry
}

func (r *repositoryInMemory) All() ([]*Entry, error) {
	return r.entries, nil
}

func (r *repositoryInMemory) Current() (*Entry, error) {
	if len(r.entries) == 0 {
		return &Entry{}, ErrNotFound
	}

	return r.entries[len(r.entries)-1], nil
}

func (r *repositoryInMemory) Create(entry *Entry) error {
	r.entries = append(r.entries, entry)
	return nil
}
