package command

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/virtualtam/walric/cmd/walric/formatter"
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

			writer := formatter.FormatSubmissionAsTab(os.Stdout, submission)
			writer.Flush()
		},
	}

	return cmd
}
