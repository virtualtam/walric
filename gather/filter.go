package gather

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func filterImagePosts(posts []*reddit.Post) []*reddit.Post {
	var imagePosts []*reddit.Post

	for _, post := range posts {
		mediaURL, err := url.Parse(post.URL)
		if err != nil {
			continue
		}

		if maybeImageURL(mediaURL) {
			imagePosts = append(imagePosts, post)
		}
	}

	return imagePosts
}

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

	// 2. check file extension
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
