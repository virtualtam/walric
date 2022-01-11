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

func TestValidatorCreate(t *testing.T) {
	testCases := []struct {
		tname                string
		repositorySubreddits []*Subreddit
		subreddit            *Subreddit
		wantErr              error
	}{
		// nominal cases
		{
			tname:     "new subreddit",
			subreddit: &Subreddit{Name: "FromSpaceWithLove"},
		},
		{
			tname: "duplicate subreddit",
			repositorySubreddits: []*Subreddit{
				{Name: "FromSpaceWithLove"},
			},
			subreddit: &Subreddit{Name: "FromSpaceWithLove"},
			wantErr:   ErrNameAlreadyRegistered,
		},

		// error cases
		{
			tname:     "empty name",
			subreddit: &Subreddit{Name: ""},
			wantErr:   ErrNameEmpty,
		},
		{
			tname:     "empty name (whitespace)",
			subreddit: &Subreddit{Name: "   "},
			wantErr:   ErrNameEmpty,
		},
		{
			tname:     "non-default ID",
			subreddit: &Subreddit{ID: 12, Name: "NonZero"},
			wantErr:   ErrIDInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tt.repositorySubreddits)
			currentID := repository.currentID
			validator := newValidator(repository)

			err := validator.Create(tt.subreddit)

			if tt.wantErr != nil {
				if err == nil {
					t.Error("expected an error but got none")
				} else if !errors.Is(err, tt.wantErr) {
					t.Errorf("want error %q, got %q", tt.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error but got %q", err)
			}

			subreddit, err := validator.ByID(currentID)
			if err != nil {
				t.Errorf("failed to retrieve subreddit: %q", err)
			}

			if subreddit.ID != currentID {
				t.Errorf("want ID %d, got %d", currentID, subreddit.ID)
			}
			if subreddit.Name != tt.subreddit.Name {
				t.Errorf("want name %q, got %q", tt.subreddit.Name, subreddit.Name)
			}
		})
	}
}
