package redwall

import (
	"fmt"
	"time"
)

type Submission struct {
	ID int `db:"id"`

	Subreddit   *Subreddit `db:"-"`
	SubredditID int        `db:"subreddit_id"`

	PostID   string    `db:"post_id"`
	PostedAt time.Time `db:"created_utc"`

	Author        string `db:"author"`
	Title         string `db:"title"`
	ImageURL      string `db:"url"`
	ImageFilename string `db:"image_filename"`
	ImageHeightPx int    `db:"image_height_px"`
	ImageWidthPx  int    `db:"image_width_px"`
}

func (s *Submission) PostURL() string {
	if s.Subreddit == nil {
		return ""
	}

	return fmt.Sprintf("https://reddit.com/r/%s/comments/%s/", s.Subreddit.Name, s.PostID)
}

func (s *Submission) User() string {
	return fmt.Sprintf("u/%s", s.Author)
}
