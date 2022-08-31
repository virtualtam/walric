package history

// Repository defines the basic operations available to access and persist
// history Entries.
type Repository interface {
	// All returns all persisted Entries.
	All() ([]*Entry, error)

	// Current returns the last chosen Entry.
	Current() (*Entry, error)

	// Create creates and persists an Entry.
	Create(entry *Entry) error
}
