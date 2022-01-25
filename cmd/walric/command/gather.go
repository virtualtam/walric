package command

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/virtualtam/walric/gather"
)

// NewGatherCommand initializes a CLI command to gather top submissions from the
// configured Subreddits.
func NewGatherCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gather",
		Short: "Gather media from top Reddit submissions",
		Run: func(cmd *cobra.Command, args []string) {
			redditReadOnlyClient, err := reddit.NewReadonlyClient(
				reddit.WithUserAgent(walricConfig.Reddit.UserAgent),
			)
			if err != nil {
				cobra.CheckErr(err)
			}

			gatherer := gather.NewGatherer(redditReadOnlyClient, submissionService, subredditService, walricConfig.Walric.DataDir)

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
