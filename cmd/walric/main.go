package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/virtualtam/walric/cmd/walric/command"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	rootCommand := command.NewRootCommand()

	commands := []*cobra.Command{
		command.NewCurrentCommand(),
		command.NewGatherCommand(),
		command.NewHistoryCommand(),
		command.NewInfoCommand(),
		command.NewListCandidatesCommand(),
		command.NewMigrateCommand(),
		command.NewRandomCommand(),
		command.NewSearchCommand(),
		command.NewStatsCommand(),
	}

	for _, cmd := range commands {
		rootCommand.AddCommand(cmd)
	}

	cobra.CheckErr(rootCommand.Execute())
}
