package cmd

import (
	"fmt"
	"log"

	"github.com/cartermckinnon/filter-feed/pkg/config"
	"github.com/cartermckinnon/filter-feed/pkg/fetch"
	"github.com/cartermckinnon/filter-feed/pkg/filter"
	"github.com/cartermckinnon/filter-feed/pkg/util"
	"github.com/integrii/flaggy"
)

type fetchCommand struct {
	flaggyCmd *flaggy.Subcommand
	feedUrl   *string
	filter    *string
	fetcher   fetch.FeedFetcher
}

func NewFetchCommand() Command {
	subcommand := flaggy.NewSubcommand("fetch")
	subcommand.Description = "Fetch a feed from a filter-feed server"

	var feedUrl string
	subcommand.AddPositionalValue(&feedUrl, "FEED_URL", 1, true, "Feed URL")

	var filter string
	subcommand.String(&filter, "f", "filter", "Filter to apply to the feed. Can be a JSON string or a file path.")

	return &fetchCommand{
		flaggyCmd: subcommand,
		feedUrl:   &feedUrl,
		filter:    &filter,
		fetcher:   fetch.NewFeedFetcher(&config.RedisConfig{Enabled: util.NewFalse()}),
	}
}

func (c *fetchCommand) GetFlaggySubcommand() *flaggy.Subcommand {
	return c.flaggyCmd
}

func (c *fetchCommand) Run() error {
	filters, err := unmarshalFilters(*c.filter)
	if err != nil {
		return fmt.Errorf("failed to unmarshal filter=%s, %v", *c.filter, err)
	}

	feed, err := c.fetcher.FetchFeed(*c.feedUrl)
	if err != nil {
		return fmt.Errorf("failed to fetch feedUrl=%s, %v", *c.feedUrl, err)
	}

	log.Printf("fetched feedType=%s feedURL=%s ", feed.FeedType, *c.feedUrl)

	items, err := filter.FilterItems(*feed, filters)
	if err != nil {
		return fmt.Errorf("error filtering feedURL=%s %v", *c.feedUrl, err)
	}
	feed.Items = items

	content, err := util.ConvertFeedToString(feed)
	if err != nil {
		return fmt.Errorf("error converting feedURL=%s %v", *c.feedUrl, err)
	}

	log.Println(content)
	return nil
}
