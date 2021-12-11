package command

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// NewHistoryCommand initializes a CLI command to display the history of the
// selected entries.
func NewHistoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Display the history of selected entries",
		Run: func(cmd *cobra.Command, args []string) {
			history, err := historyService.All()
			if err != nil {
				cobra.CheckErr(err)
			}

			writer := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)

			for _, entry := range history {
				fmt.Fprintf(
					writer,
					"%s\t%s\t%d x %d\t%s\n",
					entry.Date.UTC().Format("2006-01-02 15:04:05 MST"),
					entry.Submission.PostID,
					entry.Submission.ImageWidthPx,
					entry.Submission.ImageHeightPx,
					entry.Submission.Title,
				)
			}

			writer.Flush()
		},
	}

	return cmd
}
