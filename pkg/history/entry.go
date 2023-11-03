package history

import (
	"time"

	"github.com/virtualtam/walric/pkg/submission"
)

// Entry represents a Submission that was selected as being suitable as a
// wallpaper for a given monitor setup.
type Entry struct {
	ID   int
	Date time.Time

	Submission *submission.Submission
}

// NewEntry initializes and returns a new history Entry.
func NewEntry(sub *submission.Submission) (*Entry, error) {
	if sub.ID < 1 {
		return &Entry{}, ErrSubmissionIDNegativeOrZero
	}

	now := time.Now().UTC()

	return &Entry{
		Date:       now,
		Submission: sub,
	}, nil
}
