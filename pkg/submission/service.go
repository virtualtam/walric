package submission

import (
	"strings"

	"github.com/virtualtam/walric/pkg/monitor"
	"github.com/virtualtam/walric/pkg/subreddit"
)

// Service handles domain operations for Submission management.
type Service struct {
	r Repository

	subredditService *subreddit.Service
}

// NewService creates and initializes a Submission Service.
func NewService(repository Repository, subredditService *subreddit.Service) *Service {
	return &Service{
		r:                repository,
		subredditService: subredditService,
	}
}

func (s *Service) ByID(id int) (*Submission, error) {
	submission := &Submission{ID: id}

	if err := submission.requirePositiveID(); err != nil {
		return &Submission{}, err
	}

	submission, err := s.r.SubmissionGetByID(id)
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

func (s *Service) ByMinResolution(minResolution *monitor.Resolution) ([]*Submission, error) {
	if minResolution.HeightPx < 1 || minResolution.WidthPx < 1 {
		return []*Submission{}, ErrResolutionInvalid
	}

	submissions, err := s.r.SubmissionGetByMinResolution(minResolution)
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

func (s *Service) ByPostID(postID string) (*Submission, error) {
	submission := &Submission{PostID: postID}
	submission.Normalize()

	if err := submission.requirePostID(); err != nil {
		return &Submission{}, err
	}

	submission, err := s.r.SubmissionGetByPostID(postID)
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

func (s *Service) Create(submission *Submission) error {
	submission.Normalize()

	if err := submission.ValidateForAddition(s.r); err != nil {
		return err
	}

	return s.r.SubmissionCreate(submission)
}

func (s *Service) Search(text string) ([]*Submission, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return []*Submission{}, ErrSearchTextEmpty
	}

	submissions, err := s.r.SubmissionSearch(text)
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
	if minResolution.HeightPx < 1 || minResolution.WidthPx < 1 {
		return &Submission{}, ErrResolutionInvalid
	}

	submission, err := s.r.SubmissionGetRandom(minResolution)
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
