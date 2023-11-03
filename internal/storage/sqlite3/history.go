package sqlite3

import (
	"time"

	"github.com/virtualtam/walric/pkg/history"
)

type Entry struct {
	ID   int       `db:"id"`
	Date time.Time `db:"date"`

	SubmissionID int `db:"submission_id"`
}

func newEntry(entry *history.Entry) *Entry {
	return &Entry{
		ID:           entry.ID,
		Date:         entry.Date,
		SubmissionID: entry.Submission.ID,
	}
}
