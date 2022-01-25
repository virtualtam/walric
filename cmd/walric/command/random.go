package command

import (
	"os"

	"github.com/spf13/cobra"
	walric "github.com/virtualtam/walric"
	"github.com/virtualtam/walric/cmd/walric/formatter"
	"github.com/virtualtam/walric/monitor"
)

// NewRandomCommand initializes a CLI command to select a random submission
// suitable for the current monitor setup.
func NewRandomCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "random",
		Short: "Select a random submission suitable for the current monitor setup",
		Run: func(cmd *cobra.Command, args []string) {
			monitors, err := monitor.ConnectedMonitors(xRandRScreenNo)
			if err != nil {
				cobra.CheckErr(err)
			}

			wallpaperResolution := walric.WallpaperResolution(monitors)

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
