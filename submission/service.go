package submission

import (
	"github.com/virtualtam/walric/monitor"
	"github.com/virtualtam/walric/subreddit"
)

var _ Repository = &Service{}

type Service struct {
	Repository

	subredditService *subreddit.Service
}

func (s *Service) ByID(id int) (*Submission, error) {
	submission, err := s.Repository.ByID(id)
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
	submission, err := s.Repository.ByPostID(postID)
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
	submissions, err := s.Repository.Search(text)
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
	submissions, err := s.Repository.ByMinResolution(minResolution)
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
	submission, err := s.Repository.Random(minResolution)
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

func NewService(repository Repository, subredditService *subreddit.Service) *Service {
	validator := newValidator(repository)

	return &Service{
		Repository:       validator,
		subredditService: subredditService,
	}
}
