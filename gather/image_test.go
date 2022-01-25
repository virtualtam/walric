package gather

import (
	"net/url"
	"testing"
)

func TestMaybeImageURL(t *testing.T) {
	testCases := []struct {
		tname  string
		rawURL string
		want   bool
	}{
		// accepted URLs
		{
			tname:  "image from Reddit",
			rawURL: "https://i.redd.it/9vby1uakau521.jpg",
			want:   true,
		},
		{
			tname:  "image from Imgur (1)",
			rawURL: "https://i.imgur.com/btn0DzA.jpg",
			want:   true,
		},
		{
			tname:  "image from Imgur (2)",
			rawURL: "https://imgur.com/AxcguyH.jpg",
			want:   true,
		},

		// rejected URLs
		{
			tname:  "GIF from gfycat",
			rawURL: "https://gfycat.com/ablegiganticislandwhistler-phyllis-smith-oscar-nunez-creed-bratton",
			want:   false,
		},
		{
			tname:  "audio from Spotify",
			rawURL: "https://open.spotify.com/episode/2i2db3uaCuEiFo6WqcPQGP",
			want:   false,
		},
		{
			tname:  "video from Reddit",
			rawURL: "https://v.redd.it/b3w4hk0bcuy51",
			want:   false,
		},
		{
			tname:  "video from Youtube",
			rawURL: "https://youtu.be/RDYYVGAKqqQ",
			want:   false,
		},
		{
			tname:  "Reddit image gallery",
			rawURL: "https://www.reddit.com/gallery/rk6hzc",
			want:   false,
		},
		{
			tname:  "GIF image",
			rawURL: "https://domain.tld/path/image.gif",
			want:   false,
		},
		{
			tname:  "GIFV image",
			rawURL: "https://domain.tld/path/image.gifv",
			want:   false,
		},
		{
			tname:  "MP4 video",
			rawURL: "https://domain.tld/path/movie.mp4",
			want:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			mediaURL, err := url.Parse(tc.rawURL)

			if err != nil {
				t.Errorf("failed to parse URL: %q", err)
				return
			}

			got := maybeImageURL(mediaURL)

			if got != tc.want {
				t.Errorf("want %t, got %t", tc.want, got)
			}
		})
	}
}
