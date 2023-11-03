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

			listPostOptions := &reddit.ListPostOptions{
				ListOptions: reddit.ListOptions{Limit: walricConfig.Walric.SubmissionLimit},
				Time:        walricConfig.Walric.TimeFilter,
			}

			gatherService := gather.NewService(redditClient, submissionService, walricConfig.Walric.DataDir, listPostOptions)

			ctx := context.Background()

			err = gatherService.GatherTopImageSubmissions(ctx, walricConfig.Walric.Subreddits)
			if err != nil {
				cobra.CheckErr(err)
			}
		},
	}

	return cmd
}
