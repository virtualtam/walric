package gather

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Gatherer struct {
	client  *reddit.Client
	dataDir string
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

		for _, post := range posts {
			fmt.Println(post.SubredditName, post.Title)

			// TODO download file
			// TODO save to database
		}
	}

	return nil
}

func NewGatherer(client *reddit.Client, dataDir string) *Gatherer {
	return &Gatherer{
		client:  client,
		dataDir: dataDir,
	}
}
