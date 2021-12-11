package redwall

import "errors"

var _ SubredditRepository = &subredditValidator{}

type subredditValidationFn func(subreddit *Subreddit) error

type subredditValidator struct {
	SubredditRepository
}

func (v *subredditValidator) runValidationFns(subreddit *Subreddit, fns ...subredditValidationFn) error {
	for _, fn := range fns {
		if err := fn(subreddit); err != nil {
			return err
		}
	}

	return nil
}

func (v *subredditValidator) requirePositiveID(subreddit *Subreddit) error {
	if subreddit.ID < 0 {
		return errors.New("Negative ID")
	}

	return nil
}

func (v *subredditValidator) ByID(id int) (*Subreddit, error) {
	subreddit := &Subreddit{ID: id}

	err := v.runValidationFns(
		subreddit,
		v.requirePositiveID,
	)
	if err != nil {
		return &Subreddit{}, err
	}

	return v.SubredditRepository.ByID(id)
}

func newSubredditValidator(subredditRepository SubredditRepository) *subredditValidator {
	return &subredditValidator{
		SubredditRepository: subredditRepository,
	}
}
