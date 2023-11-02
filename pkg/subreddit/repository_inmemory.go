package subreddit

import "errors"

var _ Repository = &RepositoryInMemory{}

// repositoryInMemory provides an in-memory Repository for testing.
type RepositoryInMemory struct {
	currentID  int
	subreddits []*Subreddit
}

func (r *RepositoryInMemory) SubredditCreate(subreddit *Subreddit) error {
	subreddit.ID = r.currentID
	r.currentID++

	r.subreddits = append(r.subreddits, subreddit)

	return nil
}

func (r *RepositoryInMemory) SubredditGetAll() ([]*Subreddit, error) {
	return r.subreddits, nil
}

func (r *RepositoryInMemory) SubredditGetStats() ([]SubredditStats, error) {
	return []SubredditStats{}, errors.New("not implemented")
}

func (r *RepositoryInMemory) SubredditGetByID(id int) (*Subreddit, error) {
	for _, subreddit := range r.subreddits {
		if subreddit.ID == id {
			return subreddit, nil
		}
	}

	return &Subreddit{}, ErrNotFound
}

func (r *RepositoryInMemory) SubredditGetByName(name string) (*Subreddit, error) {
	for _, subreddit := range r.subreddits {
		if subreddit.Name == name {
			return subreddit, nil
		}
	}

	return &Subreddit{}, ErrNotFound
}

func (r *RepositoryInMemory) SubredditIsNameRegistered(name string) (bool, error) {
	for _, subreddit := range r.subreddits {
		if subreddit.Name == name {
			return true, nil
		}
	}

	return false, nil
}

// NewRepositoryInMemory initializes and returns an in-memory repository for
// testing.
func NewRepositoryInMemory(subreddits []*Subreddit) *RepositoryInMemory {
	return &RepositoryInMemory{
		currentID:  len(subreddits) + 1,
		subreddits: subreddits,
	}
}
