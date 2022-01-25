package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/virtualtam/walric/cmd/walric/formatter"
)

// NewCurrentCommand initializes a CLI command to display information for the
// currently selected entry.
func NewCurrentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Display information about the currently selected entry",
		Run: func(cmd *cobra.Command, args []string) {
			entry, err := historyService.Current()
			if err != nil {
				cobra.CheckErr(err)
			}

			fmt.Println("Current image, selected on", formatter.FormatDateAsUTC(entry.Date))
			fmt.Println()

			writer := formatter.FormatSubmissionAsTab(os.Stdout, entry.Submission)
			writer.Flush()
		},
	}

	return cmd
}
