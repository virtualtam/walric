package redwall

type SubmissionRepository interface {
	ByID(id int) (*Submission, error)
	ByPostID(postID string) (*Submission, error)
	ByMinResolution(minResolution *Resolution) ([]*Submission, error)
	ByTitle(text string) ([]*Submission, error)
	Random(minResolution *Resolution) (*Submission, error)
}
