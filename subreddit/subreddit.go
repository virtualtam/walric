package subreddit

type Subreddit struct {
	ID   int
	Name string
}

type SubredditStats struct {
	Name        string
	Submissions int
}
