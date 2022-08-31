package subreddit

import "errors"

// Service handles domain operations for Subreddit management.
type Service struct {
	*validator
}

// NewService creates and initializes a Subreddit Service.
func NewService(repository Repository) *Service {
	validator := newValidator(repository)

	return &Service{
		validator: validator,
	}
}

func (s *Service) GetOrCreateByName(name string) (*Subreddit, error) {
	subreddit, err := s.ByName(name)

	if errors.Is(err, ErrNotFound) {
		subreddit = &Subreddit{Name: name}
		if err = s.Create(subreddit); err != nil {
			return &Subreddit{}, err
		}

		subreddit, err = s.ByName(name)
		if err != nil {
			return &Subreddit{}, err
		}
	} else if err != nil {
		return &Subreddit{}, err
	}

	return subreddit, nil
}
