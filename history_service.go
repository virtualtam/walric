package redwall

var _ HistoryRepository = &HistoryService{}

type HistoryService struct {
	HistoryRepository
}

func NewHistoryService(historyRepository HistoryRepository) *HistoryService {
	validator := newHistoryValidator(historyRepository)

	return &HistoryService{
		HistoryRepository: validator,
	}
}
