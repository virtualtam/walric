package redwall

type SubmissionRepository interface {
	ByPostID(postID string) (*Submission, error)
}
