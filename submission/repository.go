package submission

import "github.com/virtualtam/redwall2/monitor"

type Repository interface {
	ByID(id int) (*Submission, error)
	ByPostID(postID string) (*Submission, error)
	ByMinResolution(minResolution *monitor.Resolution) ([]*Submission, error)
	ByTitle(text string) ([]*Submission, error)
	Random(minResolution *monitor.Resolution) (*Submission, error)
}
