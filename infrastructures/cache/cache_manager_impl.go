package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCacheManager struct {
	client *redis.Client
}

func NewRedisCacheManager(client *redis.Client) CacheManager {
	return &redisCacheManager{client: client}
}

func (r *redisCacheManager) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("cache get failed: %w", err)
	}
	return val, nil
}

func (r *redisCacheManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache set marshal failed: %w", err)
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *redisCacheManager) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.Set(ctx, key, value, ttl)
}

func (r *redisCacheManager) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisCacheManager) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("cache exists check failed: %w", err)
	}
	return result > 0, nil
}

func (r *redisCacheManager) Publish(ctx context.Context, channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("cache publish marshal failed: %w", err)
	}
	return r.client.Publish(ctx, channel, data).Err()
}

func (r *redisCacheManager) Subscribe(ctx context.Context, channel string) (<-chan string, error) {
	sub := r.client.Subscribe(ctx, channel)
	ch := make(chan string, 100)

	go func() {
		defer close(ch)
		for msg := range sub.Channel() {
			select {
			case ch <- msg.Payload:
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}
