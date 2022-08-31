package submission

import (
	"errors"
	"testing"
	"time"

	"github.com/virtualtam/walric/pkg/monitor"
)

func TestValidatorByID(t *testing.T) {
	testCases := []struct {
		tname                 string
		repositorySubmissions []*Submission
		id                    int
		want                  *Submission
		wantErr               error
	}{
		// nominal cases
		{
			tname: "existing ID",
			repositorySubmissions: []*Submission{
				{ID: 1, Title: "Sunday Afternoon In The Park [1981x2015]"},
				{ID: 2, Title: "Moroccan Sunset [1995x800]"},
				{ID: 3, Title: "Laguna Sunrise [1972x408]"},
			},
			id:   3,
			want: &Submission{ID: 3, Title: "Laguna Sunrise [1972x408]"},
		},
		{
			tname:   "not found",
			id:      7000,
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

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tc.repositorySubmissions)
			validator := newValidator(repository)

			submission, err := validator.ByID(tc.id)

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

			if submission.ID != tc.want.ID {
				t.Errorf("want ID %d, got %d", tc.want.ID, submission.ID)
			}
			if submission.Title != tc.want.Title {
				t.Errorf("want name %q, got %q", tc.want.Title, submission.Title)
			}
		})
	}
}

func TestValidatorByPostID(t *testing.T) {
	testCases := []struct {
		tname                 string
		repositorySubmissions []*Submission
		postID                string
		want                  *Submission
		wantErr               error
	}{
		// nominal cases
		{
			tname: "existing ID",
			repositorySubmissions: []*Submission{
				{PostID: "safprk", Title: "Sunday Afternoon In The Park [1981x2015]"},
				{PostID: "morsun", Title: "Moroccan Sunset [1995x800]"},
				{PostID: "lagsun", Title: "Laguna Sunrise [1972x408]"},
			},
			postID: "lagsun",
			want:   &Submission{PostID: "lagsun", Title: "Laguna Sunrise [1972x408]"},
		},
		{
			tname:   "not found",
			postID:  "rnyday",
			wantErr: ErrNotFound,
		},

		// error cases
		{
			tname:   "empty PostID",
			postID:  "",
			wantErr: ErrPostIDEmpty,
		},
		{
			tname:   "whitespace (empty) PostID",
			postID:  "       ",
			wantErr: ErrPostIDEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tc.repositorySubmissions)
			validator := newValidator(repository)

			submission, err := validator.ByPostID(tc.postID)

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

			if submission.ID != tc.want.ID {
				t.Errorf("want ID %d, got %d", tc.want.ID, submission.ID)
			}
			if submission.Title != tc.want.Title {
				t.Errorf("want name %q, got %q", tc.want.Title, submission.Title)
			}
		})
	}
}

