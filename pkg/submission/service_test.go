package submission

import (
	"errors"
	"testing"
	"time"

	"github.com/virtualtam/walric/pkg/monitor"
)

func TestServiceByID(t *testing.T) {
	testCases := []struct {
		tname                 string
		repositorySubreddits  []*Subreddit
		repositorySubmissions []*Submission
		id                    int
		want                  *Submission
		wantErr               error
	}{
		// nominal cases
		{
			tname: "existing ID",
			repositorySubreddits: []*Subreddit{
				{ID: 1, Name: "astrophotography"},
			},
			repositorySubmissions: []*Submission{
				{ID: 1, Subreddit: &Subreddit{ID: 1}, Title: "Messier 31 - The Andromeda Galaxy"},
				{ID: 2, Subreddit: &Subreddit{ID: 1}, Title: "The Owl Nebula and Surfboard Galaxy"},
			},
			id: 2,
			want: &Submission{
				ID:        2,
				Title:     "The Owl Nebula and Surfboard Galaxy",
				Subreddit: &Subreddit{ID: 1, Name: "astrophotography"},
			},
		},
		{
			tname: "unknown ID",
			repositorySubmissions: []*Submission{
				{ID: 1, Title: "Messier 31 - The Andromeda Galaxy"},
				{ID: 2, Title: "The Owl Nebula and Surfboard Galaxy"},
			},
			id:      649,
			wantErr: ErrSubmissionNotFound,
		},
		{
			tname: "unknown ID",
			repositorySubmissions: []*Submission{
				{ID: 1, Title: "Messier 31 - The Andromeda Galaxy"},
				{ID: 2, Title: "The Owl Nebula and Surfboard Galaxy"},
			},
			id:      649,
			wantErr: ErrSubmissionNotFound,
		},
		{
			tname:   "empty repository",
			id:      649,
			wantErr: ErrSubmissionNotFound,
		},

		// error cases
		{
			tname:   "ID is negative",
			id:      -548,
			wantErr: ErrSubmissionIDInvalid,
		},
		{
			tname:   "ID equals zero",
			id:      0,
			wantErr: ErrSubmissionIDInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tc.repositorySubmissions, tc.repositorySubreddits)
			service := NewService(repository)

			submission, err := service.ByID(tc.id)

			if tc.wantErr != nil {
				if err == nil {
					t.Error("expected an errot but got none")
				} else if !errors.Is(err, tc.wantErr) {
					t.Errorf("want error %q, got %q", tc.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error, got %q", err)
				return
			}

			assertSubmissionEquals(t, tc.want, submission)
			assertSubmissionSubredditEquals(t, tc.want, submission)
		})
	}
}

