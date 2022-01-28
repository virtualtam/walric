package submission

import "github.com/virtualtam/walric/monitor"

// Repository defines the basic operations available to access and persist
// Reddit Submissions.
type Repository interface {
	// ByID returns the Submission for a given ID.
	ByID(id int) (*Submission, error)

	// ByPostID returns the Submission for a given Reddit post ID.
	ByPostID(postID string) (*Submission, error)

	// Search returns all submissions whose title contains the specified text.
	// The search SHOULD BE case-insensitive.
	Search(text string) ([]*Submission, error)

	// ByMinResolution returns all Submissions whose attached image's resolution
	// is greater or equal to the specified constraints.
	ByMinResolution(minResolution *monitor.Resolution) ([]*Submission, error)

	// Random returns a randomly selected Submission whose attached image's
	// resolution is greater or equal to the specified constraints.
	// The Submission SHOULD NOT already be present in the History.
	Random(minResolution *monitor.Resolution) (*Submission, error)

	// Create creates and persists a Submission.
	Create(submission *Submission) error
}
