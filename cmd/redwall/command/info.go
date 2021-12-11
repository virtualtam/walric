package command

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// NewInfoCommand initializes a CLI command to display information on a given
// Submission.
func NewInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info [POST_ID]",
		Short: "Display information about a given submission",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			submission, err := submissionService.ByPostID(args[0])
			if err != nil {
				cobra.CheckErr(err)
			}

			writer := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)

			fmt.Fprintf(writer, "Title\t%s\t\n", submission.Title)
			fmt.Fprintf(writer, "Author\t%s\t\n", submission.User())
			fmt.Fprintf(writer, "Date\t%s\t\n", submission.PostedAt.UTC().Format("2006-01-02 15:04:05 MST"))
			fmt.Fprintf(writer, "Post URL\t%s\t\n", submission.PostURL())
			fmt.Fprintf(writer, "Image URL\t%s\t\n", submission.ImageURL)
			fmt.Fprintf(writer, "Image Size\t%d x %d\t\n", submission.ImageWidthPx, submission.ImageHeightPx)
			fmt.Fprintf(writer, "Filename\t%s\t\n", submission.ImageFilename)

			writer.Flush()
		},
	}

	return cmd
}
