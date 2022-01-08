package history

import "errors"

var (
	ErrNotFound                   error = errors.New("not found")
	ErrSubmissionIDNegativeOrZero error = errors.New("Submission ID is negative or zero")
)
