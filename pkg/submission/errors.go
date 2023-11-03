package submission

import "errors"

var (
	ErrSubmissionIDInvalid               error = errors.New("submission: invalid ID")
	ErrSubmissionNotFound                error = errors.New("submission: not found")
	ErrSubmissionPostIDAlreadyRegistered error = errors.New("submission: post ID already registered")
	ErrSubmissionPostIDEmpty             error = errors.New("submission: empty post ID")
	ErrSubmissionSearchTextEmpty         error = errors.New("submission: empty search text")
	ErrSubmissionTitleEmpty              error = errors.New("submission: empty title")

	ErrSubredditIDInvalid             error = errors.New("subreddit: invalid ID")
	ErrSubredditNameAlreadyRegistered error = errors.New("subreddit: name already registered")
	ErrSubredditNameEmpty             error = errors.New("subreddit: empty name")
	ErrSubredditNotFound              error = errors.New("subreddit: not found")
)
