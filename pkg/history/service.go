package history

import (
	"time"

	"github.com/virtualtam/walric/pkg/submission"
)

// Service handles domain operations for Entry management.
type Service struct {
	*validator

	submissionService *submission.Service
}

func (s *Service) All() ([]*Entry, error) {
	entries, err := s.validator.All()
	if err != nil {
		return []*Entry{}, err
	}

	for _, entry := range entries {
		submission, err := s.submissionService.ByID(entry.SubmissionID)
		if err != nil {
			return []*Entry{}, err
		}

		entry.Submission = submission
	}

	return entries, nil
}

func (s *Service) Current() (*Entry, error) {
	entry, err := s.validator.Current()
	if err != nil {
		return &Entry{}, err
	}

	submission, err := s.submissionService.ByID(entry.SubmissionID)
	if err != nil {
		return &Entry{}, err
	}

	entry.Submission = submission

	return entry, nil
}

func (s *Service) Save(submission *submission.Submission) error {
	now := time.Now().UTC()

	entry := &Entry{
		Date:         now,
		SubmissionID: submission.ID,
	}

	return s.Create(entry)
}

// NewService creates and initializes an Entry Service.
func NewService(repository Repository, submissionService *submission.Service) *Service {
	validator := newValidator(repository)

	return &Service{
		validator:         validator,
		submissionService: submissionService,
	}
}
