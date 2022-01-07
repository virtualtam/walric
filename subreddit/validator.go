package subreddit

import "errors"

var _ Repository = &validator{}

type validationFn func(subreddit *Subreddit) error

type validator struct {
	Repository
}

func (v *validator) runValidationFns(subreddit *Subreddit, fns ...validationFn) error {
	for _, fn := range fns {
		if err := fn(subreddit); err != nil {
			return err
		}
	}

	return nil
}

func (v *validator) requirePositiveID(subreddit *Subreddit) error {
	if subreddit.ID < 0 {
		return errors.New("Negative ID")
	}

	return nil
}

func (v *validator) ByID(id int) (*Subreddit, error) {
	subreddit := &Subreddit{ID: id}

	err := v.runValidationFns(
		subreddit,
		v.requirePositiveID,
	)
	if err != nil {
		return &Subreddit{}, err
	}

	return v.Repository.ByID(id)
}

func newValidator(repository Repository) *validator {
	return &validator{
		Repository: repository,
	}
}
