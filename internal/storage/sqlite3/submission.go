package sqlite3

import (
	"time"

	"github.com/virtualtam/walric/pkg/submission"
)

type DBSubmission struct {
	ID int `db:"id"`

	SubredditID int `db:"subreddit_id"`

	// Reddit post metadata
	Author    string    `db:"author"`
	Permalink string    `db:"permalink"`
	PostID    string    `db:"post_id"`
	PostedAt  time.Time `db:"created_utc"`
	Score     int       `db:"score"`
	Title     string    `db:"title"`

	// Attached image metadata
	ImageDomain string `db:"domain"`
	ImageURL    string `db:"url"`
	ImageNSFW   bool   `db:"over_18"`

	// Local image metadata
	ImageFilename string `db:"image_filename"`
	ImageHeightPx int    `db:"image_height_px"`
	ImageWidthPx  int    `db:"image_width_px"`
}

func newDBSubmission(sub *submission.Submission) *DBSubmission {
	return &DBSubmission{
		ID:            sub.ID,
		SubredditID:   sub.Subreddit.ID,
		Author:        sub.Author,
		Permalink:     sub.Permalink,
		PostID:        sub.PostID,
		PostedAt:      sub.PostedAt,
		Score:         sub.Score,
		Title:         sub.Title,
		ImageDomain:   sub.ImageDomain,
		ImageURL:      sub.ImageURL,
		ImageNSFW:     sub.ImageNSFW,
		ImageFilename: sub.ImageFilename,
		ImageHeightPx: sub.ImageHeightPx,
		ImageWidthPx:  sub.ImageWidthPx,
	}
}

func (s *DBSubmission) AsSubmission(sr *submission.Subreddit) *submission.Submission {
	return &submission.Submission{
		ID:            s.ID,
		Subreddit:     sr,
		Author:        s.Author,
		Permalink:     s.Permalink,
		PostID:        s.PostID,
		PostedAt:      s.PostedAt,
		Score:         s.Score,
		Title:         s.Title,
		ImageDomain:   s.ImageDomain,
		ImageURL:      s.ImageURL,
		ImageNSFW:     s.ImageNSFW,
		ImageFilename: s.ImageFilename,
		ImageHeightPx: s.ImageHeightPx,
		ImageWidthPx:  s.ImageWidthPx,
	}
}
