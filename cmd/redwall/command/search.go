package command

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// NewSearchCommand initializes a CLI command to search Submissions.
func NewSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for submissions by title",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			submissions, err := submissionService.Search(args[0])
			if err != nil {
				cobra.CheckErr(err)
			}

			writer := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)

			for _, submission := range submissions {
				fmt.Fprintf(
					writer,
					"%s\t%s\t%d x %d\t%s\n",
					submission.Subreddit.Name,
					submission.PostID,
					submission.ImageWidthPx,
					submission.ImageHeightPx,
					submission.Title,
				)
			}

			writer.Flush()

			fmt.Println()
			fmt.Println(len(submissions), "submission(s) found")
		},
	}

	return cmd
}
