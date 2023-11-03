package history

import (
	"time"

	"github.com/virtualtam/walric/pkg/submission"
)

// Entry represents a Submission that was selected as being suitable as a
// wallpaper for a given monitor setup.
type Entry struct {
	ID   int       `db:"id"`
	Date time.Time `db:"date"`

	SubmissionID int                    `db:"submission_id"`
	Submission   *submission.Submission `db:"-"`
}

// NewEntry initializes and returns a new history Entry.
func NewEntry(sub *submission.Submission) (*Entry, error) {
	if sub.ID < 1 {
		return &Entry{}, ErrSubmissionIDNegativeOrZero
	}

	now := time.Now().UTC()

	return &Entry{
		Date:         now,
		Submission:   sub,
		SubmissionID: sub.ID,
	}, nil
}