func TestServiceByMinResolution(t *testing.T) {
	repositorySubreddits := []*Subreddit{
		{
			ID:   1,
			Name: "Dummy",
		},
	}

	testCases := []struct {
		tname                 string
		repositorySubmissions []*Submission
		minResolution         *monitor.Resolution
		want                  []*Submission
		wantErr               error
	}{
		// nominal cases
		{
			tname: "single result",
			repositorySubmissions: []*Submission{
				{
					Title:         "Sunday Afternoon In The Park [640x480]",
					ImageHeightPx: 480,
					ImageWidthPx:  640,
					Subreddit:     &Subreddit{ID: 1},
				},
				{
					Title:         "Laguna Sunrise [1920x1200]",
					ImageHeightPx: 1200,
					ImageWidthPx:  1920,
					Subreddit:     &Subreddit{ID: 1},
				},
			},
			minResolution: &monitor.Resolution{HeightPx: 1200, WidthPx: 1920},
			want: []*Submission{
				{Title: "Laguna Sunrise [1920x1200]", ImageHeightPx: 1200, ImageWidthPx: 1920},
			},
		},
		{
			tname: "multiple results",
			repositorySubmissions: []*Submission{
				{Title: "Sunday Afternoon In The Park [640x480]",
					ImageHeightPx: 480,
					ImageWidthPx:  640,
					Subreddit:     &Subreddit{ID: 1},
				},
				{Title: "Moroccan Sunset [2560x1440]",
					ImageHeightPx: 1440,
					ImageWidthPx:  2560,
					Subreddit:     &Subreddit{ID: 1},
				},
				{Title: "Laguna Sunrise [1920x1200]",
					ImageHeightPx: 1200,
					ImageWidthPx:  1920,
					Subreddit:     &Subreddit{ID: 1},
				},
			},
			minResolution: &monitor.Resolution{HeightPx: 1200, WidthPx: 1920},
			want: []*Submission{
				{Title: "Moroccan Sunset [2560x1440]", ImageHeightPx: 1440, ImageWidthPx: 2560},
				{Title: "Laguna Sunrise [1920x1200]", ImageHeightPx: 1200, ImageWidthPx: 1920},
			},
		},
		{
			tname: "no result for this resolution",
			repositorySubmissions: []*Submission{
				{Title: "Sunday Afternoon In The Park [640x480]", ImageHeightPx: 480, ImageWidthPx: 640},
				{Title: "Laguna Sunrise [1920x1200]", ImageHeightPx: 1200, ImageWidthPx: 1920},
			},
			minResolution: &monitor.Resolution{HeightPx: 1440, WidthPx: 2560},
		},
		{
			tname:         "no result (empty repository)",
			minResolution: &monitor.Resolution{HeightPx: 1200, WidthPx: 1920},
		},

		// error cases
		{
			tname:         "negative resolution height",
			minResolution: &monitor.Resolution{HeightPx: -1200, WidthPx: 1920},
			wantErr:       monitor.ErrResolutionInvalid,
		},
		{
			tname:         "negative resolution width",
			minResolution: &monitor.Resolution{HeightPx: 1200, WidthPx: -1920},
			wantErr:       monitor.ErrResolutionInvalid,
		},
		{
			tname:         "negative resolution height and width",
			minResolution: &monitor.Resolution{HeightPx: -1200, WidthPx: -1920},
			wantErr:       monitor.ErrResolutionInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tc.repositorySubmissions, repositorySubreddits)
			service := NewService(repository)

			submissions, err := service.ByMinResolution(tc.minResolution)

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

			if len(submissions) != len(tc.want) {
				t.Errorf("want %d submissions, got %d", len(tc.want), len(submissions))
			}
		})
	}
}

func TestServiceByPostID(t *testing.T) {
	testCases := []struct {
		tname                 string
		repositorySubreddits  []*Subreddit
		repositorySubmissions []*Submission
		postID                string
		want                  *Submission
		wantErr               error
	}{
		// nominal cases
		{
			tname: "existing ID",
			repositorySubreddits: []*Subreddit{
				{ID: 1, Name: "astrophotography"},
			},
			repositorySubmissions: []*Submission{
				{ID: 1, PostID: "m31aga", Subreddit: &Subreddit{ID: 1}, Title: "Messier 31 - The Andromeda Galaxy"},
				{ID: 2, PostID: "owlsrf", Subreddit: &Subreddit{ID: 1}, Title: "The Owl Nebula and Surfboard Galaxy"},
			},
			postID: "owlsrf",
			want: &Submission{
				ID:        2,
				PostID:    "owlsrf",
				Title:     "The Owl Nebula and Surfboard Galaxy",
				Subreddit: &Subreddit{ID: 1, Name: "astrophotography"},
			},
		},

		// error cases
		{
			tname: "unknown ID",
			repositorySubmissions: []*Submission{
				{ID: 1, Title: "Messier 31 - The Andromeda Galaxy"},
				{ID: 2, Title: "The Owl Nebula and Surfboard Galaxy"},
			},
			postID:  "unkwn",
			wantErr: ErrSubmissionNotFound,
		},
		{
			tname:   "empty repository",
			postID:  "unkwn",
			wantErr: ErrSubmissionNotFound,
		},
		{
			tname:   "empty PostID",
			postID:  "",
			wantErr: ErrSubmissionPostIDEmpty,
		},
		{
			tname:   "whitespace (empty) PostID",
			postID:  "       ",
			wantErr: ErrSubmissionPostIDEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tc.repositorySubmissions, tc.repositorySubreddits)
			service := NewService(repository)

			submission, err := service.ByPostID(tc.postID)

			if tc.wantErr != nil {
				if err == nil {
					t.Error("expected an errot but got none")
				} else if !errors.Is(err, tc.wantErr) {
					t.Errorf("want error %q, got %q", tc.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error, got %q", err)
				return
			}

			assertSubmissionEquals(t, tc.want, submission)
			assertSubmissionSubredditEquals(t, tc.want, submission)
		})
	}
}

