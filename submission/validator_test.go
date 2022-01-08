package submission

import (
	"errors"
	"testing"

	"github.com/virtualtam/redwall2/monitor"
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

	for _, tt := range testCases {
		t.Run(tt.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tt.repositorySubmissions)
			validator := newValidator(repository)

			submission, err := validator.ByID(tt.id)

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

			if submission.ID != tt.want.ID {
				t.Errorf("want ID %d, got %d", tt.want.ID, submission.ID)
			}
			if submission.Title != tt.want.Title {
				t.Errorf("want name %q, got %q", tt.want.Title, submission.Title)
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

	for _, tt := range testCases {
		t.Run(tt.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tt.repositorySubmissions)
			validator := newValidator(repository)

			submission, err := validator.ByPostID(tt.postID)

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

			if submission.ID != tt.want.ID {
				t.Errorf("want ID %d, got %d", tt.want.ID, submission.ID)
			}
			if submission.Title != tt.want.Title {
				t.Errorf("want name %q, got %q", tt.want.Title, submission.Title)
			}
		})
	}
}

func TestValidatorByTitle(t *testing.T) {
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

	for _, tt := range testCases {
		t.Run(tt.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tt.repositorySubmissions)
			validator := newValidator(repository)

			submissions, err := validator.ByTitle(tt.text)

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

			if len(submissions) != len(tt.want) {
				t.Errorf("want %d submissions, got %d", len(tt.want), len(submissions))
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

	for _, tt := range testCases {
		t.Run(tt.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tt.repositorySubmissions)
			validator := newValidator(repository)

			submissions, err := validator.ByMinResolution(tt.minResolution)

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

			if len(submissions) != len(tt.want) {
				t.Errorf("want %d submissions, got %d", len(tt.want), len(submissions))
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

	for _, tt := range testCases {
		t.Run(tt.tname, func(t *testing.T) {
			repository := NewRepositoryInMemory(tt.repositorySubmissions)
			validator := newValidator(repository)

			submission, err := validator.Random(tt.minResolution)

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

			if submission.Title != tt.want.Title {
				t.Errorf("want name %q, got %q", tt.want.Title, submission.Title)
			}
		})
	}
}
