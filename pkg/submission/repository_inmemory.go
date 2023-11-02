package submission

import (
	"math/rand"
	"strings"

	"github.com/virtualtam/walric/pkg/monitor"
)

var _ Repository = &RepositoryInMemory{}

// repositoryInMemory provides an in-memory Repository for testing.
type RepositoryInMemory struct {
	currentID   int
	submissions []*Submission
}

func (r *RepositoryInMemory) SubmissionGetByID(id int) (*Submission, error) {
	for _, submission := range r.submissions {
		if submission.ID == id {
			return submission, nil
		}
	}

	return &Submission{}, ErrNotFound
}

func (r *RepositoryInMemory) SubmissionGetByPostID(postID string) (*Submission, error) {
	for _, submission := range r.submissions {
		if submission.PostID == postID {
			return submission, nil
		}
	}

	return &Submission{}, ErrNotFound
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
		return &Submission{}, ErrNotFound
	}

	candidates, err := r.SubmissionGetByMinResolution(minResolution)
	if err != nil {
		return &Submission{}, nil
	}

	if len(candidates) == 0 {
		return &Submission{}, ErrNotFound
	}

	index := rand.Intn(len(candidates))

	return candidates[index], nil
}

func (r *RepositoryInMemory) SubmissionCreate(submission *Submission) error {
	submission.ID = r.currentID
	r.currentID++

	r.submissions = append(r.submissions, submission)

	return nil
}

func NewRepositoryInMemory(submissions []*Submission) *RepositoryInMemory {
	return &RepositoryInMemory{
		currentID:   len(submissions) + 1,
		submissions: submissions,
	}
}