func TestServiceCreate(t *testing.T) {
	testCases := []struct {
		tname                 string
		repositorySubmissions []*Submission
		repositorySubreddits  []*Subreddit
		submission            *Submission
		wantErr               error
	}{
		// nominal cases
		{
			tname: "new submission",
			repositorySubreddits: []*Subreddit{
				{
					ID:   25,
					Name: "Dummy",
				},
			},
			submission: &Submission{
				Subreddit:     &Subreddit{ID: 25},
				Author:        "janedoe",
				Permalink:     "r/dummy/comments/newnew/new_submission/",
				PostID:        "newnew",
				PostedAt:      time.Now().UTC(),
				Score:         12,
				Title:         "New Submission [800x600]",
				ImageDomain:   "i.redd.it",
				ImageURL:      "https://i.redd.it/newsub001.jpg",
				ImageNSFW:     false,
				ImageFilename: "/data/walric/dummy/newnew-newsub001.jpg",
				ImageHeightPx: 600,
				ImageWidthPx:  800,
			},
		},

		// error cases
		{
			tname: "negative subreddit ID",
			submission: &Submission{
				Subreddit: &Subreddit{ID: -583},
			},
			wantErr: ErrSubredditIDInvalid,
		},
		{
			tname: "subreddit ID equals zero",
			submission: &Submission{
				Subreddit: &Subreddit{ID: 0},
			},
			wantErr: ErrSubredditIDInvalid,
		},
		{
			tname: "non-default ID",
			submission: &Submission{
				Subreddit: &Subreddit{ID: 12},
				ID:        179,
				PostID:    "nondft",
				Title:     "Non-default [0x0]",
			},
			wantErr: ErrSubmissionIDInvalid,
		},
		{
			tname: "empty PostID",
			submission: &Submission{
				Subreddit: &Subreddit{ID: 12},
			},
			wantErr: ErrSubmissionPostIDEmpty,
		},
		{
			tname: "empty PostID (whitespace)",
			submission: &Submission{
				Subreddit: &Subreddit{ID: 12},
				PostID:    "     ",
			},
			wantErr: ErrSubmissionPostIDEmpty,
		},
		{
			tname: "duplicate PostID",
			repositorySubmissions: []*Submission{
				{
					Subreddit: &Subreddit{ID: 12},
					ID:        1,
					PostID:    "dupdup",
				},
			},
			submission: &Submission{
				Subreddit: &Subreddit{ID: 12},
				PostID:    "dupdup",
			},
			wantErr: ErrSubmissionPostIDAlreadyRegistered,
		},
		{
			tname: "empty title",
			submission: &Submission{
				Subreddit: &Subreddit{ID: 12},
				PostID:    "notitl",
			},
			wantErr: ErrSubmissionTitleEmpty,
		},
		{
			tname: "empty title (whitespace)",
			submission: &Submission{
				Subreddit: &Subreddit{ID: 12},
				PostID:    "notitle",
				Title:     "    ",
			},
			wantErr: ErrSubmissionTitleEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tc.repositorySubmissions, tc.repositorySubreddits)
			currentID := repository.submissionCurrentID
			service := NewService(repository)

			err := service.Create(tc.submission)

			if tc.wantErr != nil {
				if err == nil {
					t.Error("expected an error but got none")
				} else if !errors.Is(err, tc.wantErr) {
					t.Errorf("want error %q, got %q", tc.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error but got %q", err)
				return
			}

			submission, err := service.ByID(currentID)

			if err != nil {
				t.Errorf("failed to retrieve submission: %q", err)
				return
			}

			if submission.ID != currentID {
				t.Errorf("want ID %d, got %d", currentID, submission.ID)
			}
			if submission.Title != tc.submission.Title {
				t.Errorf("want title %q, got %q", tc.submission.Title, submission.Title)
			}
		})
	}
}

