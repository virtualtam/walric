package redwall

type SubredditRepository interface {
	All() ([]Subreddit, error)
	Stats() ([]SubredditStats, error)
}
