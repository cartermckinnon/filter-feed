package fetch

import (
	"net/http"

	"github.com/mmcdole/gofeed"
)

type feedFetcherImpl struct {
	http   *http.Client
	parser *gofeed.Parser
}

func NewFeedFetcher() FeedFetcher {
	return &feedFetcherImpl{
		http:   &http.Client{},
		parser: gofeed.NewParser(),
	}
}

func (ff feedFetcherImpl) FetchFeed(feedUrl string) (*gofeed.Feed, error) {
	resp, err := ff.http.Get(feedUrl)
	if err != nil {
		return nil, err
	}
	feed, err := ff.parser.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return feed, nil
}
