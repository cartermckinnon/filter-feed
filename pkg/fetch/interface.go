package fetch

import "github.com/mmcdole/gofeed"

type FeedFetcher interface {
	FetchFeed(url string) (*gofeed.Feed, error)
}