func TestServiceRandom(t *testing.T) {
	repositorySubreddits := []*Subreddit{
		{
			ID:   1,
			Name: "Dummy",
		},
	}

	testCases := []struct {
		tname                 string
		repositorySubmissions []*Submission
		minResolution         *monitor.Resolution
		want                  *Submission
		wantErr               error
	}{
		// nominal cases
		{
			tname: "random result",
			repositorySubmissions: []*Submission{
				{
					Title:         "Sunday Afternoon In The Park [640x480]",
					ImageHeightPx: 480,
					ImageWidthPx:  640,
					Subreddit:     &Subreddit{ID: 1},
				},
				{
					Title:         "Laguna Sunrise [1920x1200]",
					ImageHeightPx: 1200,
					ImageWidthPx:  1920,
					Subreddit:     &Subreddit{ID: 1},
				},
			},
			minResolution: &monitor.Resolution{HeightPx: 1200, WidthPx: 1920},
			want:          &Submission{Title: "Laguna Sunrise [1920x1200]", ImageHeightPx: 1200, ImageWidthPx: 1920},
		},
		{
			tname:         "not found (empty repository)",
			minResolution: &monitor.Resolution{HeightPx: 1200, WidthPx: 1920},
			wantErr:       ErrSubmissionNotFound,
		},
		{
			tname: "not found (no result)",
			repositorySubmissions: []*Submission{
				{Title: "Sunday Afternoon In The Park [640x480]", ImageHeightPx: 480, ImageWidthPx: 640},
				{Title: "Laguna Sunrise [1920x1200]", ImageHeightPx: 1200, ImageWidthPx: 1920},
			},
			minResolution: &monitor.Resolution{HeightPx: 1440, WidthPx: 2560},
			wantErr:       ErrSubmissionNotFound,
		},

		// error cases
		{
			tname:         "negative resolution height",
			minResolution: &monitor.Resolution{HeightPx: -1200, WidthPx: 1920},
			wantErr:       monitor.ErrResolutionInvalid,
		},
		{
			tname:         "negative resolution width",
			minResolution: &monitor.Resolution{HeightPx: 1200, WidthPx: -1920},
			wantErr:       monitor.ErrResolutionInvalid,
		},
		{
			tname:         "negative resolution height and width",
			minResolution: &monitor.Resolution{HeightPx: -1200, WidthPx: -1920},
			wantErr:       monitor.ErrResolutionInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tc.repositorySubmissions, repositorySubreddits)
			service := NewService(repository)

			submission, err := service.Random(tc.minResolution)

			if tc.wantErr != nil {
				if err == nil {
					t.Error("expected an error but got none")
				} else if !errors.Is(err, tc.wantErr) {
					t.Errorf("want error %q, got %q", tc.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error but got %q", err)
				return
			}

			if submission.Title != tc.want.Title {
				t.Errorf("want name %q, got %q", tc.want.Title, submission.Title)
			}
		})
	}
}

func TestServiceSearch(t *testing.T) {
	testCases := []struct {
		tname                 string
		repositorySubreddits  []*Subreddit
		repositorySubmissions []*Submission
		text                  string
		want                  []*Submission
		wantErr               error
	}{
		// nominal cases
		{
			tname: "single result",
			repositorySubreddits: []*Subreddit{
				{ID: 1, Name: "astrophotography"},
			},
			repositorySubmissions: []*Submission{
				{ID: 1, PostID: "m31aga", Subreddit: &Subreddit{ID: 1}, Title: "Messier 31 - The Andromeda Galaxy"},
				{ID: 2, PostID: "owlsrf", Subreddit: &Subreddit{ID: 1}, Title: "The Owl Nebula and Surfboard Galaxy"},
			},
			text: "nebula",
			want: []*Submission{
				{
					ID:        2,
					PostID:    "owlsrf",
					Title:     "The Owl Nebula and Surfboard Galaxy",
					Subreddit: &Subreddit{ID: 1, Name: "astrophotography"},
				},
			},
		},
		{
			tname: "multiple results",
			repositorySubreddits: []*Subreddit{
				{ID: 1, Name: "astrophotography"},
			},
			repositorySubmissions: []*Submission{
				{ID: 1, PostID: "m31aga", Subreddit: &Subreddit{ID: 1}, Title: "Messier 31 - The Andromeda Galaxy"},
				{ID: 2, PostID: "owlsrf", Subreddit: &Subreddit{ID: 1}, Title: "The Owl Nebula and Surfboard Galaxy"},
			},
			text: "galaxy",
			want: []*Submission{
				{
					ID:        1,
					PostID:    "m31aga",
					Title:     "Messier 31 - The Andromeda Galaxy",
					Subreddit: &Subreddit{ID: 1, Name: "astrophotography"},
				},
				{
					ID:        2,
					PostID:    "owlsrf",
					Title:     "The Owl Nebula and Surfboard Galaxy",
					Subreddit: &Subreddit{ID: 1, Name: "astrophotography"},
				},
			},
		},

		// error cases
		{
			tname:   "empty text",
			text:    "",
			wantErr: ErrSubmissionSearchTextEmpty,
		},
		{
			tname:   "whitespace (empty) text",
			text:    "       ",
			wantErr: ErrSubmissionSearchTextEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tc.repositorySubmissions, tc.repositorySubreddits)
			service := NewService(repository)

			submissions, err := service.Search(tc.text)

			if tc.wantErr != nil {
				if err == nil {
					t.Error("expected an errot but got none")
				} else if !errors.Is(err, tc.wantErr) {
					t.Errorf("want error %q, got %q", tc.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error, got %q", err)
				return
			}

			if len(submissions) != len(tc.want) {
				t.Errorf("want %d submissions, got %d", len(tc.want), len(submissions))
				return
			}

			for index, want := range tc.want {
				assertSubmissionEquals(t, want, submissions[index])
				assertSubmissionSubredditEquals(t, want, submissions[index])
			}
		})
	}
}

