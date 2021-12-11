package redwall

var _ HistoryRepository = &historyValidator{}

type historyValidator struct {
	HistoryRepository
}

func newHistoryValidator(historyRepository HistoryRepository) *historyValidator {
	return &historyValidator{
		HistoryRepository: historyRepository,
	}
}
