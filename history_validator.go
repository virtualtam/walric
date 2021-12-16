package redwall

import "errors"

var _ HistoryRepository = &historyValidator{}

type historyValidationFn func(entry *HistoryEntry) error

type historyValidator struct {
	HistoryRepository
}

func (v *historyValidator) runValidationFns(entry *HistoryEntry, fns ...historyValidationFn) error {
	for _, fn := range fns {
		if err := fn(entry); err != nil {
			return err
		}
	}

	return nil
}

func (v *historyValidator) requirePositiveSubmissionID(entry *HistoryEntry) error {
	if entry.SubmissionID < 1 {
		return errors.New("Negative submission ID")
	}

	return nil
}

func (v *historyValidator) Create(entry *HistoryEntry) error {
	err := v.runValidationFns(
		entry,
		v.requirePositiveSubmissionID,
	)
	if err != nil {
		return err
	}

	return v.HistoryRepository.Create(entry)
}

func newHistoryValidator(historyRepository HistoryRepository) *historyValidator {
	return &historyValidator{
		HistoryRepository: historyRepository,
	}
}
