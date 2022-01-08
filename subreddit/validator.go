package subreddit

import "strings"

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
	if subreddit.ID <= 0 {
		return ErrIDInvalid
	}

	return nil
}

func (v *validator) normalizeName(subreddit *Subreddit) error {
	subreddit.Name = strings.TrimSpace(subreddit.Name)

	return nil
}

func (v *validator) requireName(subreddit *Subreddit) error {
	if subreddit.Name == "" {
		return ErrNameEmpty
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

func (v *validator) ByName(name string) (*Subreddit, error) {
	subreddit := &Subreddit{Name: name}

	err := v.runValidationFns(
		subreddit,
		v.normalizeName,
		v.requireName,
	)
	if err != nil {
		return &Subreddit{}, err
	}

	return v.Repository.ByName(subreddit.Name)
}

func newValidator(repository Repository) *validator {
	return &validator{
		Repository: repository,
	}
}
