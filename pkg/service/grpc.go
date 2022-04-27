package service

import (
	"context"
	"log"

	api "github.com/cmckn/filter-feed/pkg/api/v1"
	"github.com/cmckn/filter-feed/pkg/fetch"
	"github.com/cmckn/filter-feed/pkg/filter"
	"github.com/cmckn/filter-feed/pkg/util"
	"google.golang.org/grpc"
)

func RegisterFilterFeedServer(grpc *grpc.Server, fetcher fetch.FeedFetcher) {
	api.RegisterFilterFeedServer(grpc, &server{
		fetcher: fetcher,
	})
}

type server struct {
	api.UnimplementedFilterFeedServer
	fetcher fetch.FeedFetcher
}

func (s *server) FetchFeed(ctx context.Context, in *api.FetchFeedRequest) (*api.FetchFeedResponse, error) {
	feed, err := s.fetcher.FetchFeed(in.FeedUrl)
	if err != nil {
		log.Printf("error fetching feedUrl=%s", in.FeedUrl)
		return nil, err
	}

	log.Printf("fetched feedType=%s feedUrl=%s ", in.FeedUrl, feed.FeedType)

	items, err := filter.FilterItems(*feed, in.Filters)
	if err != nil {
		log.Printf("error filtering feedUrl=%s %v", in.FeedUrl, err)
		return nil, err
	}
	feed.Items = items

	content, err := util.ConvertFeedToString(feed)
	if err != nil {
		log.Printf("error converting feedUrl=%s %v", in.FeedUrl, err)
		return nil, err
	}

	return &api.FetchFeedResponse{
		Content: content,
	}, nil
}
