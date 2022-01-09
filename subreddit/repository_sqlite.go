package subreddit

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var _ Repository = &RepositorySQLite{}

type RepositorySQLite struct {
	db *sqlx.DB
}

func (r *RepositorySQLite) All() ([]*Subreddit, error) {
	rows, err := r.db.Queryx("SELECT id, name from subreddits ORDER BY name COLLATE NOCASE")

	if err != nil {
		return []*Subreddit{}, err
	}

	subreddits := []*Subreddit{}

	for rows.Next() {
		subreddit := &Subreddit{}

		if err := rows.StructScan(subreddit); err != nil {
			return []*Subreddit{}, err
		}

		subreddits = append(subreddits, subreddit)
	}

	return subreddits, nil
}

func (r *RepositorySQLite) ByID(id int) (*Subreddit, error) {
	subreddit := &Subreddit{}

	err := r.db.QueryRowx("SELECT id, name FROM subreddits WHERE id=?", id).StructScan(subreddit)
	if errors.Is(err, sql.ErrNoRows) {
		return &Subreddit{}, ErrNotFound
	}
	if err != nil {
		return &Subreddit{}, err
	}

	return subreddit, nil
}

func (r *RepositorySQLite) ByName(name string) (*Subreddit, error) {
	subreddit := &Subreddit{}

	err := r.db.QueryRowx("SELECT id, name FROM subreddits WHERE name=?", name).StructScan(subreddit)
	if errors.Is(err, sql.ErrNoRows) {
		return &Subreddit{}, ErrNotFound
	}
	if err != nil {
		return &Subreddit{}, err
	}

	return subreddit, nil
}

func (r *RepositorySQLite) Stats() ([]SubredditStats, error) {
	rows, err := r.db.Queryx(`SELECT sr.name as name, COUNT(sm.post_id) as submissions
FROM subreddits AS sr
LEFT JOIN submissions AS sm ON sr.id = sm.subreddit_id
GROUP BY sr.name
ORDER BY sr.name COLLATE NOCASE
`)

	if err != nil {
		return []SubredditStats{}, err
	}

	subredditStats := []SubredditStats{}

	for rows.Next() {
		stats := SubredditStats{}

		if err := rows.StructScan(&stats); err != nil {
			return []SubredditStats{}, err
		}

		subredditStats = append(subredditStats, stats)
	}

	return subredditStats, nil
}

func NewRepositorySQLite(db *sqlx.DB) *RepositorySQLite {
	return &RepositorySQLite{
		db: db,
	}
}
