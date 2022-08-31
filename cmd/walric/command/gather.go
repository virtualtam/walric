package command

import (
	"context"

	"github.com/sethjones/go-reddit/v2/reddit"
	"github.com/spf13/cobra"
	"github.com/virtualtam/walric/pkg/gather"
)

// NewGatherCommand initializes a CLI command to gather top submissions from the
// configured Subreddits.
func NewGatherCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gather",
		Short: "Gather media from top Reddit submissions",
		Run: func(cmd *cobra.Command, args []string) {
			redditClient, err := reddit.NewClient(
				reddit.Credentials{
					ID:     walricConfig.Reddit.ClientID,
					Secret: walricConfig.Reddit.ClientSecret,
				},
				reddit.WithApplicationOnlyOAuth(true),
				reddit.WithUserAgent(walricConfig.Reddit.UserAgent),
			)
			if err != nil {
				cobra.CheckErr(err)
			}

			gatherer := gather.NewGatherer(redditClient, submissionService, subredditService, walricConfig.Walric.DataDir)

			listPostOptions := &reddit.ListPostOptions{
				ListOptions: reddit.ListOptions{Limit: walricConfig.Walric.SubmissionLimit},
				Time:        walricConfig.Walric.TimeFilter,
			}

			ctx := context.Background()

			err = gatherer.GatherTopImageSubmissions(ctx, walricConfig.Walric.Subreddits, listPostOptions)
			if err != nil {
				cobra.CheckErr(err)
			}
		},
	}

	return cmd
}
