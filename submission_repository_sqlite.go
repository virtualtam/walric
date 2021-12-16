package redwall

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var _ SubmissionRepository = &SubmissionRepositorySQLite{}

type SubmissionRepositorySQLite struct {
	db               *sqlx.DB
	subredditService *SubredditService
}

func (r *SubmissionRepositorySQLite) ByID(id int) (*Submission, error) {
	submission := &Submission{}

	err := r.db.QueryRowx(`
SELECT
  id,
  author,
  created_utc,
  image_filename,
  image_height_px,
  image_width_px,
  post_id,
  subreddit_id,
  title,
  url
FROM submissions WHERE id=?`,
		id,
	).StructScan(submission)

	if err != nil {
		return &Submission{}, err
	}

	subreddit, err := r.subredditService.ByID(submission.SubredditID)
	if err != nil {
		return &Submission{}, err
	}

	submission.Subreddit = subreddit

	return submission, nil
}

func (r *SubmissionRepositorySQLite) ByMinResolution(minResolution *Resolution) ([]*Submission, error) {
	rows, err := r.db.Queryx(`
SELECT
  sm.id,
  sm.author,
  sm.created_utc,
  sm.image_filename,
  sm.image_height_px,
  sm.image_width_px,
  sm.post_id,
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

		subreddit, err := r.subredditService.ByID(submission.SubredditID)
		if err != nil {
			return []*Submission{}, err
		}

		submission.Subreddit = subreddit

		submissions = append(submissions, submission)
	}

	return submissions, nil
}

func (r *SubmissionRepositorySQLite) ByPostID(postID string) (*Submission, error) {
	submission := &Submission{}

	err := r.db.QueryRowx(`
SELECT
  id,
  author,
  created_utc,
  image_filename,
  image_height_px,
  image_width_px,
  post_id,
  subreddit_id,
  title,
  url
FROM submissions WHERE post_id=?`,
		postID,
	).StructScan(submission)

	if err != nil {
		return &Submission{}, err
	}

	subreddit, err := r.subredditService.ByID(submission.SubredditID)
	if err != nil {
		return &Submission{}, err
	}

	submission.Subreddit = subreddit

	return submission, nil
}

func (r *SubmissionRepositorySQLite) ByTitle(searchText string) ([]*Submission, error) {
	searchPattern := fmt.Sprintf("%%%s%%", searchText)

	rows, err := r.db.Queryx(`
SELECT
  id,
  author,
  created_utc,
  image_filename,
  image_height_px,
  image_width_px,
  post_id,
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

		subreddit, err := r.subredditService.ByID(submission.SubredditID)
		if err != nil {
			return []*Submission{}, err
		}

		submission.Subreddit = subreddit

		submissions = append(submissions, submission)
	}

	return submissions, nil
}

func (r *SubmissionRepositorySQLite) Random(minResolution *Resolution) (*Submission, error) {
	submission := &Submission{}

	err := r.db.QueryRowx(`
SELECT
  id,
  author,
  created_utc,
  image_filename,
  image_height_px,
  image_width_px,
  post_id,
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

	if err != nil {
		return &Submission{}, err
	}

	subreddit, err := r.subredditService.ByID(submission.SubredditID)
	if err != nil {
		return &Submission{}, err
	}

	submission.Subreddit = subreddit

	return submission, nil
}

func NewSubmissionRepositorySQLite(db *sqlx.DB, subredditService *SubredditService) *SubmissionRepositorySQLite {
	return &SubmissionRepositorySQLite{
		db:               db,
		subredditService: subredditService,
	}
}
