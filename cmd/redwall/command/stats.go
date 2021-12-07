package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewStatsCommand initializes a CLI command to display statistics.
func NewStatsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Display statistics about gathered submissions",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello")
		},
	}

	return cmd
}
