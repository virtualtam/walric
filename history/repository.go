package history

type Repository interface {
	All() ([]*Entry, error)
	Current() (*Entry, error)

	Create(entry *Entry) error
}
