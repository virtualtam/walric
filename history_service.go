package redwall

import "time"

var _ HistoryRepository = &HistoryService{}

type HistoryService struct {
	HistoryRepository
}

func (s *HistoryService) Save(submission *Submission) error {
	now := time.Now().UTC()

	entry := &HistoryEntry{
		Date:         now,
		SubmissionID: submission.ID,
	}

	return s.Create(entry)
}

func NewHistoryService(historyRepository HistoryRepository) *HistoryService {
	validator := newHistoryValidator(historyRepository)

	return &HistoryService{
		HistoryRepository: validator,
	}
}
