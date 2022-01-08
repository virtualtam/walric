package submission

import (
	"errors"
	"testing"

	"github.com/virtualtam/redwall2/subreddit"
)

func TestServiceByID(t *testing.T) {
	testCases := []struct {
		tname                 string
		repositorySubreddits  []*subreddit.Subreddit
		repositorySubmissions []*Submission
		id                    int
		want                  *Submission
		wantErr               error
	}{
		{
			tname: "existing ID",
			repositorySubreddits: []*subreddit.Subreddit{
				{ID: 1, Name: "astrophotography"},
			},
			repositorySubmissions: []*Submission{
				{ID: 1, SubredditID: 1, Title: "Messier 31 - The Andromeda Galaxy"},
				{ID: 2, SubredditID: 1, Title: "The Owl Nebula and Surfboard Galaxy"},
			},
			id: 2,
			want: &Submission{
				ID:          2,
				SubredditID: 1,
				Title:       "The Owl Nebula and Surfboard Galaxy",
				Subreddit:   &subreddit.Subreddit{ID: 1, Name: "astrophotography"},
			},
		},
		{
			tname: "unknown ID",
			repositorySubmissions: []*Submission{
				{ID: 1, Title: "Messier 31 - The Andromeda Galaxy"},
				{ID: 2, Title: "The Owl Nebula and Surfboard Galaxy"},
			},
			id:      649,
			wantErr: ErrNotFound,
		},
		{
			tname: "unknown ID",
			repositorySubmissions: []*Submission{
				{ID: 1, Title: "Messier 31 - The Andromeda Galaxy"},
				{ID: 2, Title: "The Owl Nebula and Surfboard Galaxy"},
			},
			id:      649,
			wantErr: ErrNotFound,
		},
		{
			tname:   "empty repository",
			id:      649,
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.tname, func(t *testing.T) {
			subredditRepository := subreddit.NewRepositoryInMemory(tt.repositorySubreddits)
			subredditService := subreddit.NewService(subredditRepository)

			repository := NewRepositoryInMemory(tt.repositorySubmissions)
			service := NewService(repository, subredditService)

			submission, err := service.ByID(tt.id)

			if tt.wantErr != nil {
				if err == nil {
					t.Error("expected an errot but got none")
				}

				if !errors.Is(err, tt.wantErr) {
					t.Errorf("want error %q, got %q", tt.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error, got %q", err)
			}

			assertSubmissionEquals(t, tt.want, submission)
			assertSubmissionSubredditEquals(t, tt.want, submission)
		})
	}
}

func TestServiceByPostID(t *testing.T) {
	testCases := []struct {
		tname                 string
		repositorySubreddits  []*subreddit.Subreddit
		repositorySubmissions []*Submission
		postID                string
		want                  *Submission
		wantErr               error
	}{
		{
			tname: "existing ID",
			repositorySubreddits: []*subreddit.Subreddit{
				{ID: 1, Name: "astrophotography"},
			},
			repositorySubmissions: []*Submission{
				{ID: 1, PostID: "m31aga", SubredditID: 1, Title: "Messier 31 - The Andromeda Galaxy"},
				{ID: 2, PostID: "owlsrf", SubredditID: 1, Title: "The Owl Nebula and Surfboard Galaxy"},
			},
			postID: "owlsrf",
			want: &Submission{
				ID:          2,
				PostID:      "owlsrf",
				SubredditID: 1,
				Title:       "The Owl Nebula and Surfboard Galaxy",
				Subreddit:   &subreddit.Subreddit{ID: 1, Name: "astrophotography"},
			},
		},
		{
			tname: "unknown ID",
			repositorySubmissions: []*Submission{
				{ID: 1, Title: "Messier 31 - The Andromeda Galaxy"},
				{ID: 2, Title: "The Owl Nebula and Surfboard Galaxy"},
			},
			postID:  "unkwn",
			wantErr: ErrNotFound,
		},
		{
			tname:   "empty repository",
			postID:  "unkwn",
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.tname, func(t *testing.T) {
			subredditRepository := subreddit.NewRepositoryInMemory(tt.repositorySubreddits)
			subredditService := subreddit.NewService(subredditRepository)

			repository := NewRepositoryInMemory(tt.repositorySubmissions)
			service := NewService(repository, subredditService)

			submission, err := service.ByPostID(tt.postID)

			if tt.wantErr != nil {
				if err == nil {
					t.Error("expected an errot but got none")
				}

				if !errors.Is(err, tt.wantErr) {
					t.Errorf("want error %q, got %q", tt.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error, got %q", err)
			}

			assertSubmissionEquals(t, tt.want, submission)
			assertSubmissionSubredditEquals(t, tt.want, submission)
		})
	}
}

func TestServiceSearch(t *testing.T) {
	testCases := []struct {
		tname                 string
		repositorySubreddits  []*subreddit.Subreddit
		repositorySubmissions []*Submission
		text                  string
		want                  []*Submission
		wantErr               error
	}{
		{
			tname: "single result",
			repositorySubreddits: []*subreddit.Subreddit{
				{ID: 1, Name: "astrophotography"},
			},
			repositorySubmissions: []*Submission{
				{ID: 1, PostID: "m31aga", SubredditID: 1, Title: "Messier 31 - The Andromeda Galaxy"},
				{ID: 2, PostID: "owlsrf", SubredditID: 1, Title: "The Owl Nebula and Surfboard Galaxy"},
			},
			text: "nebula",
			want: []*Submission{
				{
					ID:          2,
					PostID:      "owlsrf",
					SubredditID: 1,
					Title:       "The Owl Nebula and Surfboard Galaxy",
					Subreddit:   &subreddit.Subreddit{ID: 1, Name: "astrophotography"},
				},
			},
		},
		{
			tname: "multiple results",
			repositorySubreddits: []*subreddit.Subreddit{
				{ID: 1, Name: "astrophotography"},
			},
			repositorySubmissions: []*Submission{
				{ID: 1, PostID: "m31aga", SubredditID: 1, Title: "Messier 31 - The Andromeda Galaxy"},
				{ID: 2, PostID: "owlsrf", SubredditID: 1, Title: "The Owl Nebula and Surfboard Galaxy"},
			},
			text: "galaxy",
			want: []*Submission{
				{
					ID:          1,
					PostID:      "m31aga",
					SubredditID: 1,
					Title:       "Messier 31 - The Andromeda Galaxy",
					Subreddit:   &subreddit.Subreddit{ID: 1, Name: "astrophotography"},
				},
				{
					ID:          2,
					PostID:      "owlsrf",
					SubredditID: 1,
					Title:       "The Owl Nebula and Surfboard Galaxy",
					Subreddit:   &subreddit.Subreddit{ID: 1, Name: "astrophotography"},
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.tname, func(t *testing.T) {
			subredditRepository := subreddit.NewRepositoryInMemory(tt.repositorySubreddits)
			subredditService := subreddit.NewService(subredditRepository)

			repository := NewRepositoryInMemory(tt.repositorySubmissions)
			service := NewService(repository, subredditService)

			submissions, err := service.Search(tt.text)

			if tt.wantErr != nil {
				if err == nil {
					t.Error("expected an errot but got none")
				}

				if !errors.Is(err, tt.wantErr) {
					t.Errorf("want error %q, got %q", tt.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error, got %q", err)
			}

			if len(submissions) != len(tt.want) {
				t.Errorf("want %d submissions, got %d", len(tt.want), len(submissions))
			}

			for index, want := range tt.want {
				assertSubmissionEquals(t, want, submissions[index])
				assertSubmissionSubredditEquals(t, want, submissions[index])
			}
		})
	}
}

func assertSubmissionEquals(t *testing.T, want, got *Submission) {
	t.Helper()

	if got.ID != want.ID {
		t.Errorf("want ID %d, got %d", want.ID, got.ID)
	}
	if got.Title != want.Title {
		t.Errorf("want title %q, got %q", want.Title, got.Title)
	}
	if got.SubredditID != want.SubredditID {
		t.Errorf("want subreddit ID %d, got %d", want.SubredditID, got.SubredditID)
	}
}

func assertSubmissionSubredditEquals(t *testing.T, want, got *Submission) {
	t.Helper()

	if got.Subreddit == nil {
		t.Error("expected subreddit metadata to be set but it is not")
		return
	}

	if got.Subreddit.ID != want.Subreddit.ID {
		t.Errorf("want subreddit ID %d, got %d", want.Subreddit.ID, got.Subreddit.ID)
	}
	if got.Subreddit.Name != want.Subreddit.Name {
		t.Errorf("want subreddit name %q, got %q", want.Subreddit.Name, got.Subreddit.Name)
	}
}
