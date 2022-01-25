package command

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/virtualtam/walric/cmd/walric/config"
	"github.com/virtualtam/walric/history"
	"github.com/virtualtam/walric/submission"
	"github.com/virtualtam/walric/subreddit"
)

const (
	defaultDebugMode bool = false
)

var (
	configPath string
	debugMode  bool

	walricConfig *config.Config

	historyService    *history.Service
	submissionService *submission.Service
	subredditService  *subreddit.Service
)

// NewRootCommand initializes the main CLI entrypoint and common command flags.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "walric",
		Short: "Walric",
		Long: `Walric, the front wallpaper to your monitor(s)

Walric helps you manage a collection of curated wallpapers, courtesy of the Reddit community.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if debugMode {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}

			cfg, err := config.LoadTOML(configPath)
			if err != nil {
				return err
			}
			walricConfig = cfg

			db, err := sqlx.Open("sqlite3", walricConfig.DatabasePath())
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
