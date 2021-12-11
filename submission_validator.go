package redwall

import (
	"errors"
	"strings"
)

var _ SubmissionRepository = &submissionValidator{}

type submissionValidationFn func(*Submission) error

type submissionValidator struct {
	SubmissionRepository
}

func (v *submissionValidator) runValidationFns(submission *Submission, fns ...submissionValidationFn) error {
	for _, fn := range fns {
		if err := fn(submission); err != nil {
			return err
		}
	}

	return nil
}

func (v *submissionValidator) normalizePostID(submission *Submission) error {
	submission.PostID = strings.TrimSpace(submission.PostID)

	return nil
}

func (v *submissionValidator) requirePostID(submission *Submission) error {
	if submission.PostID == "" {
		return errors.New("Empty post ID")
	}

	return nil
}

func (v *submissionValidator) ByPostID(postID string) (*Submission, error) {
	submission := &Submission{PostID: postID}

	err := v.runValidationFns(
		submission,
		v.normalizePostID,
		v.requirePostID,
	)

	if err != nil {
		return &Submission{}, err
	}

	return v.SubmissionRepository.ByPostID(postID)
}

func newSubmissionValidator(submissionRepository SubmissionRepository) *submissionValidator {
	return &submissionValidator{
		SubmissionRepository: submissionRepository,
	}
}
