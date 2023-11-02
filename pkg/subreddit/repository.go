package subreddit

// ValidationRepository provides methods for Subreddit validation.
type ValidationRepository interface {
	// SubredditIsNameRegistered returns whether this Subreddit was previously saved.
	SubredditIsNameRegistered(name string) (bool, error)
}

// Repository defines the basic operations available to access and persist
// Subreddits.
type Repository interface {
	ValidationRepository

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
