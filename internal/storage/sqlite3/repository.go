package sqlite3

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/virtualtam/walric/pkg/history"
	"github.com/virtualtam/walric/pkg/monitor"
	"github.com/virtualtam/walric/pkg/submission"
	"github.com/virtualtam/walric/pkg/subreddit"
)

var _ history.Repository = &RepositorySQLite{}
var _ submission.Repository = &RepositorySQLite{}
var _ subreddit.Repository = &RepositorySQLite{}

// RepositorySQLite provides a SQLite3 database persistence layer for
// Subreddits.
type RepositorySQLite struct {
	db *sqlx.DB
}

func (r *RepositorySQLite) HistoryGetAll() ([]*history.Entry, error) {
	rows, err := r.db.Queryx("SELECT date, submission_id FROM history ORDER BY date")
	if err != nil {
		return []*history.Entry{}, err
	}

	entries := []*history.Entry{}

	for rows.Next() {
		entry := &history.Entry{}

		if err := rows.StructScan(&entry); err != nil {
			return []*history.Entry{}, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *RepositorySQLite) HistoryGetCurrent() (*history.Entry, error) {
	entry := &history.Entry{}

	err := r.db.QueryRowx("SELECT date, submission_id FROM history ORDER BY date desc LIMIT 1").StructScan(entry)
	if errors.Is(err, sql.ErrNoRows) {
		return &history.Entry{}, history.ErrNotFound
	}
	if err != nil {
		return &history.Entry{}, err
	}

	return entry, nil
}

func (r *RepositorySQLite) HistoryCreate(entry *history.Entry) error {
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

func (r *RepositorySQLite) SubmissionGetByID(id int) (*submission.Submission, error) {
	s := &submission.Submission{}

	err := r.db.QueryRowx(`
SELECT
  id,
  author,
  created_utc,
  domain,
  image_filename,
  image_height_px,
  image_width_px,
  over_18,
  permalink,
  post_id,
  score,
  subreddit_id,
  title,
  url
FROM submissions WHERE id=?`,
		id,
	).StructScan(s)

	if errors.Is(err, sql.ErrNoRows) {
		return &submission.Submission{}, submission.ErrNotFound
	}
	if err != nil {
		return &submission.Submission{}, err
	}

	return s, nil
}

func (r *RepositorySQLite) SubmissionGetByMinResolution(minResolution *monitor.Resolution) ([]*submission.Submission, error) {
	rows, err := r.db.Queryx(`
SELECT
  sm.id,
  sm.author,
  sm.created_utc,
  sm.domain,
  sm.image_filename,
  sm.image_height_px,
  sm.image_width_px,
  sm.over_18,
  sm.permalink,
  sm.post_id,
  sm.score,
  sm.subreddit_id,
  sm.title,
  sm.url
FROM submissions sm
LEFT JOIN subreddits sub ON sm.subreddit_id=sub.id
WHERE image_height_px >= ?
AND   image_width_px  >= ?
ORDER BY sub.name COLLATE NOCASE, sm.created_utc
`,
		minResolution.HeightPx,
		minResolution.WidthPx,
	)

	if err != nil {
		return []*submission.Submission{}, err
	}

	submissions := []*submission.Submission{}

	for rows.Next() {
		s := &submission.Submission{}

		if err := rows.StructScan(s); err != nil {
			return []*submission.Submission{}, err
		}

		submissions = append(submissions, s)
	}

	return submissions, nil
}

func (r *RepositorySQLite) SubmissionGetByPostID(postID string) (*submission.Submission, error) {
	s := &submission.Submission{}

	err := r.db.QueryRowx(`
SELECT
  id,
  author,
  created_utc,
  domain,
  image_filename,
  image_height_px,
  image_width_px,
  over_18,
  permalink,
  post_id,
  score,
  subreddit_id,
  title,
  url
FROM submissions WHERE post_id=?`,
		postID,
	).StructScan(s)

	if errors.Is(err, sql.ErrNoRows) {
		return &submission.Submission{}, submission.ErrNotFound
	}
	if err != nil {
		return &submission.Submission{}, err
	}

	return s, nil
}

func (r *RepositorySQLite) SubmissionIsPostIDRegistered(postID string) (bool, error) {
	s := &submission.Submission{}

	err := r.db.QueryRowx(
		"SELECT id FROM submissions WHERE post_id=?",
		postID,
	).StructScan(s)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *RepositorySQLite) SubmissionSearch(text string) ([]*submission.Submission, error) {
	searchPattern := fmt.Sprintf("%%%s%%", text)

	rows, err := r.db.Queryx(`
SELECT
  id,
  author,
  created_utc,
  domain,
  image_filename,
  image_height_px,
  image_width_px,
  over_18,
  permalink,
  post_id,
  score,
  subreddit_id,
  title,
  url
FROM submissions
WHERE title LIKE ? COLLATE NOCASE
ORDER BY created_utc
`,
		searchPattern,
	)

	if err != nil {
		return []*submission.Submission{}, err
	}

	submissions := []*submission.Submission{}

	for rows.Next() {
		s := &submission.Submission{}

		if err := rows.StructScan(s); err != nil {
			return []*submission.Submission{}, err
		}

		submissions = append(submissions, s)
	}

	return submissions, nil
}

func (r *RepositorySQLite) SubmissionGetRandom(minResolution *monitor.Resolution) (*submission.Submission, error) {
	s := &submission.Submission{}

	err := r.db.QueryRowx(`
SELECT
  id,
  author,
  created_utc,
  domain,
  image_filename,
  image_height_px,
  image_width_px,
  over_18,
  permalink,
  post_id,
  score,
  subreddit_id,
  title,
  url
FROM submissions
WHERE image_height_px >= ?
AND   image_width_px  >= ?
AND   id NOT IN (SELECT submission_id from history)
ORDER BY RANDOM() LIMIT 1
`,
		minResolution.HeightPx,
		minResolution.WidthPx,
	).StructScan(s)

	if errors.Is(err, sql.ErrNoRows) {
		return &submission.Submission{}, submission.ErrNotFound
	}
	if err != nil {
		return &submission.Submission{}, err
	}

	return s, nil
}

func (r *RepositorySQLite) SubmissionCreate(s *submission.Submission) error {
	_, err := r.db.NamedExec(`
INSERT INTO submissions(
	subreddit_id,
	author,
	permalink,
	post_id,
	created_utc,
	score,
	title,
	domain,
	url,
	over_18,
	image_filename,
	image_height_px,
	image_width_px
)
VALUES (
	:subreddit_id,
	:author,
	:permalink,
	:post_id,
	:created_utc,
	:score,
	:title,
	:domain,
	:url,
	:over_18,
	:image_filename,
	:image_height_px,
	:image_width_px
)`,
		s,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *RepositorySQLite) SubredditCreate(s *subreddit.Subreddit) error {
	_, err := r.db.NamedExec("INSERT INTO subreddits(name) VALUES(:name)", s)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositorySQLite) SubredditGetAll() ([]*subreddit.Subreddit, error) {
	rows, err := r.db.Queryx("SELECT id, name from subreddits ORDER BY name COLLATE NOCASE")

	if err != nil {
		return []*subreddit.Subreddit{}, err
	}

	subreddits := []*subreddit.Subreddit{}

	for rows.Next() {
		s := &subreddit.Subreddit{}

		if err := rows.StructScan(s); err != nil {
			return []*subreddit.Subreddit{}, err
		}

		subreddits = append(subreddits, s)
	}

	return subreddits, nil
}

func (r *RepositorySQLite) SubredditGetByID(id int) (*subreddit.Subreddit, error) {
	s := &subreddit.Subreddit{}

	err := r.db.QueryRowx("SELECT id, name FROM subreddits WHERE id=?", id).StructScan(s)
	if errors.Is(err, sql.ErrNoRows) {
		return &subreddit.Subreddit{}, subreddit.ErrNotFound
	}
	if err != nil {
		return &subreddit.Subreddit{}, err
	}

	return s, nil
}

func (r *RepositorySQLite) SubredditGetByName(name string) (*subreddit.Subreddit, error) {
	s := &subreddit.Subreddit{}

	err := r.db.QueryRowx("SELECT id, name FROM subreddits WHERE name=?", name).StructScan(s)
	if errors.Is(err, sql.ErrNoRows) {
		return &subreddit.Subreddit{}, subreddit.ErrNotFound
	}
	if err != nil {
		return &subreddit.Subreddit{}, err
	}

	return s, nil
}

func (r *RepositorySQLite) SubredditIsNameRegistered(name string) (bool, error) {
	s := &subreddit.Subreddit{}

	err := r.db.QueryRowx("SELECT id FROM subreddits WHERE name=?", name).StructScan(s)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *RepositorySQLite) SubredditGetStats() ([]subreddit.SubredditStats, error) {
	rows, err := r.db.Queryx(`SELECT sr.name as name, COUNT(sm.post_id) as submissions
FROM subreddits AS sr
LEFT JOIN submissions AS sm ON sr.id = sm.subreddit_id
GROUP BY sr.name
ORDER BY sr.name COLLATE NOCASE
`)

	if err != nil {
		return []subreddit.SubredditStats{}, err
	}

	subredditStats := []subreddit.SubredditStats{}

	for rows.Next() {
		stats := subreddit.SubredditStats{}

		if err := rows.StructScan(&stats); err != nil {
			return []subreddit.SubredditStats{}, err
		}

		subredditStats = append(subredditStats, stats)
	}

	return subredditStats, nil
}

// NewRepositorySQLite initializes and returns a SQLite3 repository to persist
// and manage Subreddits.
func NewRepositorySQLite(db *sqlx.DB) *RepositorySQLite {
	return &RepositorySQLite{
		db: db,
	}
}
