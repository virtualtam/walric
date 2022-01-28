package submission

import (
	"fmt"
	"time"

	"github.com/virtualtam/walric/subreddit"
)

// Submission represents the metadata for a Reddit post with an image
// attachment, and the metadata for the corresponding local file.
type Submission struct {
	ID int `db:"id"`

	Subreddit   *subreddit.Subreddit `db:"-"`
	SubredditID int                  `db:"subreddit_id"`

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

// PermalinkURL returns the Reddit permalink for this submission's post.
func (s *Submission) PermalinkURL() string {
	return fmt.Sprintf("https://reddit.com%s", s.Permalink)
}

// User returns the Reddit-formatted username for this submission's author.
func (s *Submission) User() string {
	return fmt.Sprintf("u/%s", s.Author)
}
