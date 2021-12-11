package redwall

import (
	"github.com/jmoiron/sqlx"
)

var _ SubmissionRepository = &SubmissionRepositorySQLite{}

type SubmissionRepositorySQLite struct {
	db               *sqlx.DB
	subredditService *SubredditService
}

func (r *SubmissionRepositorySQLite) ByPostID(postID string) (*Submission, error) {
	submission := &Submission{}

	err := r.db.QueryRowx(`
SELECT
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

func NewSubmissionRepositorySQLite(db *sqlx.DB, subredditService *SubredditService) *SubmissionRepositorySQLite {
	return &SubmissionRepositorySQLite{
		db:               db,
		subredditService: subredditService,
	}
}
