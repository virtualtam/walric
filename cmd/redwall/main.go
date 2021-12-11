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

	infoCommand := command.NewInfoCommand()
	rootCommand.AddCommand(infoCommand)

	statsCommand := command.NewStatsCommand()
	rootCommand.AddCommand(statsCommand)

	cobra.CheckErr(rootCommand.Execute())
}
