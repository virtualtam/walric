package redwall

type HistoryRepository interface {
	All() ([]HistoryEntry, error)
	Current() (*HistoryEntry, error)

	Create(entry *HistoryEntry) error
}
