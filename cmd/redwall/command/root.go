package command

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/virtualtam/redwall2/cmd/redwall/config"
	"github.com/virtualtam/redwall2/history"
	"github.com/virtualtam/redwall2/submission"
	"github.com/virtualtam/redwall2/subreddit"
)

const (
	defaultDebugMode bool = false
)

var (
	configPath string
	debugMode  bool

	historyService    *history.Service
	submissionService *submission.Service
	subredditService  *subreddit.Service
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

			redwallConfig, err := config.LoadTOML(configPath)
			if err != nil {
				return err
			}

			db, err := sqlx.Open("sqlite3", redwallConfig.DatabasePath())
			if err != nil {
				return err
			}

			subredditRepository := subreddit.NewRepositorySQLite(db)
			subredditService = subreddit.NewService(subredditRepository)

			submissionRepository := submission.NewRepositorySQLite(db)
			submissionService = submission.NewService(submissionRepository, subredditService)

			historyRepository := history.NewRepositorySQLite(db)
			historyService = history.NewService(historyRepository, submissionService)

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(
		&configPath,
		"config",
		"",
		"Configuration file",
	)
	cmd.PersistentFlags().BoolVar(
		&debugMode,
		"debug",
		defaultDebugMode,
		"Enable debugging",
	)

	return cmd
}
