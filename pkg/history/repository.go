package history

// Repository defines the basic operations available to access and persist
// history Entries.
type Repository interface {
	// HistoryGetAll returns all persisted Entries.
	HistoryGetAll() ([]*Entry, error)

	// HistoryGetCurrent returns the last chosen Entry.
	HistoryGetCurrent() (*Entry, error)

	// HistoryCreate creates and persists an Entry.
	HistoryCreate(entry *Entry) error
}
