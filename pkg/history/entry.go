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
