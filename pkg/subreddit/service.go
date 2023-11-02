package subreddit

import "errors"

// Service handles domain operations for Subreddit management.
type Service struct {
	r Repository
}

// NewService creates and initializes a Subreddit Service.
func NewService(r Repository) *Service {
	return &Service{
		r: r,
	}
}

func (s *Service) ByID(id int) (*Subreddit, error) {
	sr := &Subreddit{ID: id}

	if err := sr.requirePositiveID(); err != nil {
		return &Subreddit{}, err
	}

	return s.r.SubredditGetByID(id)
}

func (s *Service) ByName(name string) (*Subreddit, error) {
	sr := &Subreddit{Name: name}
	sr.Normalize()

	if err := sr.requireName(); err != nil {
		return &Subreddit{}, err
	}

	return s.r.SubredditGetByName(sr.Name)
}

func (s *Service) Create(sr *Subreddit) error {
	sr.Normalize()

	if err := sr.ValidateForAddition(s.r); err != nil {
		return err
	}

	return s.r.SubredditCreate(sr)
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

func (s *Service) Stats() ([]SubredditStats, error) {
	return s.r.SubredditGetStats()
}
