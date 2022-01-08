package submission

import "errors"

var (
	ErrIDInvalid         error = errors.New("submission: invalid ID")
	ErrNotFound          error = errors.New("submission: not found")
	ErrPostIDEmpty       error = errors.New("submission: empty post ID")
	ErrResolutionInvalid error = errors.New("submission: invalid resolution")
	ErrSearchTextEmpty   error = errors.New("submission: empty search text")
)
