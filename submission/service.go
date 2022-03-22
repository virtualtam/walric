package submission

import (
	"github.com/virtualtam/walric/monitor"
	"github.com/virtualtam/walric/subreddit"
)

// Service handles domain operations for Submission management.
type Service struct {
	*validator

	subredditService *subreddit.Service
}

func (s *Service) ByID(id int) (*Submission, error) {
	submission, err := s.validator.ByID(id)
	if err != nil {
		return &Submission{}, err
	}

	subreddit, err := s.subredditService.ByID(submission.SubredditID)
	if err != nil {
		return &Submission{}, err
	}

	submission.Subreddit = subreddit

	return submission, nil
}

func (s *Service) ByPostID(postID string) (*Submission, error) {
	submission, err := s.validator.ByPostID(postID)
	if err != nil {
		return &Submission{}, err
	}

	subreddit, err := s.subredditService.ByID(submission.SubredditID)
	if err != nil {
		return &Submission{}, err
	}

	submission.Subreddit = subreddit

	return submission, nil
}

func (s *Service) Search(text string) ([]*Submission, error) {
	submissions, err := s.validator.Search(text)
	if err != nil {
		return []*Submission{}, err
	}

	for _, submission := range submissions {
		subreddit, err := s.subredditService.ByID(submission.SubredditID)
		if err != nil {
			return []*Submission{}, err
		}

		submission.Subreddit = subreddit
	}

	return submissions, nil
}

func (s *Service) ByMinResolution(minResolution *monitor.Resolution) ([]*Submission, error) {
	submissions, err := s.validator.ByMinResolution(minResolution)
	if err != nil {
		return []*Submission{}, err
	}

	for _, submission := range submissions {
		subreddit, err := s.subredditService.ByID(submission.SubredditID)
		if err != nil {
			return []*Submission{}, err
		}

		submission.Subreddit = subreddit
	}

	return submissions, nil
}

func (s *Service) Random(minResolution *monitor.Resolution) (*Submission, error) {
	submission, err := s.validator.Random(minResolution)
	if err != nil {
		return &Submission{}, err
	}

	subreddit, err := s.subredditService.ByID(submission.SubredditID)
	if err != nil {
		return &Submission{}, err
	}

	submission.Subreddit = subreddit

	return submission, err
}

// NewService creates and initializes a Submission Service.
func NewService(repository Repository, subredditService *subreddit.Service) *Service {
	validator := newValidator(repository)

	return &Service{
		validator:        validator,
		subredditService: subredditService,
	}
}
