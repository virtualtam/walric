package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/virtualtam/redwall2/cmd/redwall/command"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	rootCommand := command.NewRootCommand()

	commands := []*cobra.Command{
		command.NewCurrentCommand(),
		command.NewHistoryCommand(),
		command.NewInfoCommand(),
		command.NewListCandidatesCommand(),
		command.NewRandomCommand(),
		command.NewSearchCommand(),
		command.NewStatsCommand(),
	}

	for _, cmd := range commands {
		rootCommand.AddCommand(cmd)
	}

	cobra.CheckErr(rootCommand.Execute())
}
