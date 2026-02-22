package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type CacheConfig struct {
	RedisURL string
}

func NewRedisClient(config CacheConfig) (*redis.Client, error) {
	opts, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
