package command

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// NewStatsCommand initializes a CLI command to display statistics.
func NewStatsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Display statistics about gathered submissions",
		Run: func(cmd *cobra.Command, args []string) {
			stats, err := subredditService.Stats()
			if err != nil {
				cobra.CheckErr(err)
			}

			writer := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)

			fmt.Fprintln(writer, "Count\tSubreddit\t")
			fmt.Fprintln(writer, "-----\t---------\t")
			fmt.Fprintln(writer, "\t\t")

			var total int

			for _, subredditStats := range stats {
				total += subredditStats.Submissions
				fmt.Fprintf(writer, "%d\t%s\t\n", subredditStats.Submissions, subredditStats.Name)
			}

			fmt.Fprintln(writer, "\t\t")
			fmt.Fprintf(writer, "%d\t%s\t\n", total, "TOTAL")

			writer.Flush()
		},
	}

	return cmd
}
