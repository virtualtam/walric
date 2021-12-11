package redwall

type SubmissionRepository interface {
	ByID(id int) (*Submission, error)
	ByPostID(postID string) (*Submission, error)
}