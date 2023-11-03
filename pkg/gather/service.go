package gather

import (
	"context"
	"errors"
	"image"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/sethjones/go-reddit/v2/reddit"
	"github.com/sourcegraph/conc/pool"
	"github.com/virtualtam/walric/pkg/submission"
)

const (
	nWorkers = 4
)

// Service handles domain operations for gathering image files from Reddit.
type Service struct {
	logger zerolog.Logger

	client            *reddit.Client
	submissionService *submission.Service
	dataDir           string
	listPostOptions   *reddit.ListPostOptions
}

// NewService creates and initializes a new Service.
func NewService(rootLogger zerolog.Logger, client *reddit.Client, submissionService *submission.Service, dataDir string, listPostOptions *reddit.ListPostOptions) *Service {
	return &Service{
		logger: rootLogger.With().Str("service", "gather").Logger(),

		client:            client,
		submissionService: submissionService,
		dataDir:           dataDir,
		listPostOptions:   listPostOptions,
	}
}

func (s *Service) filterPosts(posts []*reddit.Post) ([]*reddit.Post, error) {
	var imagePosts []*reddit.Post

	for _, post := range posts {
		postLogger := s.logger.With().
			Str("post_id", post.ID).
			Str("post_title", post.Title).
			Str("subreddit", post.SubredditName).
			Logger()

		// check whether the post's URL is likely to point to an image file
		mediaURL, err := url.Parse(post.URL)
		if err != nil {
			postLogger.Error().
				Err(err).
				Str("post_url", post.URL).
				Msg("failed to parse URL")
			continue
		}

		if !maybeImageURL(mediaURL) {
			postLogger.Debug().Msg("submission does not contain an image")
			continue
		}

		// check whether the post was already saved
		_, err = s.submissionService.ByPostID(post.ID)

		if err == nil {
			postLogger.Debug().Msg("submission already saved")
			continue
		}

		if err != submission.ErrSubmissionNotFound {
			postLogger.Error().Err(err).Msg("database: failed to query submission information")
			return []*reddit.Post{}, err
		}

		// perform a HTTP HEAD request to ensure the URL points to a supported
		// image file
		ok, err := isSupportedImageURL(http.DefaultClient, mediaURL)

		if err != nil {
			postLogger.Error().
				Err(err).
				Str("post_url", post.URL).
				Msg("failed to retrieve remote file metadata")
			continue
		}

		if !ok {
			postLogger.Debug().Msg("unsupported image file format")
			continue
		}

		imagePosts = append(imagePosts, post)
	}

	return imagePosts, nil
}

func (s *Service) gatherImageSubmission(sr *submission.Subreddit, subredditName string, subredditDir string, post *reddit.Post) error {
	gatherLogger := s.logger.With().Str("subreddit", subredditName).Logger()

	postImage, err := newPostImage(subredditDir, post)
	if err != nil {
		gatherLogger.Error().
			Err(err).
			Str("post_url", post.URL).
			Msg("failed to fetch image metadata")
		return err
	}

	if err := postImage.Download(); err != nil {
		gatherLogger.Error().
			Err(err).
			Str("post_url", post.URL).
			Msgf("failed to download image")
		return err
	}

	err = postImage.GetResolutionFromFile()
	if errors.Is(err, image.ErrFormat) {
		gatherLogger.Warn().
			Str("filepath", postImage.filePath).
			Msgf("unknown or unsupported image file format")

		if err := os.Remove(postImage.filePath); err != nil {
			gatherLogger.Error().
				Err(err).
				Str("filepath", postImage.filePath).
				Msg("failed to remove unsupported image file")
			return err
		}

		gatherLogger.Warn().
			Str("filepath", postImage.filePath).
			Msg("unsupported image file removed")

		return nil
	} else if err != nil {
		gatherLogger.Error().
			Err(err).
			Str("filepath", postImage.filePath).
			Msg("failed to get image resolution")
		return err
	}

	imageURL, err := url.Parse(post.URL)
	if err != nil {
		gatherLogger.Error().
			Err(err).
			Stringer("image_url", imageURL).
			Msg("failed to parse image URL")
		return err
	}

	dbSubmission := &submission.Submission{
		Subreddit:     sr,
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
		gatherLogger.Error().
			Err(err).
			Str("post_id", post.ID).
			Str("post_title", post.Title).
			Msgf("failed to create submission")
		return err
	}

	gatherLogger.Info().
		Str("post_id", post.ID).
		Str("post_title", post.Title).
		Msg("submission saved to database")
	return nil
}

func (s *Service) gatherImageSubmissions(ctx context.Context, subredditName string, posts []*reddit.Post) error {
	gatherLogger := s.logger.With().Str("subreddit", subredditName).Logger()

	subredditDir := filepath.Join(s.dataDir, subredditName)
	if err := os.MkdirAll(subredditDir, os.ModePerm); err != nil {
		gatherLogger.Error().
			Err(err).
			Str("subreddit_dir", subredditDir).
			Msg("failed to create directory")
		return err
	}

	sr, err := s.submissionService.SubredditGetOrCreateByName(subredditName)
	if err != nil {
		gatherLogger.Error().
			Err(err).
			Msg("failed to query database")
	}

	workerPool := pool.New().WithErrors().WithMaxGoroutines(nWorkers)
	for _, post := range posts {
		workerPost := post
		workerPool.Go(func() error {
			return s.gatherImageSubmission(sr, subredditName, subredditDir, workerPost)
		})
	}
	if err := workerPool.Wait(); err != nil {
		gatherLogger.Error().
			Err(err).
			Msg("failed to download some submissions")
	}

	return nil
}

// GatherTopImageSubmissions gathers images for the top N submissions for the
// configured subreddits.
func (s *Service) GatherTopImageSubmissions(ctx context.Context, subredditNames []string) error {
	s.logger.Info().
		Int("gather_limit", s.listPostOptions.Limit).
		Str("gather_range", s.listPostOptions.Time).
		Msg("gathering Reddit posts containing images")

	for _, subredditName := range subredditNames {
		gatherLogger := s.logger.With().Str("subreddit", subredditName).Logger()

		topPosts, _, err := s.client.Subreddit.TopPosts(
			ctx,
			subredditName,
			s.listPostOptions,
		)

		if err != nil {
			gatherLogger.Error().
				Err(err).
				Msg("failed to retrieve posts")
			return err
		}

		gatherLogger.Debug().
			Int("n_posts", len(topPosts)).
			Msg("found top posts")

		posts, err := s.filterPosts(topPosts)
		if err != nil {
			gatherLogger.Error().Err(err).Msg("failed to filter posts")
			return err
		}

		if len(posts) == 0 {
			gatherLogger.Info().Msgf("found no new posts or no post containing images")
			continue
		}

		gatherLogger.Info().
			Int("n_posts", len(topPosts)).
			Int("n_image_posts", len(posts)).
			Msg("found new posts containing images")

		if err := s.gatherImageSubmissions(ctx, subredditName, posts); err != nil {
			return err
		}
	}

	return nil
}
