package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
	"github.com/integrii/flaggy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type fetchCommand struct {
	flaggyCmd *flaggy.Subcommand
	endpoint  *string
	feedUrl   *string
	filter    *string
}

func NewFetchCommand() Command {
	subcommand := flaggy.NewSubcommand("fetch")
	subcommand.Description = "Fetch a feed from a filter-feed server"

	endpoint := "localhost:8080"
	subcommand.String(&endpoint, "e", "endpoint", "gRPC service endpoint")

	var feedUrl string
	subcommand.AddPositionalValue(&feedUrl, "FEED_URL", 1, true, "Feed URL")

	var filter string
	subcommand.String(&filter, "f", "filter", "Filter to apply to the feed. Can be a JSON string or a file path.")

	return &fetchCommand{
		flaggyCmd: subcommand,
		endpoint:  &endpoint,
		feedUrl:   &feedUrl,
		filter:    &filter,
	}
}

func (c *fetchCommand) GetFlaggySubcommand() *flaggy.Subcommand {
	return c.flaggyCmd
}

func (c *fetchCommand) Run() error {
	conn, err := grpc.Dial(*c.endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server=%s, %v", *c.endpoint, err)
	}
	defer conn.Close()

	client := api.NewFilterFeedClient(conn)

	filters, err := unmarshalFilters(*c.filter)
	if err != nil {
		log.Fatalf("failed to unmarshal filter=%s, %v", *c.filter, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	r, err := client.FetchFeed(ctx, &api.FetchFeedRequest{
		FeedUrl: *c.feedUrl,
		Filters: filters,
	})
	if err != nil {
		log.Fatalf("failed to fetch feedUrl=%s, %v", *c.feedUrl, err)
	}
	fmt.Printf("%s", r.Content)
	return nil
}

func unmarshalFilters(filter string) ([]*api.FilterSpec, error) {
	if strings.HasPrefix(filter, "file://") {
		bytes, err := ioutil.ReadFile(filter[7:])
		if err != nil {
			return nil, err
		}
		filter = string(bytes)
	}
	if strings.HasPrefix(filter, "{") {
		var f api.FilterSpec
		if err := protojson.Unmarshal([]byte(filter), &f); err != nil {
			return nil, err
		}
		return []*api.FilterSpec{&f}, nil
	}
	if strings.HasPrefix(filter, "[") {
		var fs api.FilterSpecs
		if err := protojson.Unmarshal([]byte(filter), &fs); err != nil {
			return nil, err
		}
		return fs.Specs, nil
	}
	return nil, nil
}
