package redwall

import "github.com/jmoiron/sqlx"

type SubredditRepositorySQLite struct {
	db *sqlx.DB
}

func (r *SubredditRepositorySQLite) All() ([]Subreddit, error) {
	rows, err := r.db.Queryx("SELECT id, name from subreddits ORDER BY name COLLATE NOCASE")

	if err != nil {
		return []Subreddit{}, err
	}

	subreddits := []Subreddit{}

	for rows.Next() {
		subreddit := Subreddit{}

		if err := rows.StructScan(&subreddit); err != nil {
			return []Subreddit{}, err
		}

		subreddits = append(subreddits, subreddit)
	}

	return subreddits, nil
}

func (r *SubredditRepositorySQLite) ByID(id int) (*Subreddit, error) {
	subreddit := &Subreddit{}

	err := r.db.QueryRowx("SELECT id, name FROM subreddits WHERE id=?", id).StructScan(subreddit)
	if err != nil {
		return &Subreddit{}, err
	}

	return subreddit, nil
}

func (r *SubredditRepositorySQLite) ByName(name string) (*Subreddit, error) {
	subreddit := &Subreddit{}

	err := r.db.QueryRowx("SELECT id, name FROM subreddits WHERE name=?", name).StructScan(subreddit)
	if err != nil {
		return &Subreddit{}, err
	}

	return subreddit, nil
}

func (r *SubredditRepositorySQLite) Stats() ([]SubredditStats, error) {
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

func NewSubredditRepositorySQLite(db *sqlx.DB) *SubredditRepositorySQLite {
	return &SubredditRepositorySQLite{
		db: db,
	}
}
