package submission

import "github.com/virtualtam/walric/pkg/monitor"

// Repository defines the basic operations available to access and persist
// Reddit Submissions.
type Repository interface {
	// SubmissionGetByID returns the Submission for a given ID.
	SubmissionGetByID(id int) (*Submission, error)

	// SubmissionGetByMinResolution returns all Submissions whose attached image's resolution
	// is greater or equal to the specified constraints.
	SubmissionGetByMinResolution(minResolution *monitor.Resolution) ([]*Submission, error)

	// SubmissionGetByPostID returns the Submission for a given Reddit post ID.
	SubmissionGetByPostID(postID string) (*Submission, error)

	// SubmissionSearch returns all submissions whose title contains the specified text.
	// The search SHOULD BE case-insensitive.
	SubmissionSearch(text string) ([]*Submission, error)

	// SubmissionGetRandom returns a randomly selected Submission whose attached image's
	// resolution is greater or equal to the specified constraints.
	// The Submission SHOULD NOT already be present in the History.
	SubmissionGetRandom(minResolution *monitor.Resolution) (*Submission, error)

	// SubmissionCreate creates and persists a Submission.
	SubmissionCreate(submission *Submission) error
}
