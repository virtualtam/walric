package redwall

type HistoryRepository interface {
	All() ([]HistoryEntry, error)
}