func TestServiceSubredditByID(t *testing.T) {
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
			wantErr: ErrSubredditNotFound,
		},

		// error cases
		{
			tname:   "ID is negative",
			id:      -548,
			wantErr: ErrSubredditIDInvalid,
		},
		{
			tname:   "ID equals zero",
			id:      0,
			wantErr: ErrSubredditIDInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(nil, tc.repositorySubreddits)
			service := NewService(repository)

			subreddit, err := service.subredditByID(tc.id)

			if tc.wantErr != nil {
				if err == nil {
					t.Error("expected an error but got none")
				} else if !errors.Is(err, tc.wantErr) {
					t.Errorf("want error %q, got %q", tc.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error but got %q", err)
				return
			}

			if subreddit.ID != tc.want.ID {
				t.Errorf("want ID %d, got %d", tc.want.ID, subreddit.ID)
			}
			if subreddit.Name != tc.want.Name {
				t.Errorf("want name %q, got %q", tc.want.Name, subreddit.Name)
			}
		})
	}
}

func TestServiceSubredditByName(t *testing.T) {
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
			wantErr: ErrSubredditNotFound,
		},

		// error cases
		{
			tname:   "empty name",
			name:    "",
			wantErr: ErrSubredditNameEmpty,
		},
		{
			tname:   "empty name (whitespace)",
			name:    "     ",
			wantErr: ErrSubredditNameEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(nil, tc.repositorySubreddits)
			validator := NewService(repository)

			subreddit, err := validator.SubredditByName(tc.name)

			if tc.wantErr != nil {
				if err == nil {
					t.Error("expected an error but got none")
				} else if !errors.Is(err, tc.wantErr) {
					t.Errorf("want error %q, got %q", tc.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error but got %q", err)
				return
			}

			if subreddit.ID != tc.want.ID {
				t.Errorf("want ID %d, got %d", tc.want.ID, subreddit.ID)
			}
			if subreddit.Name != tc.want.Name {
				t.Errorf("want name %q, got %q", tc.want.Name, subreddit.Name)
			}
		})
	}
}

func TestServiceSubredditCreate(t *testing.T) {
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
				{ID: 1, Name: "FromSpaceWithLove"},
			},
			subreddit: &Subreddit{Name: "FromSpaceWithLove"},
			wantErr:   ErrSubredditNameAlreadyRegistered,
		},

		// error cases
		{
			tname:     "empty name",
			subreddit: &Subreddit{Name: ""},
			wantErr:   ErrSubredditNameEmpty,
		},
		{
			tname:     "empty name (whitespace)",
			subreddit: &Subreddit{Name: "   "},
			wantErr:   ErrSubredditNameEmpty,
		},
		{
			tname:     "non-default ID",
			subreddit: &Subreddit{ID: 12, Name: "NonZero"},
			wantErr:   ErrSubredditIDInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(nil, tc.repositorySubreddits)
			currentID := repository.subredditCurrentID
			service := NewService(repository)

			err := service.SubredditCreate(tc.subreddit)

			if tc.wantErr != nil {
				if err == nil {
					t.Error("expected an error but got none")
				} else if !errors.Is(err, tc.wantErr) {
					t.Errorf("want error %q, got %q", tc.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error but got %q", err)
			}

			subreddit, err := service.subredditByID(currentID)
			if err != nil {
				t.Errorf("failed to retrieve subreddit: %q", err)
			}

			if subreddit.ID != currentID {
				t.Errorf("want ID %d, got %d", currentID, subreddit.ID)
			}
			if subreddit.Name != tc.subreddit.Name {
				t.Errorf("want name %q, got %q", tc.subreddit.Name, subreddit.Name)
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
	if got.Subreddit.ID != want.Subreddit.ID {
		t.Errorf("want subreddit ID %d, got %d", want.Subreddit.ID, got.Subreddit.ID)
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
