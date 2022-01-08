package subreddit

import "errors"

var (
	ErrIDInvalid error = errors.New("subreddit: invalid ID")
	ErrNameEmpty error = errors.New("subreddit: empty name")
	ErrNotFound  error = errors.New("subreddit: not found")
)
