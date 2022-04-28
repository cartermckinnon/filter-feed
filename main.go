package main

import (
	stdlog "log"

	"github.com/cartermckinnon/filter-feed/cmd"
	"github.com/integrii/flaggy"
)

var log = stdlog.Default()

func main() {
	flaggy.SetName("filter-feed")
	flaggy.SetDescription("for fetchin' filtered feeds")
	flaggy.SetVersion("0.0.0-dev")

	var subcommands []cmd.Command
	subcommands = append(subcommands, cmd.NewServerCommand())
	subcommands = append(subcommands, cmd.NewFetchCommand())
	subcommands = append(subcommands, cmd.NewURLCommand())
	for _, subcommand := range subcommands {
		flaggy.AttachSubcommand(subcommand.GetFlaggySubcommand(), 1) // flaggy positions start at 1
	}

	flaggy.ShowHelpOnUnexpectedEnable()

	flaggy.Parse()

	var usedCommand cmd.Command
	for _, subcommand := range subcommands {
		if subcommand.GetFlaggySubcommand().Used {
			usedCommand = subcommand
			break

		}
	}
	if usedCommand == nil {
		flaggy.ShowHelpAndExit("error: no subcommand specified")
	} else if err := usedCommand.Run(); err != nil {
		log.Fatalf("error running command=%s: %v", usedCommand.GetFlaggySubcommand().Name, err)
	}
}
