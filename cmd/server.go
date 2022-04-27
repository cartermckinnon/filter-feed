package cmd

import (
	"log"
	"net"

	"github.com/cmckn/filter-feed/pkg/fetch"
	"github.com/cmckn/filter-feed/pkg/service"
	"github.com/integrii/flaggy"
	"google.golang.org/grpc"
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

	addr := *c.address
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("failed to listen on address=%s", addr)
		return err
	}
	log.Printf("server listening on address=%s", lis.Addr())

	grpc := grpc.NewServer()
	service.RegisterFilterFeedServer(grpc, fetch.NewFeedFetcher())

	if err := grpc.Serve(lis); err != nil {
		return err
	}

	return nil
}
