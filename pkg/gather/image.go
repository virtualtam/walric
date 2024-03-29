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

	"github.com/sethjones/go-reddit/v2/reddit"
)

type postImage struct {
	url      string
	filePath string

	HeightPx int
	WidthPx  int
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

// Download downloads an image locally.
func (i *postImage) Download() error {
	resp, err := http.Get(i.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: %s", resp.Status)
	}

	out, err := os.Create(i.filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// GetResolutionFromFile retrieves an image's resolution (height, width).
func (i *postImage) GetResolutionFromFile() error {
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

// isSupportedImageURL performs a HTTP HEAD request to retrieve the Content-Type
// header for the remote file, and determine whether the type of the remote file
// is a  supported image format.
func isSupportedImageURL(client *http.Client, mediaURL *url.URL) (bool, error) {
	response, err := client.Head(mediaURL.String())

	if err != nil {
		return false, err
	}

	contentType := response.Header.Get("Content-Type")

	switch contentType {
	case "application/octet-stream", "image/jpeg", "image/png":
		return true, nil
	}

	return false, nil
}
