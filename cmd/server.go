package cmd

import (
	"github.com/cartermckinnon/filter-feed/pkg/fetch"
	"github.com/cartermckinnon/filter-feed/pkg/service"
	"github.com/integrii/flaggy"
)

type command struct {
	flaggyCmd *flaggy.Subcommand
	address   *string
}

func NewServerCommand() Command {
	subcommand := flaggy.NewSubcommand("server")
	subcommand.Description = "Run the gRPC server"

	address := ":8080"
	subcommand.String(&address, "a", "address", "gRPC service address")

	return &command{
		flaggyCmd: subcommand,
		address:   &address,
	}
}

func (c *command) GetFlaggySubcommand() *flaggy.Subcommand {
	return c.flaggyCmd
}

func (c *command) Run() error {
	return service.RunHTTPService(*c.address, fetch.NewFeedFetcher())
}
