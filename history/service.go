package history

import (
	"time"

	"github.com/virtualtam/redwall2/submission"
)

var _ Repository = &Service{}

type Service struct {
	Repository
}

func (s *Service) Save(submission *submission.Submission) error {
	now := time.Now().UTC()

	entry := &Entry{
		Date:         now,
		SubmissionID: submission.ID,
	}

	return s.Create(entry)
}

func NewService(repository Repository) *Service {
	validator := newValidator(repository)

	return &Service{
		Repository: validator,
	}
}
