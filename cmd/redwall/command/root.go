package command

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	redwall "github.com/virtualtam/redwall2"
)

const (
	defaultDebugMode bool = false
)

var (
	debugMode bool

	submissionService *redwall.SubmissionService
	subredditService  *redwall.SubredditService
)

// NewRootCommand initializes the main CLI entrypoint and common command flags.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redwall",
		Short: "Redwall",
		Long: `Redwall, the front wallpaper to your monitor(s)

Redwall helps you manage a collection of curated wallpapers, courtesy of the Reddit community.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if debugMode {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}

			userHome, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			configPath := filepath.Join(userHome, ".config", "redwall.toml")

			configBytes, err := os.ReadFile(configPath)
			if err != nil {
				return err
			}

			config := &redwall.Config{}
			_, err = toml.Decode(string(configBytes), config)
			if err != nil {
				return err
			}

			db, err := sqlx.Open("sqlite3", config.DatabasePath())
			if err != nil {
				return err
			}

			subredditRepository := redwall.NewSubredditRepositorySQLite(db)
			subredditService = redwall.NewSubredditService(subredditRepository)

			submissionRepository := redwall.NewSubmissionRepositorySQLite(db, subredditService)
			submissionService = redwall.NewSubmissionService(submissionRepository)

			return nil
		},
	}

	cmd.PersistentFlags().BoolVar(
		&debugMode,
		"debug",
		defaultDebugMode,
		"Enable debugging",
	)

	return cmd
}
