package gather

import (
	"context"
	"fmt"
	"image"
	"os"
	"path/filepath"

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
		posts, _, err := g.client.Subreddit.TopPosts(
			ctx,
			subredditName,
			listPostOptions,
		)

		if err != nil {
			return err
		}

		posts = filterImagePosts(posts)

		if len(posts) == 0 {
			continue
		}

		subredditDir := filepath.Join(g.dataDir, subredditName)
		if err := os.MkdirAll(subredditDir, os.ModePerm); err != nil {
			return err
		}

		dbSubreddit, err := g.subredditService.ByName(subredditName)
		if err == subreddit.ErrNotFound {
			dbSubreddit = &subreddit.Subreddit{Name: subredditName}
			if err = g.subredditService.Create(dbSubreddit); err != nil {
				return err
			}

			dbSubreddit, err = g.subredditService.ByName(subredditName)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		for _, post := range posts {
			_, err := g.submissionService.ByPostID(post.ID)

			if err != submission.ErrNotFound {
				return err
			}

			postImage, err := newPostImage(subredditDir, post)
			if err != nil {
				return err
			}

			if err := postImage.Download(); err != nil {
				return err
			}

			err = postImage.UpdateResolution()
			if err == image.ErrFormat {
				if err := os.Remove(postImage.filePath); err != nil {
					return err
				}

				continue
			}
			if err != nil {
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
				ImageURL:      post.URL,
				ImageNSFW:     post.NSFW,
				ImageFilename: postImage.filePath,
				ImageHeightPx: postImage.HeightPx,
				ImageWidthPx:  postImage.WidthPx,
			}

			fmt.Printf("%#v\n", dbSubmission)

			// TODO save to database
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
