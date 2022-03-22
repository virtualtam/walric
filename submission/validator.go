package submission

import (
	"errors"
	"strings"

	"github.com/virtualtam/walric/monitor"
)

type validationFn func(*Submission) error

type validator struct {
	Repository
}

func (v *validator) runValidationFns(submission *Submission, fns ...validationFn) error {
	for _, fn := range fns {
		if err := fn(submission); err != nil {
			return err
		}
	}

	return nil
}

func (v *validator) requirePositiveSubredditID(submission *Submission) error {
	if submission.SubredditID <= 0 {
		return ErrSubredditIDInvalid
	}

	return nil
}

func (v *validator) requireDefaultID(submission *Submission) error {
	if submission.ID != 0 {
		return ErrIDInvalid
	}

	return nil
}

func (v *validator) requirePositiveID(submission *Submission) error {
	if submission.ID <= 0 {
		return ErrIDInvalid
	}

	return nil
}

func (v *validator) ensurePostIDIsNotRegistered(submission *Submission) error {
	_, err := v.ByPostID(submission.PostID)

	if errors.Is(err, ErrNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	return ErrPostIDAlreadyRegistered
}

func (v *validator) normalizePostID(submission *Submission) error {
	submission.PostID = strings.TrimSpace(submission.PostID)

	return nil
}

func (v *validator) requirePostID(submission *Submission) error {
	if submission.PostID == "" {
		return ErrPostIDEmpty
	}

	return nil
}

func (v *validator) normalizeTitle(submission *Submission) error {
	submission.Title = strings.TrimSpace(submission.Title)

	return nil
}

func (v *validator) requireTitle(submission *Submission) error {
	if submission.Title == "" {
		return ErrTitleEmpty
	}

	return nil
}

func (v *validator) ByID(id int) (*Submission, error) {
	submission := &Submission{ID: id}

	err := v.runValidationFns(
		submission,
		v.requirePositiveID,
	)

	if err != nil {
		return &Submission{}, err
	}

	return v.Repository.ByID(id)
}

func (v *validator) ByPostID(postID string) (*Submission, error) {
	submission := &Submission{PostID: postID}

	err := v.runValidationFns(
		submission,
		v.normalizePostID,
		v.requirePostID,
	)

	if err != nil {
		return &Submission{}, err
	}

	return v.Repository.ByPostID(postID)
}

func (v *validator) Search(text string) ([]*Submission, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return []*Submission{}, ErrSearchTextEmpty
	}

	return v.Repository.Search(text)
}

func (v *validator) ByMinResolution(minResolution *monitor.Resolution) ([]*Submission, error) {
	if minResolution.HeightPx < 1 || minResolution.WidthPx < 1 {
		return []*Submission{}, ErrResolutionInvalid
	}

	return v.Repository.ByMinResolution(minResolution)
}

func (v *validator) Random(minResolution *monitor.Resolution) (*Submission, error) {
	if minResolution.HeightPx < 1 || minResolution.WidthPx < 1 {
		return &Submission{}, ErrResolutionInvalid
	}

	return v.Repository.Random(minResolution)
}

func (v *validator) Create(submission *Submission) error {
	err := v.runValidationFns(
		submission,
		v.requirePositiveSubredditID,
		v.requireDefaultID,
		v.normalizePostID,
		v.requirePostID,
		v.ensurePostIDIsNotRegistered,
		v.normalizeTitle,
		v.requireTitle,
	)

	if err != nil {
		return err
	}

	return v.Repository.Create(submission)
}

func newValidator(repository Repository) *validator {
	return &validator{
		Repository: repository,
	}
}
