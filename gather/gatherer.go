package gather

import (
	"context"
	"fmt"
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

		fmt.Println(dbSubreddit)

		for _, post := range posts {
			fmt.Println(post.SubredditName, post.Title)

			// TODO download file
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
