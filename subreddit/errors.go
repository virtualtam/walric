package subreddit

import "errors"

var (
	ErrIDInvalid             error = errors.New("subreddit: invalid ID")
	ErrNameAlreadyRegistered error = errors.New("subreddit: name already used")
	ErrNameEmpty             error = errors.New("subreddit: empty name")
	ErrNotFound              error = errors.New("subreddit: not found")
)
