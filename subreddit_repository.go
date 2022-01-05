package redwall

type SubredditRepository interface {
	All() ([]Subreddit, error)
	Stats() ([]SubredditStats, error)

	ByID(id int) (*Subreddit, error)
	ByName(name string) (*Subreddit, error)
}
