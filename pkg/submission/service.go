package submission

import (
	"errors"
	"strings"

	"github.com/virtualtam/walric/pkg/monitor"
)

// Service handles domain operations for Submission and Subreddit management.
type Service struct {
	r Repository
}

// NewService creates and initializes a Submission and Subreddit Service.
func NewService(repository Repository) *Service {
	return &Service{
		r: repository,
	}
}

// ByPostID returns the Submission matching a given ID.
func (s *Service) ByID(id int) (*Submission, error) {
	submission := &Submission{ID: id}

	if err := submission.requirePositiveID(); err != nil {
		return &Submission{}, err
	}

	submission, err := s.r.SubmissionGetByID(id)
	if err != nil {
		return &Submission{}, err
	}

	subreddit, err := s.subredditByID(submission.Subreddit.ID)
	if err != nil {
		return &Submission{}, err
	}

	submission.Subreddit = subreddit

	return submission, nil
}

// ByMinResolution returns all Submissions whose size is greater or equal to
// the provided minimum resolution.
func (s *Service) ByMinResolution(minResolution *monitor.Resolution) ([]*Submission, error) {
	if err := minResolution.Validate(); err != nil {
		return []*Submission{}, err
	}

	submissions, err := s.r.SubmissionGetByMinResolution(minResolution)
	if err != nil {
		return []*Submission{}, err
	}

	for _, submission := range submissions {
		subreddit, err := s.subredditByID(submission.Subreddit.ID)
		if err != nil {
			return []*Submission{}, err
		}

		submission.Subreddit = subreddit
	}

	return submissions, nil
}

// ByPostID returns the Submission matching a given post ID.
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

	subreddit, err := s.subredditByID(submission.Subreddit.ID)
	if err != nil {
		return &Submission{}, err
	}

	submission.Subreddit = subreddit

	return submission, nil
}

// Creates creates a new Submission.
func (s *Service) Create(submission *Submission) error {
	submission.Normalize()

	if err := submission.ValidateForAddition(s.r); err != nil {
		return err
	}

	return s.r.SubmissionCreate(submission)
}

// Search returns all Submissions whose title match the search string.
func (s *Service) Search(text string) ([]*Submission, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return []*Submission{}, ErrSubmissionSearchTextEmpty
	}

	submissions, err := s.r.SubmissionSearch(text)
	if err != nil {
		return []*Submission{}, err
	}

	for _, submission := range submissions {
		subreddit, err := s.subredditByID(submission.Subreddit.ID)
		if err != nil {
			return []*Submission{}, err
		}

		submission.Subreddit = subreddit
	}

	return submissions, nil
}

// Random returns a randomly selected Submission with a size greater or equal
// to the provided minimum resolution.
func (s *Service) Random(minResolution *monitor.Resolution) (*Submission, error) {
	if err := minResolution.Validate(); err != nil {
		return &Submission{}, err
	}

	submission, err := s.r.SubmissionGetRandom(minResolution)
	if err != nil {
		return &Submission{}, err
	}

	subreddit, err := s.subredditByID(submission.Subreddit.ID)
	if err != nil {
		return &Submission{}, err
	}

	submission.Subreddit = subreddit

	return submission, err
}

// Stats returns statistics about how many Submissions were gathered per Subreddit.
func (s *Service) Stats() ([]SubredditStats, error) {
	return s.r.SubredditGetStats()
}

func (s *Service) subredditByID(id int) (*Subreddit, error) {
	sr := &Subreddit{ID: id}

	if err := sr.requirePositiveID(); err != nil {
		return &Subreddit{}, err
	}

	return s.r.SubredditGetByID(id)
}

// SubredditCreate creates a new Subreddit.
func (s *Service) SubredditCreate(sr *Subreddit) error {
	sr.Normalize()

	if err := sr.ValidateForAddition(s.r); err != nil {
		return err
	}

	return s.r.SubredditCreate(sr)
}

// SubredditGetOrCreateByName returns an existing Subreddit or creates it otherwise.
func (s *Service) SubredditGetOrCreateByName(name string) (*Subreddit, error) {
	subreddit, err := s.SubredditByName(name)

	if errors.Is(err, ErrSubredditNotFound) {
		subreddit = &Subreddit{Name: name}
		if err = s.SubredditCreate(subreddit); err != nil {
			return &Subreddit{}, err
		}

		subreddit, err = s.SubredditByName(name)
		if err != nil {
			return &Subreddit{}, err
		}
	} else if err != nil {
		return &Subreddit{}, err
	}

	return subreddit, nil
}

// SubredditByName returns the SUbreddit for a given name.
func (s *Service) SubredditByName(name string) (*Subreddit, error) {
	sr := &Subreddit{Name: name}
	sr.Normalize()

	if err := sr.requireName(); err != nil {
		return &Subreddit{}, err
	}

	return s.r.SubredditGetByName(sr.Name)
}
