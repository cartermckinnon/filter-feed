package fetch

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cartermckinnon/filter-feed/pkg/config"
	"github.com/go-redis/redis/v8"
	"github.com/mmcdole/gofeed"
)

type feedFetcherImpl struct {
	http   *http.Client
	parser *gofeed.Parser
	redis  *redis.Client
	ttl    time.Duration
}

func NewFeedFetcher(redisConfig *config.RedisConfig) FeedFetcher {
	ff := &feedFetcherImpl{
		http:   &http.Client{},
		parser: gofeed.NewParser(),
	}
	if *redisConfig.Enabled {
		log.Printf("using redis at address=%s for feed caching", *redisConfig.Address)
		ff.redis = redis.NewClient(&redis.Options{
			Addr:     *redisConfig.Address,
			DB:       *redisConfig.DB,
			Username: *redisConfig.Username,
			Password: *redisConfig.Password,
		})
		ttl, err := time.ParseDuration(*redisConfig.TTL)
		if err != nil {
			log.Fatalf("failed to parse redisTTL=%v: %v", redisConfig.TTL, err)
		}
		ff.ttl = ttl
	}
	return ff
}

func (ff feedFetcherImpl) FetchFeed(feedUrl string) (*gofeed.Feed, error) {
	if ff.redis == nil {
		return ff.parser.ParseURL(feedUrl)
	}
	cached, err := ff.redis.Get(context.Background(), feedUrl).Result()
	if err != nil {
		resp, err := ff.http.Get(feedUrl)
		if err != nil {
			return nil, err
		}
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = ff.redis.Set(context.Background(), feedUrl, bytes, ff.ttl).Err()
		if err != nil {
			return nil, err
		}
		log.Printf("downloaded feedUrl=%s", feedUrl)
		return ff.parser.Parse(strings.NewReader(string(bytes)))
	}
	log.Printf("using cached feedUrl=%s", feedUrl)
	return ff.parser.Parse(strings.NewReader(cached))
}
