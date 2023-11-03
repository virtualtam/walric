package submission

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/virtualtam/walric/pkg/monitor"
)

var _ Repository = &RepositoryInMemory{}

// repositoryInMemory provides an in-memory Repository for testing.
type RepositoryInMemory struct {
	submissionCurrentID int
	submissions         []*Submission

	subredditCurrentID int
	subreddits         []*Subreddit
}

func NewRepositoryInMemory(submissions []*Submission, subreddits []*Subreddit) *RepositoryInMemory {
	return &RepositoryInMemory{
		submissionCurrentID: len(submissions) + 1,
		submissions:         submissions,

		subredditCurrentID: len(subreddits) + 1,
		subreddits:         subreddits,
	}
}

func (r *RepositoryInMemory) SubmissionGetByID(id int) (*Submission, error) {
	for _, submission := range r.submissions {
		if submission.ID == id {
			return submission, nil
		}
	}

	return &Submission{}, ErrSubmissionNotFound
}

func (r *RepositoryInMemory) SubmissionGetByPostID(postID string) (*Submission, error) {
	for _, submission := range r.submissions {
		if submission.PostID == postID {
			return submission, nil
		}
	}

	return &Submission{}, ErrSubmissionNotFound
}

func (r *RepositoryInMemory) SubmissionIsPostIDRegistered(postID string) (bool, error) {
	for _, submission := range r.submissions {
		if submission.PostID == postID {
			return true, nil
		}
	}

	return false, nil
}

func (r *RepositoryInMemory) SubmissionSearch(text string) ([]*Submission, error) {
	results := []*Submission{}

	for _, submission := range r.submissions {
		if strings.Contains(strings.ToLower(submission.Title), strings.ToLower(text)) {
			results = append(results, submission)
		}
	}

	return results, nil
}

func (r *RepositoryInMemory) SubmissionGetByMinResolution(minResolution *monitor.Resolution) ([]*Submission, error) {
	candidates := []*Submission{}
	for _, submission := range r.submissions {
		if submission.ImageHeightPx >= minResolution.HeightPx && submission.ImageWidthPx >= minResolution.WidthPx {
			candidates = append(candidates, submission)
		}
	}

	return candidates, nil
}

func (r *RepositoryInMemory) SubmissionGetRandom(minResolution *monitor.Resolution) (*Submission, error) {
	if len(r.submissions) == 0 {
		return &Submission{}, ErrSubmissionNotFound
	}

	candidates, err := r.SubmissionGetByMinResolution(minResolution)
	if err != nil {
		return &Submission{}, nil
	}

	if len(candidates) == 0 {
		return &Submission{}, ErrSubmissionNotFound
	}

	index := rand.Intn(len(candidates))

	return candidates[index], nil
}

func (r *RepositoryInMemory) SubmissionCreate(submission *Submission) error {
	submission.ID = r.submissionCurrentID
	r.submissionCurrentID++

	r.submissions = append(r.submissions, submission)

	return nil
}

func (r *RepositoryInMemory) SubredditCreate(subreddit *Subreddit) error {
	subreddit.ID = r.subredditCurrentID
	r.subredditCurrentID++

	r.subreddits = append(r.subreddits, subreddit)

	return nil
}

func (r *RepositoryInMemory) SubredditGetAll() ([]*Subreddit, error) {
	return r.subreddits, nil
}

func (r *RepositoryInMemory) SubredditGetStats() ([]SubredditStats, error) {
	return []SubredditStats{}, errors.New("not implemented")
}

func (r *RepositoryInMemory) SubredditGetByID(id int) (*Subreddit, error) {
	for _, subreddit := range r.subreddits {
		if subreddit.ID == id {
			return subreddit, nil
		}
	}

	return &Subreddit{}, ErrSubredditNotFound
}

func (r *RepositoryInMemory) SubredditGetByName(name string) (*Subreddit, error) {
	for _, subreddit := range r.subreddits {
		if subreddit.Name == name {
			return subreddit, nil
		}
	}

	return &Subreddit{}, ErrSubredditNotFound
}

func (r *RepositoryInMemory) SubredditIsNameRegistered(name string) (bool, error) {
	for _, subreddit := range r.subreddits {
		if subreddit.Name == name {
			return true, nil
		}
	}

	return false, nil
}
