package subreddit

// Subreddit represents a Reddit subreddit.
type Subreddit struct {
	ID   int
	Name string
}

// SubredditStats holds the aggregated usage statistics for a given Subreddit.
type SubredditStats struct {
	Name        string
	Submissions int
}
