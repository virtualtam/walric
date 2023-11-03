package sqlite3

import (
	"time"

	"github.com/virtualtam/walric/pkg/history"
)

type DBEntry struct {
	ID   int       `db:"id"`
	Date time.Time `db:"date"`

	SubmissionID int `db:"submission_id"`
}

func newDBEntry(entry *history.Entry) *DBEntry {
	return &DBEntry{
		ID:           entry.ID,
		Date:         entry.Date,
		SubmissionID: entry.Submission.ID,
	}
}
