package history

import (
	"time"

	"github.com/virtualtam/redwall2/submission"
)

type Entry struct {
	ID   int       `db:"id"`
	Date time.Time `db:"date"`

	SubmissionID int                    `db:"submission_id"`
	Submission   *submission.Submission `db:"-"`
}
