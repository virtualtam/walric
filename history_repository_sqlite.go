package redwall

import "github.com/jmoiron/sqlx"

var _ HistoryRepository = &HistoryRepositorySQLite{}

type HistoryRepositorySQLite struct {
	db                *sqlx.DB
	submissionService *SubmissionService
}

func (r *HistoryRepositorySQLite) All() ([]HistoryEntry, error) {
	rows, err := r.db.Queryx("SELECT date, submission_id FROM history ORDER BY date")
	if err != nil {
		return []HistoryEntry{}, err
	}

	history := []HistoryEntry{}

	for rows.Next() {
		entry := HistoryEntry{}

		if err := rows.StructScan(&entry); err != nil {
			return []HistoryEntry{}, err
		}

		submission, err := r.submissionService.ByID(entry.SubmissionID)
		if err != nil {
			return []HistoryEntry{}, err
		}

		entry.Submission = submission

		history = append(history, entry)
	}

	return history, nil
}

func NewHistoryRepositorySQLite(db *sqlx.DB, submissionService *SubmissionService) *HistoryRepositorySQLite {
	return &HistoryRepositorySQLite{
		db:                db,
		submissionService: submissionService,
	}
}
