package subreddit

import "errors"

var _ Repository = &RepositoryInMemory{}

type RepositoryInMemory struct {
	currentID  int
	subreddits []*Subreddit
}

func (r *RepositoryInMemory) All() ([]*Subreddit, error) {
	return r.subreddits, nil
}

func (r *RepositoryInMemory) Stats() ([]SubredditStats, error) {
	return []SubredditStats{}, errors.New("not implemented")
}

func (r *RepositoryInMemory) ByID(id int) (*Subreddit, error) {
	for _, subreddit := range r.subreddits {
		if subreddit.ID == id {
			return subreddit, nil
		}
	}

	return &Subreddit{}, ErrNotFound
}

func (r *RepositoryInMemory) ByName(name string) (*Subreddit, error) {
	for _, subreddit := range r.subreddits {
		if subreddit.Name == name {
			return subreddit, nil
		}
	}

	return &Subreddit{}, ErrNotFound
}

func (r *RepositoryInMemory) Create(subreddit *Subreddit) error {
	subreddit.ID = r.currentID
	r.currentID++

	r.subreddits = append(r.subreddits, subreddit)

	return nil
}

func NewRepositoryInMemory(subreddits []*Subreddit) *RepositoryInMemory {
	return &RepositoryInMemory{
		currentID:  len(subreddits) + 1,
		subreddits: subreddits,
	}
}
