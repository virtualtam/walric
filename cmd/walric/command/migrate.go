package command

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"

	"github.com/virtualtam/walric/internal/storage/sqlite3/migrations"
)

var _ migrate.Logger = &migrateLogger{}

type migrateLogger struct {
	verbose bool
}

func (l migrateLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l migrateLogger) Verbose() bool {
	return l.verbose
}

// NewMigrateCommand initializes a CLI command to create database tables and run
// SQL migrations.
func NewMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Initialize database and run migrations",
		Run: func(cmd *cobra.Command, args []string) {
			if err := os.MkdirAll(walricConfig.DataDir(), os.ModePerm); err != nil {
				cobra.CheckErr(err)
			}

			migrationsSource, err := iofs.New(migrations.MigrationsFS, ".")
			if err != nil {
				cobra.CheckErr(err)
			}

			migrater, err := migrate.NewWithSourceInstance(
				"iofs",
				migrationsSource,
				fmt.Sprintf("sqlite3://%s", walricConfig.DatabasePath()),
			)
			if err != nil {
				cobra.CheckErr(err)
			}

			migrater.Log = migrateLogger{
				verbose: debugMode,
			}

			err = migrater.Up()
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("Database schema already up to date")
			} else if err != nil {
				cobra.CheckErr(err)
			}
		},
	}

	return cmd
}
