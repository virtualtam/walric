package history

import (
	"errors"
	"testing"
	"time"
)

func TestValidatorCreate(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		tname             string
		repositoryEntries []*Entry
		entry             *Entry
		wantErr           error
	}{
		// nominal cases
		{
			tname: "new entry",
			entry: &Entry{
				Date:         now,
				SubmissionID: 856,
			},
		},
		{
			tname: "new duplicate entry",
			repositoryEntries: []*Entry{
				{
					ID:           1,
					Date:         now,
					SubmissionID: 856,
				},
			},
			entry: &Entry{
				Date:         now,
				SubmissionID: 856,
			},
		},

		// error cases
		{
			tname:   "submission ID is negative",
			entry:   &Entry{ID: -67},
			wantErr: ErrSubmissionIDNegativeOrZero,
		},
		{
			tname:   "submission ID equals zero",
			entry:   &Entry{ID: 0},
			wantErr: ErrSubmissionIDNegativeOrZero,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.tname, func(t *testing.T) {
			nEntries := len(tt.repositoryEntries)

			repository := &repositoryInMemory{
				entries: tt.repositoryEntries,
			}
			validator := newValidator(repository)

			err := validator.Create(tt.entry)

			if tt.wantErr != nil {
				if err == nil {
					t.Error("expected an error but got none")
				}

				if !errors.Is(err, tt.wantErr) {
					t.Errorf("want error %q, got %q", tt.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error but got %q", err)
			}

			wantNEntries := nEntries + 1
			if len(repository.entries) != wantNEntries {
				t.Errorf("want %d entries, got %d", wantNEntries, len(repository.entries))
			}

			entry, err := repository.Current()
			if err != nil {
				t.Errorf("failed to retrieve entry: %q", err)
			}

			if entry.SubmissionID != tt.entry.SubmissionID {
				t.Errorf("want submission ID %d, got %d", tt.entry.SubmissionID, entry.SubmissionID)
			}
		})
	}
}
