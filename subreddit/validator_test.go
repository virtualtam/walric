package subreddit

import (
	"errors"
	"testing"
)

func TestValidatorByID(t *testing.T) {
	testCases := []struct {
		tname                string
		repositorySubreddits []*Subreddit
		id                   int
		want                 *Subreddit
		wantErr              error
	}{
		// nominal cases
		{
			tname: "return by ID",
			repositorySubreddits: []*Subreddit{
				{ID: 1, Name: "astrophotography"},
				{ID: 2, Name: "FromSpaceWithLove"},
				{ID: 3, Name: "Museum"},
			},
			id:   2,
			want: &Subreddit{ID: 2, Name: "FromSpaceWithLove"},
		},
		{
			tname:   "not found",
			id:      9362,
			wantErr: ErrNotFound,
		},

		// error cases
		{
			tname:   "ID is negative",
			id:      -548,
			wantErr: ErrIDInvalid,
		},
		{
			tname:   "ID equals zero",
			id:      0,
			wantErr: ErrIDInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tt.repositorySubreddits)
			validator := newValidator(repository)

			subreddit, err := validator.ByID(tt.id)

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

			if subreddit.ID != tt.want.ID {
				t.Errorf("want ID %d, got %d", tt.want.ID, subreddit.ID)
			}
			if subreddit.Name != tt.want.Name {
				t.Errorf("want name %q, got %q", tt.want.Name, subreddit.Name)
			}
		})
	}
}

func TestValidatorByName(t *testing.T) {
	testCases := []struct {
		tname                string
		repositorySubreddits []*Subreddit
		name                 string
		want                 *Subreddit
		wantErr              error
	}{
		// nominal cases
		{
			tname: "return by ID",
			repositorySubreddits: []*Subreddit{
				{ID: 1, Name: "astrophotography"},
				{ID: 2, Name: "FromSpaceWithLove"},
				{ID: 3, Name: "Museum"},
			},
			name: "FromSpaceWithLove",
			want: &Subreddit{ID: 2, Name: "FromSpaceWithLove"},
		},
		{
			tname:   "not found",
			name:    "Unknown",
			wantErr: ErrNotFound,
		},

		// error cases
		{
			tname:   "empty name",
			name:    "",
			wantErr: ErrNameEmpty,
		},
		{
			tname:   "empty name (whitespace)",
			name:    "     ",
			wantErr: ErrNameEmpty,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tt.repositorySubreddits)
			validator := newValidator(repository)

			subreddit, err := validator.ByName(tt.name)

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

			if subreddit.ID != tt.want.ID {
				t.Errorf("want ID %d, got %d", tt.want.ID, subreddit.ID)
			}
			if subreddit.Name != tt.want.Name {
				t.Errorf("want name %q, got %q", tt.want.Name, subreddit.Name)
			}
		})
	}
}
