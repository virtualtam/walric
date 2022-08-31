package history

type validationFn func(entry *Entry) error

// validator provides a normalization and validation layer for history Entry
// Repository operations.
type validator struct {
	Repository
}

func (v *validator) runValidationFns(entry *Entry, fns ...validationFn) error {
	for _, fn := range fns {
		if err := fn(entry); err != nil {
			return err
		}
	}

	return nil
}

func (v *validator) requirePositiveSubmissionID(entry *Entry) error {
	if entry.SubmissionID < 1 {
		return ErrSubmissionIDNegativeOrZero
	}

	return nil
}

func (v *validator) Create(entry *Entry) error {
	err := v.runValidationFns(
		entry,
		v.requirePositiveSubmissionID,
	)
	if err != nil {
		return err
	}

	return v.Repository.Create(entry)
}

func newValidator(repository Repository) *validator {
	return &validator{
		Repository: repository,
	}
}
