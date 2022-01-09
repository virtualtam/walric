package history

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var _ Repository = &RepositorySQLite{}

type RepositorySQLite struct {
	db *sqlx.DB
}

func (r *RepositorySQLite) All() ([]*Entry, error) {
	rows, err := r.db.Queryx("SELECT date, submission_id FROM history ORDER BY date")
	if err != nil {
		return []*Entry{}, err
	}

	history := []*Entry{}

	for rows.Next() {
		entry := &Entry{}

		if err := rows.StructScan(&entry); err != nil {
			return []*Entry{}, err
		}

		history = append(history, entry)
	}

	return history, nil
}

func (r *RepositorySQLite) Current() (*Entry, error) {
	entry := &Entry{}

	err := r.db.QueryRowx("SELECT date, submission_id FROM history ORDER BY date desc LIMIT 1").StructScan(entry)
	if errors.Is(err, sql.ErrNoRows) {
		return &Entry{}, ErrNotFound
	}
	if err != nil {
		return &Entry{}, err
	}

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

func NewRepositorySQLite(db *sqlx.DB) *RepositorySQLite {
	return &RepositorySQLite{
		db: db,
	}
}
