package service

import (
	"encoding/base64"
	"log"
	"net/http"

	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
	"github.com/cartermckinnon/filter-feed/pkg/fetch"
	"github.com/cartermckinnon/filter-feed/pkg/filter"
	"github.com/cartermckinnon/filter-feed/pkg/util"
	"google.golang.org/protobuf/proto"
)

const FilterFeedV1Path = "/v1/f/"

var noFilters = []byte("no filters provided")

type filterFeedV1Handler struct {
	fetcher fetch.FeedFetcher
}

func (h *filterFeedV1Handler) filterFeedV1(w http.ResponseWriter, req *http.Request) {
	base64Payload := req.URL.Path[len(FilterFeedV1Path):]
	payload, err := base64.StdEncoding.DecodeString(base64Payload)
	if err != nil {
		log.Printf("error decoding base64 payload: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var ffr = api.FetchFeedRequest{}
	err = proto.Unmarshal(payload, &ffr)
	if err != nil {
		log.Printf("error unmarshalling payload: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(ffr.Filters) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		w.Write(noFilters)
		return
	}

	feed, err := h.fetcher.FetchFeed(ffr.FeedURL)
	if err != nil {
		log.Printf("error fetching feedURL=%s, %v", ffr.FeedURL, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("fetched feedType=%s feedURL=%s ", feed.FeedType, ffr.FeedURL)

	items, err := filter.FilterItems(*feed, ffr.Filters)
	if err != nil {
		log.Printf("error filtering feedURL=%s %v", ffr.FeedURL, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	feed.Items = items

	content, err := util.ConvertFeedToString(feed)
	if err != nil {
		log.Printf("error converting feedURL=%s %v", ffr.FeedURL, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}

func RunHTTPService(addr string, fetcher fetch.FeedFetcher) error {
	http.HandleFunc(FilterFeedV1Path, (&filterFeedV1Handler{
		fetcher: fetcher,
	}).filterFeedV1)
	return http.ListenAndServe(addr, nil)
}
