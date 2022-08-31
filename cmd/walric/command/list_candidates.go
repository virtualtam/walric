package command

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/virtualtam/walric/pkg/monitor"
)

const (
	defaultXRandRScreenNo int = 0
)

var (
	xRandRScreenNo int
)

// NewListCandidates initializes a CLI command to list Submissions suitable for
// the current monitor configuration.
func NewListCandidatesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-candidates",
		Short: "List submissions suitable for the current monitor setup",
		Run: func(cmd *cobra.Command, args []string) {
			monitors, err := monitor.ConnectedMonitors(xRandRScreenNo)
			if err != nil {
				cobra.CheckErr(err)
			}

			wallpaperResolution := monitor.WallpaperResolution(monitors)

			submissions, err := submissionService.ByMinResolution(wallpaperResolution)
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
		},
	}

	cmd.Flags().IntVar(
		&xRandRScreenNo,
		"xrandr-screen",
		defaultXRandRScreenNo,
		"XRandR screen",
	)

	return cmd
}
