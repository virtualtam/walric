package sqlite3

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/virtualtam/walric/pkg/history"
	"github.com/virtualtam/walric/pkg/monitor"
	"github.com/virtualtam/walric/pkg/submission"
)

var _ history.Repository = &Repository{}
var _ submission.Repository = &Repository{}

// Repository provides a SQLite3 database persistence layer for
// Subreddits.
type Repository struct {
	db *sqlx.DB
}

// NewRepository initializes and returns a SQLite3 repository to persist
// and manage Subreddits.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) HistoryGetAll() ([]*history.Entry, error) {
	rows, err := r.db.Queryx("SELECT date, submission_id FROM history ORDER BY date")
	if err != nil {
		return []*history.Entry{}, err
	}

	entries := []*history.Entry{}

	for rows.Next() {
		dbEntry := &DBEntry{}

		if err := rows.StructScan(&dbEntry); err != nil {
			return []*history.Entry{}, err
		}

		sub, err := r.SubmissionGetByID(dbEntry.SubmissionID)
		if err != nil {
			return []*history.Entry{}, err
		}

		entry := &history.Entry{
			ID:         dbEntry.ID,
			Date:       dbEntry.Date,
			Submission: sub,
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *Repository) HistoryGetCurrent() (*history.Entry, error) {
	dbEntry := &DBEntry{}

	err := r.db.QueryRowx("SELECT date, submission_id FROM history ORDER BY date desc LIMIT 1").StructScan(dbEntry)
	if errors.Is(err, sql.ErrNoRows) {
		return &history.Entry{}, history.ErrNotFound
	}
	if err != nil {
		return &history.Entry{}, err
	}

	sub, err := r.SubmissionGetByID(dbEntry.SubmissionID)
	if err != nil {
		return &history.Entry{}, err
	}

	entry := &history.Entry{
		ID:         dbEntry.ID,
		Date:       dbEntry.Date,
		Submission: sub,
	}

	return entry, nil
}

func (r *Repository) HistoryCreate(entry *history.Entry) error {
	dbEntry := newDBEntry(entry)

	_, err := r.db.NamedExec(`
INSERT INTO history(date, submission_id)
VALUES (:date, :submission_id)`,
		dbEntry,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) submissionGetQuery(query string, queryParams ...any) (*submission.Submission, error) {
	dbSubmission := &DBSubmission{}

	err := r.db.QueryRowx(query, queryParams...).StructScan(dbSubmission)

	if errors.Is(err, sql.ErrNoRows) {
		return &submission.Submission{}, submission.ErrSubmissionNotFound
	}
	if err != nil {
		return &submission.Submission{}, err
	}

	sr, err := r.SubredditGetByID(dbSubmission.SubredditID)
	if err != nil {
		return &submission.Submission{}, err
	}

	s := dbSubmission.AsSubmission(sr)

	return s, nil
}

func (r *Repository) submissionGetManyQuery(query string, queryParams ...any) ([]*submission.Submission, error) {
	rows, err := r.db.Queryx(query, queryParams...)

	if err != nil {
		return []*submission.Submission{}, err
	}

	submissions := []*submission.Submission{}

	for rows.Next() {
		dbSubmission := &DBSubmission{}

		if err := rows.StructScan(dbSubmission); err != nil {
			return []*submission.Submission{}, err
		}

		sr, err := r.SubredditGetByID(dbSubmission.SubredditID)
		if err != nil {
			return []*submission.Submission{}, err
		}

		s := dbSubmission.AsSubmission(sr)
		submissions = append(submissions, s)
	}

	return submissions, nil
}

func (r *Repository) SubmissionGetByID(id int) (*submission.Submission, error) {
	return r.submissionGetQuery(`
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
	)
}

func (r *Repository) SubmissionGetByMinResolution(minResolution *monitor.Resolution) ([]*submission.Submission, error) {
	return r.submissionGetManyQuery(`
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
}

func (r *Repository) SubmissionGetByPostID(postID string) (*submission.Submission, error) {
	return r.submissionGetQuery(`
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
	)
}

func (r *Repository) SubmissionIsPostIDRegistered(postID string) (bool, error) {
	var registered int64

	err := r.db.QueryRowx(
		"SELECT id FROM submissions WHERE post_id=?",
		postID,
	).Scan(&registered)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *Repository) SubmissionSearch(text string) ([]*submission.Submission, error) {
	searchPattern := fmt.Sprintf("%%%s%%", text)

	return r.submissionGetManyQuery(`
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
}

func (r *Repository) SubmissionGetRandom(minResolution *monitor.Resolution) (*submission.Submission, error) {
	return r.submissionGetQuery(`
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
	)
}

func (r *Repository) SubmissionCreate(s *submission.Submission) error {
	dbSubmission := newDBSubmission(s)

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
		dbSubmission,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) SubredditCreate(s *submission.Subreddit) error {
	_, err := r.db.NamedExec("INSERT INTO subreddits(name) VALUES(:name)", s)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) SubredditGetAll() ([]*submission.Subreddit, error) {
	rows, err := r.db.Queryx("SELECT id, name from subreddits ORDER BY name COLLATE NOCASE")

	if err != nil {
		return []*submission.Subreddit{}, err
	}

	subreddits := []*submission.Subreddit{}

	for rows.Next() {
		s := &submission.Subreddit{}

		if err := rows.StructScan(s); err != nil {
			return []*submission.Subreddit{}, err
		}

		subreddits = append(subreddits, s)
	}

	return subreddits, nil
}

func (r *Repository) SubredditGetByID(id int) (*submission.Subreddit, error) {
	s := &submission.Subreddit{}

	err := r.db.QueryRowx("SELECT id, name FROM subreddits WHERE id=?", id).StructScan(s)
	if errors.Is(err, sql.ErrNoRows) {
		return &submission.Subreddit{}, submission.ErrSubredditNotFound
	}
	if err != nil {
		return &submission.Subreddit{}, err
	}

	return s, nil
}

func (r *Repository) SubredditGetByName(name string) (*submission.Subreddit, error) {
	s := &submission.Subreddit{}

	err := r.db.QueryRowx("SELECT id, name FROM subreddits WHERE name=?", name).StructScan(s)
	if errors.Is(err, sql.ErrNoRows) {
		return &submission.Subreddit{}, submission.ErrSubredditNotFound
	}
	if err != nil {
		return &submission.Subreddit{}, err
	}

	return s, nil
}

func (r *Repository) SubredditIsNameRegistered(name string) (bool, error) {
	var registered int64

	err := r.db.QueryRowx("SELECT id FROM subreddits WHERE name=?", name).Scan(&registered)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *Repository) SubredditGetStats() ([]submission.SubredditStats, error) {
	rows, err := r.db.Queryx(`
SELECT sr.name as name, COUNT(sm.post_id) as submissions
FROM subreddits AS sr
LEFT JOIN submissions AS sm ON sr.id = sm.subreddit_id
GROUP BY sr.name
ORDER BY sr.name COLLATE NOCASE
`)

	if err != nil {
		return []submission.SubredditStats{}, err
	}

	subredditStats := []submission.SubredditStats{}

	for rows.Next() {
		stats := submission.SubredditStats{}

		if err := rows.StructScan(&stats); err != nil {
			return []submission.SubredditStats{}, err
		}

		subredditStats = append(subredditStats, stats)
	}

	return subredditStats, nil
}
