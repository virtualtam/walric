package subreddit

// Repository defines the basic operations available to access and persist
// Subreddits.
type Repository interface {
	// All returns all persisted Subreddits.
	All() ([]*Subreddit, error)

	// Stats returns the aggregated usage statistics for all Subreddits.
	Stats() ([]SubredditStats, error)

	// ByID returns the Subreddit for a given ID.
	ByID(id int) (*Subreddit, error)

	// ByName returns the Subreddit for a given Name.
	ByName(name string) (*Subreddit, error)

	// Create creates and persists a Subreddit.
	Create(subreddit *Subreddit) error
}
