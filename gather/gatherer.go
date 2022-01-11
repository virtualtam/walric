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

		posts = filterImagePosts(posts)

		// TODO create directory

		for _, post := range posts {
			fmt.Println(post.SubredditName, post.Title)

			// TODO download file
			// TODO save to database
		}
	}

	return nil
}

func NewGatherer(client *reddit.Client) *Gatherer {
	return &Gatherer{
		client: client,
	}
}
