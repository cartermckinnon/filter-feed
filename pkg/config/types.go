package config

type RedisConfig struct {
	Address  *string
	DB       *int
	Username *string
	Password *string
	TTL      *string
	Enabled  *bool
}
