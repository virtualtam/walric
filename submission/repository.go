package submission

import "github.com/virtualtam/redwall2/monitor"

type Repository interface {
	ByID(id int) (*Submission, error)
	ByPostID(postID string) (*Submission, error)

	Search(text string) ([]*Submission, error)

	ByMinResolution(minResolution *monitor.Resolution) ([]*Submission, error)
	Random(minResolution *monitor.Resolution) (*Submission, error)

	Create(submission *Submission) error
}
