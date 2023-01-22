package gather

import (
	"context"
	"errors"
	"image"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/sethjones/go-reddit/v2/reddit"
	"github.com/sourcegraph/conc/pool"
	"github.com/virtualtam/walric/pkg/submission"
	"github.com/virtualtam/walric/pkg/subreddit"
)

const (
	nWorkers = 4
)

// Service handles domain operations for gathering image files from Reddit.
type Service struct {
	client            *reddit.Client
	submissionService *submission.Service
	subredditService  *subreddit.Service
	dataDir           string
	listPostOptions   *reddit.ListPostOptions
}

// NewService creates and initializes a new Service.
func NewService(client *reddit.Client, submissionService *submission.Service, subredditService *subreddit.Service, dataDir string, listPostOptions *reddit.ListPostOptions) *Service {
	return &Service{
		client:            client,
		submissionService: submissionService,
		subredditService:  subredditService,
		dataDir:           dataDir,
		listPostOptions:   listPostOptions,
	}
}

func (s *Service) filterPosts(posts []*reddit.Post) ([]*reddit.Post, error) {
	var imagePosts []*reddit.Post

	for _, post := range posts {
		// check whether the post's URL is likely to point to an image file
		mediaURL, err := url.Parse(post.URL)
		if err != nil {
			log.Error().Err(err).Msgf("failed to parse URL: %s", post.URL)
			continue
		}

		if !maybeImageURL(mediaURL) {
			log.Debug().Msgf(
				"%s: submission does not contain an image: %s - %s",
				post.SubredditName,
				post.ID,
				post.Title,
			)
			continue
		}

		// check whether the post was already saved
		_, err = s.submissionService.ByPostID(post.ID)

		if err == nil {
			log.Debug().Msgf(
				"%s: submission already saved: %s - %s",
				post.SubredditName,
				post.ID,
				post.Title,
			)
			continue
		}

		if err != submission.ErrNotFound {
			log.Error().Err(err).Msgf("database: failed to query submission information")
			return []*reddit.Post{}, err
		}

		// perform a HTTP HEAD request to ensure the URL points to a supported
		// image file
		ok, err := isSupportedImageURL(http.DefaultClient, mediaURL)

		if err != nil {
			log.Error().Err(err).Msgf("failed to retrieve remote file metadata: %s", post.URL)
			continue
		}

		if !ok {
			log.Debug().Msgf(
				"%s: submission points to a file with an unsupported format: %s - %s",
				post.SubredditName,
				post.ID,
				post.Title,
			)
			continue
		}

		imagePosts = append(imagePosts, post)
	}

	return imagePosts, nil
}

func (s *Service) gatherImageSubmission(dbSubreddit *subreddit.Subreddit, subredditName string, subredditDir string, post *reddit.Post) error {
	postImage, err := newPostImage(subredditDir, post)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to fetch image metadata from URL: %s", subredditName, post.URL)
		return err
	}

	if err := postImage.Download(); err != nil {
		log.Error().Err(err).Msgf("%s: failed to download image from URL: %s", subredditName, post.URL)
		return err
	}

	err = postImage.GetResolutionFromFile()
	if errors.Is(err, image.ErrFormat) {
		log.Warn().Msgf("%s: unknown or unsupported image file format: %s", subredditName, postImage.filePath)

		if err := os.Remove(postImage.filePath); err != nil {
			log.Error().Err(err).Msgf("image: failed to remove file: %s", postImage.filePath)
			return err
		}

		log.Warn().Msgf("%s: file removed: %s", subredditName, postImage.filePath)

		return nil
	} else if err != nil {
		log.Error().Err(err).Msgf("%s: failed to get image resolution: %s", subredditName, postImage.filePath)
		return err
	}

	imageURL, err := url.Parse(post.URL)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to parse image URL: %s", subredditName, imageURL)
		return err
	}

	dbSubmission := &submission.Submission{
		Subreddit:     dbSubreddit,
		SubredditID:   dbSubreddit.ID,
		Author:        post.Author,
		Permalink:     post.Permalink,
		PostID:        post.ID,
		PostedAt:      post.Created.UTC(),
		Score:         post.Score,
		Title:         post.Title,
		ImageDomain:   imageURL.Host,
		ImageURL:      post.URL,
		ImageNSFW:     post.NSFW,
		ImageFilename: postImage.filePath,
		ImageHeightPx: postImage.HeightPx,
		ImageWidthPx:  postImage.WidthPx,
	}

	if err := s.submissionService.Create(dbSubmission); err != nil {
		log.Error().Err(err).Msgf("%s: failed to create submission entry: %s - %s", subredditName, post.ID, post.Title)
		return err
	}

	log.Info().Msgf("%s: submission saved to database: %s - %s", subredditName, post.ID, post.Title)
	return nil
}

func (s *Service) gatherImageSubmissions(ctx context.Context, subredditName string, posts []*reddit.Post) error {
	subredditDir := filepath.Join(s.dataDir, subredditName)
	if err := os.MkdirAll(subredditDir, os.ModePerm); err != nil {
		log.Error().Err(err).Msgf("failed to create directory: %s", subredditDir)
		return err
	}

	dbSubreddit, err := s.subredditService.GetOrCreateByName(subredditName)
	if err != nil {
		log.Error().Err(err).Msgf("%s: failed to query database", subredditName)
	}

	workerPool := pool.New().WithErrors().WithMaxGoroutines(nWorkers)
	for _, post := range posts {
		workerPost := post
		workerPool.Go(func() error {
			return s.gatherImageSubmission(dbSubreddit, subredditName, subredditDir, workerPost)
		})
	}
	if err := workerPool.Wait(); err != nil {
		log.Error().Err(err).Msgf("%s: failed to download some submissions", subredditName)
	}

	return nil
}

// GatherTopImageSubmissions gathers images for the top N submissions for the
// configured subreddits.
func (s *Service) GatherTopImageSubmissions(ctx context.Context, subredditNames []string) error {
	for _, subredditName := range subredditNames {
		topPosts, _, err := s.client.Subreddit.TopPosts(
			ctx,
			subredditName,
			s.listPostOptions,
		)

		if err != nil {
			log.Error().Err(err).Msgf(
				"%s: failed to retrieve top %d posts for the last %s",
				subredditName,
				s.listPostOptions.Limit,
				s.listPostOptions.Time,
			)
			return err
		}

		log.Debug().Msgf(
			"%s: found %d top posts for the last %s",
			subredditName,
			len(topPosts),
			s.listPostOptions.Time,
		)

		posts, err := s.filterPosts(topPosts)
		if err != nil {
			log.Error().Err(err).Msgf("%s: failed to filter posts", subredditName)
			return err
		}

		if len(posts) == 0 {
			log.Warn().Msgf("%s: found no new posts, or no post containing images", subredditName)
			continue
		}

		log.Info().Msgf(
			"%s: found %d new posts containing images for the last %s",
			subredditName,
			len(posts),
			s.listPostOptions.Time,
		)

		if err := s.gatherImageSubmissions(ctx, subredditName, posts); err != nil {
			return err
		}
	}

	return nil
}
