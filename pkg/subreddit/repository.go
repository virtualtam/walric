package subreddit

// Repository defines the basic operations available to access and persist
// Subreddits.
type Repository interface {
	// SubredditGetAll returns all persisted Subreddits.
	SubredditGetAll() ([]*Subreddit, error)

	// SubredditGetStats returns the aggregated usage statistics for all Subreddits.
	SubredditGetStats() ([]SubredditStats, error)

	// SubredditGetByID returns the Subreddit for a given ID.
	SubredditGetByID(id int) (*Subreddit, error)

	// SubredditGetByName returns the Subreddit for a given Name.
	SubredditGetByName(name string) (*Subreddit, error)

	// SubredditCreate creates and persists a Subreddit.
	SubredditCreate(subreddit *Subreddit) error
}
