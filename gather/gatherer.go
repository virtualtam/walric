package gather

import (
	"context"
	"fmt"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Gatherer struct {
	client *reddit.Client
}

func (g *Gatherer) GatherTopImageSubmissions(ctx context.Context, subreddits []string, listPostOptions *reddit.ListPostOptions) error {
	for _, subreddit := range subreddits {
		posts, _, err := g.client.Subreddit.TopPosts(
			ctx,
			subreddit,
			listPostOptions,
		)

		if err != nil {
			return err
		}

		for _, post := range posts {
			fmt.Println(post.SubredditName, post.Title)
		}
	}

	return nil
}

func NewGatherer(client *reddit.Client) *Gatherer {
	return &Gatherer{
		client: client,
	}
}
