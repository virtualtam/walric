package redwall

type subredditValidator struct {
	SubredditRepository
}

func newSubredditValidator(subredditRepository SubredditRepository) *subredditValidator {
	return &subredditValidator{
		SubredditRepository: subredditRepository,
	}
}
