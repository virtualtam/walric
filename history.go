package redwall

import "time"

type HistoryEntry struct {
	ID   int       `db:"id"`
	Date time.Time `db:"date"`

	SubmissionID int         `db:"submission_id"`
	Submission   *Submission `db:"-"`
}