func TestValidatorSearch(t *testing.T) {
	testCases := []struct {
		tname                 string
		repositorySubmissions []*Submission
		text                  string
		want                  []*Submission
		wantErr               error
	}{
		// nominal cases
		{
			tname: "single result",
			repositorySubmissions: []*Submission{
				{Title: "Sunday Afternoon In The Park [1981x2015]"},
				{Title: "Moroccan Sunset [1995x800]"},
				{Title: "Laguna Sunrise [1972x408]"},
			},
			text: "SUNSET",
			want: []*Submission{{Title: "Laguna Sunrise [1972x408]"}},
		},
		{
			tname: "multiple results",
			repositorySubmissions: []*Submission{
				{Title: "Sunday Afternoon In The Park [1981x2015]"},
				{Title: "Moroccan Sunset [1995x800]"},
				{Title: "Laguna Sunrise [1972x408]"},
			},
			text: "sun",
			want: []*Submission{
				{Title: "Sunday Afternoon In The Park [1981x2015]"},
				{Title: "Moroccan Sunset [1995x800]"},
				{Title: "Laguna Sunrise [1972x408]"},
			},
		},
		{
			tname: "no results",
			text:  "rain",
		},

		// error cases
		{
			tname:   "empty text",
			text:    "",
			wantErr: ErrSearchTextEmpty,
		},
		{
			tname:   "whitespace (empty) text",
			text:    "       ",
			wantErr: ErrSearchTextEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tc.repositorySubmissions)
			validator := newValidator(repository)

			submissions, err := validator.Search(tc.text)

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

func TestValidatorByMinResolution(t *testing.T) {
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
				{Title: "Sunday Afternoon In The Park [640x480]", ImageHeightPx: 480, ImageWidthPx: 640},
				{Title: "Laguna Sunrise [1920x1200]", ImageHeightPx: 1200, ImageWidthPx: 1920},
			},
			minResolution: &monitor.Resolution{HeightPx: 1200, WidthPx: 1920},
			want: []*Submission{
				{Title: "Laguna Sunrise [1920x1200]", ImageHeightPx: 1200, ImageWidthPx: 1920},
			},
		},
		{
			tname: "multiple results",
			repositorySubmissions: []*Submission{
				{Title: "Sunday Afternoon In The Park [640x480]", ImageHeightPx: 480, ImageWidthPx: 640},
				{Title: "Moroccan Sunset [2560x1440]", ImageHeightPx: 1440, ImageWidthPx: 2560},
				{Title: "Laguna Sunrise [1920x1200]", ImageHeightPx: 1200, ImageWidthPx: 1920},
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
			wantErr:       ErrResolutionInvalid,
		},
		{
			tname:         "negative resolution width",
			minResolution: &monitor.Resolution{HeightPx: 1200, WidthPx: -1920},
			wantErr:       ErrResolutionInvalid,
		},
		{
			tname:         "negative resolution height and width",
			minResolution: &monitor.Resolution{HeightPx: -1200, WidthPx: -1920},
			wantErr:       ErrResolutionInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tc.repositorySubmissions)
			validator := newValidator(repository)

			submissions, err := validator.ByMinResolution(tc.minResolution)

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

func TestValidatorRandom(t *testing.T) {
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
				{Title: "Sunday Afternoon In The Park [640x480]", ImageHeightPx: 480, ImageWidthPx: 640},
				{Title: "Laguna Sunrise [1920x1200]", ImageHeightPx: 1200, ImageWidthPx: 1920},
			},
			minResolution: &monitor.Resolution{HeightPx: 1200, WidthPx: 1920},
			want:          &Submission{Title: "Laguna Sunrise [1920x1200]", ImageHeightPx: 1200, ImageWidthPx: 1920},
		},
		{
			tname:         "not found (empty repository)",
			minResolution: &monitor.Resolution{HeightPx: 1200, WidthPx: 1920},
			wantErr:       ErrNotFound,
		},
		{
			tname: "not found (no result)",
			repositorySubmissions: []*Submission{
				{Title: "Sunday Afternoon In The Park [640x480]", ImageHeightPx: 480, ImageWidthPx: 640},
				{Title: "Laguna Sunrise [1920x1200]", ImageHeightPx: 1200, ImageWidthPx: 1920},
			},
			minResolution: &monitor.Resolution{HeightPx: 1440, WidthPx: 2560},
			wantErr:       ErrNotFound,
		},

		// error cases
		{
			tname:         "negative resolution height",
			minResolution: &monitor.Resolution{HeightPx: -1200, WidthPx: 1920},
			wantErr:       ErrResolutionInvalid,
		},
		{
			tname:         "negative resolution width",
			minResolution: &monitor.Resolution{HeightPx: 1200, WidthPx: -1920},
			wantErr:       ErrResolutionInvalid,
		},
		{
			tname:         "negative resolution height and width",
			minResolution: &monitor.Resolution{HeightPx: -1200, WidthPx: -1920},
			wantErr:       ErrResolutionInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tc.repositorySubmissions)
			validator := newValidator(repository)

			submission, err := validator.Random(tc.minResolution)

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

func TestValidatorCreate(t *testing.T) {
	testCases := []struct {
		tname                 string
		repositorySubmissions []*Submission
		submission            *Submission
		wantErr               error
	}{
		// nominal cases
		{
			tname: "new submission",
			submission: &Submission{
				SubredditID:   25,
				Author:        "janedoe",
				Permalink:     "r/dummy/comments/newnew/new_submission/",
				PostID:        "newnew",
				PostedAt:      time.Now().UTC(),
				Score:         12,
				Title:         "New Submission [800x600]",
				ImageDomain:   "i.redd.it",
				ImageURL:      "https://i.redd.it/newsub001.jpg",
				ImageNSFW:     false,
				ImageFilename: "/deta/walric/dummy/newnew-newsub001.jpg",
				ImageHeightPx: 600,
				ImageWidthPx:  800,
			},
		},

		// error cases
		{
			tname: "negative subreddit ID",
			submission: &Submission{
				SubredditID: -583,
			},
			wantErr: ErrSubredditIDInvalid,
		},
		{
			tname: "subreddit ID equals zero",
			submission: &Submission{
				SubredditID: 0,
			},
			wantErr: ErrSubredditIDInvalid,
		},
		{
			tname: "non-default ID",
			submission: &Submission{
				SubredditID: 12,
				ID:          179,
				PostID:      "nondft",
				Title:       "Non-default [0x0]",
			},
			wantErr: ErrIDInvalid,
		},
		{
			tname: "empty PostID",
			submission: &Submission{
				SubredditID: 12,
			},
			wantErr: ErrPostIDEmpty,
		},
		{
			tname: "empty PostID (whitespace)",
			submission: &Submission{
				SubredditID: 12,
				PostID:      "     ",
			},
			wantErr: ErrPostIDEmpty,
		},
		{
			tname: "duplicate PostID",
			repositorySubmissions: []*Submission{
				{
					SubredditID: 12,
					ID:          1,
					PostID:      "dupdup",
				},
			},
			submission: &Submission{
				SubredditID: 12,
				PostID:      "dupdup",
			},
			wantErr: ErrPostIDAlreadyRegistered,
		},
		{
			tname: "empty title",
			submission: &Submission{
				SubredditID: 12,
				PostID:      "notitl",
			},
			wantErr: ErrTitleEmpty,
		},
		{
			tname: "empty title (whitespace)",
			submission: &Submission{
				SubredditID: 12,
				PostID:      "notitl",
				Title:       "    ",
			},
			wantErr: ErrTitleEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tc.repositorySubmissions)
			currentID := repository.currentID
			validator := newValidator(repository)

			err := validator.Create(tc.submission)

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

			submission, err := validator.ByID(currentID)

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
