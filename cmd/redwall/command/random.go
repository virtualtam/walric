package command

import (
	"os"

	"github.com/spf13/cobra"
	redwall "github.com/virtualtam/redwall2"
	"github.com/virtualtam/redwall2/cmd/redwall/formatter"
)

// NewRandomCommand initializes a CLI command to select a random submission
// suitable for the current monitor setup.
func NewRandomCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "random",
		Short: "Select a random submission suitable for the current monitor setup",
		Run: func(cmd *cobra.Command, args []string) {
			monitors, err := redwall.ConnectedMonitors(xRandRScreenNo)
			if err != nil {
				cobra.CheckErr(err)
			}

			wallpaperResolution := redwall.WallpaperResolution(monitors)

			submission, err := submissionService.Random(wallpaperResolution)
			if err != nil {
				cobra.CheckErr(err)
			}

			if err := historyService.Save(submission); err != nil {
				cobra.CheckErr(err)
			}

			writer := formatter.FormatSubmissionAsTab(os.Stdout, submission)
			writer.Flush()
		},
	}

	return cmd
}
