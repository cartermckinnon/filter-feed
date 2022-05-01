package cmd

import (
	"github.com/cartermckinnon/filter-feed/pkg/config"
	"github.com/cartermckinnon/filter-feed/pkg/fetch"
	"github.com/cartermckinnon/filter-feed/pkg/service"
	"github.com/integrii/flaggy"
)

type command struct {
	flaggyCmd *flaggy.Subcommand
	address   *string
	redis     *config.RedisConfig
}

func NewServerCommand() Command {
	subcommand := flaggy.NewSubcommand("server")
	subcommand.Description = "Run the HTTP server"

	address := ":8080"
	subcommand.String(&address, "a", "address", "HTTP service address")

	redisAddress := ":6379"
	subcommand.String(&redisAddress, "r", "redis-address", "Redis address")

	redisDb := 0
	subcommand.Int(&redisDb, "d", "redis-db", "Redis database")

	redisUsername := ""
	subcommand.String(&redisUsername, "u", "redis-username", "Redis username")

	redisPassword := ""
	subcommand.String(&redisPassword, "p", "redis-password", "Redis password")

	redisTTL := "15m"
	subcommand.String(&redisTTL, "t", "redis-ttl", "Redis TTL")

	redisEnabled := true
	subcommand.Bool(&redisEnabled, "e", "redis-enabled", "Redis enabled")

	return &command{
		flaggyCmd: subcommand,
		address:   &address,
		redis: &config.RedisConfig{
			Address:  &redisAddress,
			DB:       &redisDb,
			Username: &redisUsername,
			Password: &redisPassword,
			TTL:      &redisTTL,
			Enabled:  &redisEnabled,
		},
	}
}

func (c *command) GetFlaggySubcommand() *flaggy.Subcommand {
	return c.flaggyCmd
}

func (c *command) Run() error {
	return service.RunHTTPService(*c.address, fetch.NewFeedFetcher(c.redis))
}
