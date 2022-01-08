package submission

import (
	"strings"

	"github.com/virtualtam/redwall2/monitor"
)

var _ Repository = &validator{}

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

func (v *validator) requirePositiveID(submission *Submission) error {
	if submission.ID <= 0 {
		return ErrIDInvalid
	}

	return nil
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

func newValidator(repository Repository) *validator {
	return &validator{
		Repository: repository,
	}
}
