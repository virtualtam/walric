package gather

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type postImage struct {
	url      string
	filePath string

	HeightPx int
	WidthPx  int
}

func (i *postImage) Download() error {
	out, err := os.Create(i.filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(i.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (i *postImage) UpdateResolution() error {
	reader, err := os.Open(i.filePath)
	if err != nil {
		return err
	}

	defer reader.Close()

	config, _, err := image.DecodeConfig(reader)
	if err != nil {
		return err
	}

	i.HeightPx = config.Height
	i.WidthPx = config.Width

	return nil
}

func newPostImage(subredditDir string, post *reddit.Post) (*postImage, error) {
	imageURL, err := url.Parse(post.URL)
	if err != nil {
		return &postImage{}, err
	}

	fileName := fmt.Sprintf("%s-%s", post.ID, filepath.Base(imageURL.Path))
	filePath := filepath.Join(subredditDir, fileName)

	return &postImage{
		url:      post.URL,
		filePath: filePath,
	}, nil
}

// maybeImageURL attempts to determine whether a URL points to a JPEG or PNG
// image file, by looking at the URL's host and path. Checks are performed
// locally, and no outgoing request is made.
func maybeImageURL(mediaURL *url.URL) bool {
	// 1. check hosting domain (exact match)
	switch mediaURL.Host {
	case "gfycat.com":
		// GIF hosting
		return false

	case "open.spotify.com":
		// audio hosting
		return false

	case "v.redd.it", "youtu.be":
		// video hosting
		return false
	}

	// 2. check hosting domain and path
	if mediaURL.Host == "www.reddit.com" && strings.HasPrefix(mediaURL.Path, "/gallery") {
		// Reddit image gallery
		return false
	}

	// 3. check file extension
	ext := strings.ToLower(filepath.Ext(filepath.Base(mediaURL.Path)))

	switch ext {
	case ".gif", ".gifv":
		return false

	case ".mp4":
		return false
	}

	// despite the previous guesses, the URL may still point to a non-image
	// file, eg if the URL does not contain a file extension
	return true
}
