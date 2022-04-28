package cmd

import (
	"encoding/base64"
	"fmt"
	"log"

	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
	"github.com/cartermckinnon/filter-feed/pkg/service"
	"github.com/integrii/flaggy"
	"google.golang.org/protobuf/proto"
)

type urlCommand struct {
	flaggyCmd     *flaggy.Subcommand
	feedUrl       *string
	filterVersion *uint
	filter        *string
}

func NewURLCommand() Command {
	subcommand := flaggy.NewSubcommand("url")
	subcommand.Description = "Generate a URL for a filtered feed"

	var feedUrl string
	subcommand.AddPositionalValue(&feedUrl, "FEED_URL", 1, true, "Feed URL.")

	var filter string
	subcommand.AddPositionalValue(&filter, "FILTER", 2, true, "Filter to apply to the feed. Can be a JSON string or a file path.")

	var filterVersion uint = 1
	subcommand.UInt(&filterVersion, "v", "filter-version", "API version of the provided filter.")

	return &urlCommand{
		flaggyCmd:     subcommand,
		feedUrl:       &feedUrl,
		filter:        &filter,
		filterVersion: &filterVersion,
	}
}

func (c *urlCommand) GetFlaggySubcommand() *flaggy.Subcommand {
	return c.flaggyCmd
}

func (c *urlCommand) Run() error {
	filters, err := unmarshalFilters(*c.filter)
	if err != nil {
		log.Fatalf("failed to unmarshal filter=%s, %v", *c.filter, err)
	}
	request := api.FetchFeedRequest{
		FeedURL: *c.feedUrl,
		Filters: filters,
	}
	bytes, err := proto.Marshal(&request)
	if err != nil {
		log.Fatalf("failed to marshal request=%v, %v", request, err)
	}
	base64Request := base64.StdEncoding.EncodeToString(bytes)
	fmt.Printf("%s%s", service.FilterFeedV1Path, base64Request)
	return nil
}
