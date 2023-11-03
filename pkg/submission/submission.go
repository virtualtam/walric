package submission

import (
	"fmt"
	"strings"
	"time"
)

// Submission represents the metadata for a Reddit post with an image
// attachment, and the metadata for the corresponding local file.
type Submission struct {
	ID int

	Subreddit *Subreddit

	// Reddit post metadata
	Author    string
	Permalink string
	PostID    string
	PostedAt  time.Time
	Score     int
	Title     string

	// Attached image metadata
	ImageDomain string
	ImageURL    string
	ImageNSFW   bool

	// Local image metadata
	ImageFilename string
	ImageHeightPx int
	ImageWidthPx  int
}

// Normalize sanitizes and normalizes all fields.
func (s *Submission) Normalize() {
	s.normalizePostID()
	s.normalizeTitle()
}

// PermalinkURL returns the Reddit permalink for this submission's post.
func (s *Submission) PermalinkURL() string {
	return fmt.Sprintf("https://reddit.com%s", s.Permalink)
}

// User returns the Reddit-formatted username for this submission's author.
func (s *Submission) User() string {
	return fmt.Sprintf("u/%s", s.Author)
}

// ValidateForAddition ensures mandatory fields are properly set when adding an
// new Submission.
func (s *Submission) ValidateForAddition(r ValidationRepository) error {
	fns := []func() error{
		s.requirePositiveSubredditID,
		s.requireDefaultID,
		s.requirePostID,
		s.ensurePostIDIsNotRegistered(r),
		s.requireTitle,
	}

	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Submission) normalizePostID() {
	s.PostID = strings.TrimSpace(s.PostID)
}

func (s *Submission) normalizeTitle() {
	s.Title = strings.TrimSpace(s.Title)
}

func (s *Submission) requirePositiveSubredditID() error {
	if s.Subreddit.ID <= 0 {
		return ErrSubredditIDInvalid
	}

	return nil
}

func (s *Submission) requireDefaultID() error {
	if s.ID != 0 {
		return ErrSubmissionIDInvalid
	}

	return nil
}

func (s *Submission) requirePositiveID() error {
	if s.ID <= 0 {
		return ErrSubmissionIDInvalid
	}

	return nil
}

func (s *Submission) ensurePostIDIsNotRegistered(r ValidationRepository) func() error {
	return func() error {
		registered, err := r.SubmissionIsPostIDRegistered(s.PostID)

		if err != nil {
			return err
		}

		if registered {
			return ErrSubmissionPostIDAlreadyRegistered
		}

		return nil
	}
}

func (s *Submission) requirePostID() error {
	if s.PostID == "" {
		return ErrSubmissionPostIDEmpty
	}

	return nil
}

func (s *Submission) requireTitle() error {
	if s.Title == "" {
		return ErrSubmissionTitleEmpty
	}

	return nil
}
