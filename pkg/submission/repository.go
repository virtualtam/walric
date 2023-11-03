package submission

import "github.com/virtualtam/walric/pkg/monitor"

// ValidationRepository provides methods for Submission validation.
type ValidationRepository interface {
	// SubmissionIsPostIDRegistered returns whether this Submission was previously saved.
	SubmissionIsPostIDRegistered(postID string) (bool, error)

	// SubredditIsNameRegistered returns whether this Subreddit was previously saved.
	SubredditIsNameRegistered(name string) (bool, error)
}

// Repository defines the basic operations available to access and persist
// Reddit Submissions.
type Repository interface {
	ValidationRepository

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

	// SubredditGetAll returns all persisted Subreddits.
	SubredditGetAll() ([]*Subreddit, error)

	// SubredditGetStats returns the aggregated usage statistics for all Subreddits.
	SubredditGetStats() ([]SubredditStats, error)

	// SubredditGetByID returns the Subreddit for a given ID.
	SubredditGetByID(id int) (*Subreddit, error)

	// SubredditGetByName returns the Subreddit for a given Name.
	SubredditGetByName(name string) (*Subreddit, error)

	// SubredditCreate creates and persists a Subreddit.
	SubredditCreate(subreddit *Subreddit) error
}
