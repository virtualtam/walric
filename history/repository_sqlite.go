package history

import (
	"github.com/jmoiron/sqlx"
	"github.com/virtualtam/redwall2/submission"
)

var _ Repository = &RepositorySQLite{}

type RepositorySQLite struct {
	db                *sqlx.DB
	submissionService *submission.Service
}

func (r *RepositorySQLite) All() ([]Entry, error) {
	rows, err := r.db.Queryx("SELECT date, submission_id FROM history ORDER BY date")
	if err != nil {
		return []Entry{}, err
	}

	history := []Entry{}

	for rows.Next() {
		entry := Entry{}

		if err := rows.StructScan(&entry); err != nil {
			return []Entry{}, err
		}

		submission, err := r.submissionService.ByID(entry.SubmissionID)
		if err != nil {
			return []Entry{}, err
		}

		entry.Submission = submission

		history = append(history, entry)
	}

	return history, nil
}

func (r *RepositorySQLite) Current() (*Entry, error) {
	entry := &Entry{}

	if err := r.db.QueryRowx("SELECT date, submission_id FROM history ORDER BY date desc LIMIT 1").StructScan(entry); err != nil {
		return &Entry{}, err
	}

	submission, err := r.submissionService.ByID(entry.SubmissionID)
	if err != nil {
		return &Entry{}, err
	}

	entry.Submission = submission

	return entry, nil
}

func (r *RepositorySQLite) Create(entry *Entry) error {
	_, err := r.db.NamedExec(`
INSERT INTO history(date, submission_id)
VALUES (:date, :submission_id)`,
		entry,
	)
	if err != nil {
		return err
	}

	return nil
}

func NewRepositorySQLite(db *sqlx.DB, submissionService *submission.Service) *RepositorySQLite {
	return &RepositorySQLite{
		db:                db,
		submissionService: submissionService,
	}
}
