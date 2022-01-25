package submission

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/virtualtam/walric/monitor"
)

var _ Repository = &RepositorySQLite{}

type RepositorySQLite struct {
	db *sqlx.DB
}

func (r *RepositorySQLite) ByID(id int) (*Submission, error) {
	submission := &Submission{}

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
	).StructScan(submission)

	if errors.Is(err, sql.ErrNoRows) {
		return &Submission{}, ErrNotFound
	}
	if err != nil {
		return &Submission{}, err
	}

	return submission, nil
}

func (r *RepositorySQLite) ByMinResolution(minResolution *monitor.Resolution) ([]*Submission, error) {
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
		return []*Submission{}, err
	}

	submissions := []*Submission{}

	for rows.Next() {
		submission := &Submission{}

		if err := rows.StructScan(submission); err != nil {
			return []*Submission{}, err
		}

		submissions = append(submissions, submission)
	}

	return submissions, nil
}

func (r *RepositorySQLite) ByPostID(postID string) (*Submission, error) {
	submission := &Submission{}

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
	).StructScan(submission)

	if errors.Is(err, sql.ErrNoRows) {
		return &Submission{}, ErrNotFound
	}
	if err != nil {
		return &Submission{}, err
	}

	return submission, nil
}

func (r *RepositorySQLite) Search(text string) ([]*Submission, error) {
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
		return []*Submission{}, err
	}

	submissions := []*Submission{}

	for rows.Next() {
		submission := &Submission{}

		if err := rows.StructScan(submission); err != nil {
			return []*Submission{}, err
		}

		submissions = append(submissions, submission)
	}

	return submissions, nil
}

func (r *RepositorySQLite) Random(minResolution *monitor.Resolution) (*Submission, error) {
	submission := &Submission{}

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
	).StructScan(submission)

	if errors.Is(err, sql.ErrNoRows) {
		return &Submission{}, ErrNotFound
	}
	if err != nil {
		return &Submission{}, err
	}

	return submission, nil
}

func (r *RepositorySQLite) Create(submission *Submission) error {
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
		submission,
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
