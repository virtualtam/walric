package submission

import "errors"

var (
	ErrIDInvalid               error = errors.New("submission: invalid ID")
	ErrNotFound                error = errors.New("submission: not found")
	ErrPostIDAlreadyRegistered error = errors.New("submission: post ID already registered")
	ErrPostIDEmpty             error = errors.New("submission: empty post ID")
	ErrResolutionInvalid       error = errors.New("submission: invalid resolution")
	ErrSearchTextEmpty         error = errors.New("submission: empty search text")
	ErrSubredditIDInvalid      error = errors.New("submission: invalid subreddit ID")
	ErrTitleEmpty              error = errors.New("submission: empty title")
)
