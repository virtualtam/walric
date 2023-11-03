package history

import (
	"github.com/virtualtam/walric/pkg/submission"
)

// Service handles domain operations for Entry history management.
type Service struct {
	r Repository

	submissionService *submission.Service
}

// NewService creates and initializes an Entry history Service.
func NewService(r Repository, submissionService *submission.Service) *Service {
	return &Service{
		r:                 r,
		submissionService: submissionService,
	}
}

// All returns the history of all saved entries.
func (s *Service) All() ([]*Entry, error) {
	entries, err := s.r.HistoryGetAll()
	if err != nil {
		return []*Entry{}, err
	}

	for _, entry := range entries {
		submission, err := s.submissionService.ByID(entry.Submission.ID)
		if err != nil {
			return []*Entry{}, err
		}

		entry.Submission = submission
	}

	return entries, nil
}

// Current returns the last selected history Entry.
func (s *Service) Current() (*Entry, error) {
	entry, err := s.r.HistoryGetCurrent()
	if err != nil {
		return &Entry{}, err
	}

	submission, err := s.submissionService.ByID(entry.Submission.ID)
	if err != nil {
		return &Entry{}, err
	}

	entry.Submission = submission

	return entry, nil
}

// Save adds a new Entry to the history.
func (s *Service) Save(entry *Entry) error {
	return s.r.HistoryCreate(entry)
}
