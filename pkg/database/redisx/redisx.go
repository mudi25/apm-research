package redisx

import (
	"context"

	"github.com/go-redis/redis/v8"
	apmgoredis "go.elastic.co/apm/module/apmgoredisv8/v2"
)

// Config holds the Redis connection configuration.
type Config struct {
	Url    string // Redis connection URL, e.g. "redis://user:password@localhost:6379/0"
	UseApm bool   // Enable Elastic APM hook if true
}

// NewClient creates a new Redis client from the given configuration.
// - Parses the connection URL
// - Initializes the client
// - Optionally adds Elastic APM hook
// - Pings Redis to verify the connection
func NewClient(ctx context.Context, cfg Config) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.Url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	// Attach Elastic APM hook if enabled
	if cfg.UseApm {
		client.AddHook(apmgoredis.NewHook())
	}

	// Test the connection with a ping
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

// Disconnect cleanly closes the Redis client connection.
// Safe to call multiple times; does nothing if client is nil.
func Disconnect(client *redis.Client) error {
	if client == nil {
		return nil
	}
	return client.Close()
}
