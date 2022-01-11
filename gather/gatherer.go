package gather

import (
	"context"
	"errors"
	"image"
	"net/url"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/virtualtam/redwall2/submission"
	"github.com/virtualtam/redwall2/subreddit"
)

type Gatherer struct {
	client            *reddit.Client
	submissionService *submission.Service
	subredditService  *subreddit.Service
	dataDir           string
}

func (g *Gatherer) GatherTopImageSubmissions(ctx context.Context, subredditNames []string, listPostOptions *reddit.ListPostOptions) error {
	for _, subredditName := range subredditNames {
		topPosts, _, err := g.client.Subreddit.TopPosts(
			ctx,
			subredditName,
			listPostOptions,
		)

		if err != nil {
			log.Error().Err(err).Msgf(
				"%s: failed to retrieve top %d posts for the last %s",
				subredditName,
				listPostOptions.Limit,
				listPostOptions.Time,
			)
			return err
		}

		log.Debug().Msgf(
			"%s: found %d top posts for the last %s",
			subredditName,
			len(topPosts),
			listPostOptions.Time,
		)

		posts := filterImagePosts(topPosts)

		if len(posts) == 0 {
			log.Warn().Msgf("%s: found no posts containing images", subredditName)
			continue
		}

		log.Info().Msgf(
			"%s: found %d top posts containing images for the last %s",
			subredditName,
			len(posts),
			listPostOptions.Time,
		)

		subredditDir := filepath.Join(g.dataDir, subredditName)
		if err := os.MkdirAll(subredditDir, os.ModePerm); err != nil {
			log.Error().Err(err).Msgf("failed to create directory: %s", subredditDir)
			return err
		}

		dbSubreddit, err := g.subredditService.ByName(subredditName)

		if errors.Is(err, subreddit.ErrNotFound) {
			log.Info().Msgf("%s: save subreddit information to database", subredditName)

			dbSubreddit = &subreddit.Subreddit{Name: subredditName}
			if err = g.subredditService.Create(dbSubreddit); err != nil {
				log.Error().Err(err).Msgf("%s: failed to create database entry", subredditName)
				return err
			}

			dbSubreddit, err = g.subredditService.ByName(subredditName)
			if err != nil {
				log.Error().Err(err).Msg("database: failed to retrieve subreddit")
				return err
			}
		} else if err != nil {
			log.Error().Err(err).Msg("database: failed to query subreddit information")
			return err
		}

		for _, post := range posts {
			_, err := g.submissionService.ByPostID(post.ID)

			if err == nil {
				log.Info().Msgf(
					"%s: submission already saved: %s - %s",
					subredditName,
					post.ID,
					post.Title,
				)
				continue
			}

			if err != submission.ErrNotFound {
				log.Error().Err(err).Msgf("database: failed to query submission information")
				return err
			}

			postImage, err := newPostImage(subredditDir, post)
			if err != nil {
				log.Error().Err(err).Msgf("%s: failed to fetch image metadata from URL: %s", subredditName, post.URL)
				return err
			}

			if err := postImage.Download(); err != nil {
				log.Error().Err(err).Msgf("%s: failed to download image from URL: %s", subredditName, post.URL)
				return err
			}

			err = postImage.UpdateResolution()
			if err == image.ErrFormat {
				log.Warn().Msgf("%s: unknown or unsupported image file format: %s", subredditName, postImage.filePath)

				if err := os.Remove(postImage.filePath); err != nil {
					log.Error().Err(err).Msgf("image: failed to remove file: %s", postImage.filePath)
					return err
				}

				log.Warn().Msgf("%s: file removed: %s", subredditName, postImage.filePath)

				continue
			}
			if err != nil {
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

			if err := g.submissionService.Create(dbSubmission); err != nil {
				log.Error().Err(err).Msgf("%s / %s: failed to create submission entry", subredditName, post.ID)
				return err
			}

			log.Info().Msgf("%s: submission saved to database: %s - %s", subredditName, post.ID, post.Title)
		}
	}

	return nil
}

func NewGatherer(client *reddit.Client, submissionService *submission.Service, subredditService *subreddit.Service, dataDir string) *Gatherer {
	return &Gatherer{
		client:            client,
		submissionService: submissionService,
		subredditService:  subredditService,
		dataDir:           dataDir,
	}
}
