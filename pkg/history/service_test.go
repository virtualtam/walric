package history

import (
	"errors"
	"testing"
	"time"

	"github.com/virtualtam/walric/pkg/submission"
)

func TestServiceCreate(t *testing.T) {
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
				Date:       now,
				Submission: &submission.Submission{ID: 856},
			},
		},
		{
			tname: "new duplicate entry",
			repositoryEntries: []*Entry{
				{
					ID:         1,
					Date:       now,
					Submission: &submission.Submission{ID: 856},
				},
			},
			entry: &Entry{
				Date:       now,
				Submission: &submission.Submission{ID: 856},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			nEntries := len(tc.repositoryEntries)

			repository := &repositoryInMemory{
				entries: tc.repositoryEntries,
			}
			service := NewService(repository, nil)

			err := service.Save(tc.entry)

			if tc.wantErr != nil {
				if err == nil {
					t.Error("expected an error but got none")
				} else if !errors.Is(err, tc.wantErr) {
					t.Errorf("want error %q, got %q", tc.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error, got %q", err)
				return
			}

			wantNEntries := nEntries + 1
			if len(repository.entries) != wantNEntries {
				t.Errorf("want %d entries, got %d", wantNEntries, len(repository.entries))
				return
			}

			entry, err := repository.HistoryGetCurrent()
			if err != nil {
				t.Errorf("failed to retrieve entry: %q", err)
				return
			}

			if entry.Submission.ID != tc.entry.Submission.ID {
				t.Errorf("want submission ID %d, got %d", tc.entry.Submission.ID, entry.Submission.ID)
			}
		})
	}
}
