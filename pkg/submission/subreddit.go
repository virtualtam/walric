package submission

import (
	"strings"
)

// Subreddit represents a Reddit subreddit.
type Subreddit struct {
	ID   int
	Name string
}

// Normalize sanitizes and normalizes all fields.
func (sr *Subreddit) Normalize() {
	sr.normalizeName()
}

// ValidateForAddition ensures mandatory fields are properly set when adding an
// new Subreddit.
func (sr *Subreddit) ValidateForAddition(r ValidationRepository) error {
	fns := []func() error{
		sr.requireDefaultID,
		sr.requireName,
		sr.ensureNameIsNotRegistered(r),
	}

	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (sr *Subreddit) ensureNameIsNotRegistered(r ValidationRepository) func() error {
	return func() error {
		registered, err := r.SubredditIsNameRegistered(sr.Name)
		if err != nil {
			return err
		}

		if registered {
			return ErrSubredditNameAlreadyRegistered
		}

		return nil
	}
}

func (sr *Subreddit) normalizeName() {
	sr.Name = strings.TrimSpace(sr.Name)
}

func (sr *Subreddit) requireDefaultID() error {
	if sr.ID != 0 {
		return ErrSubredditIDInvalid
	}

	return nil
}

func (sr *Subreddit) requirePositiveID() error {
	if sr.ID <= 0 {
		return ErrSubredditIDInvalid
	}

	return nil
}

func (sr *Subreddit) requireName() error {
	if sr.Name == "" {
		return ErrSubredditNameEmpty
	}

	return nil
}
