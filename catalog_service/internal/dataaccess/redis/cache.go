package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("not found")
)

type RedisCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) RedisCache {
	return &redisCache{
		client: client,
	}
}

func (rc *redisCache) Get(ctx context.Context, key string) (string, error) {
	data, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}

	return data, nil
}
func (rc *redisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return rc.client.Set(ctx, key, value, ttl).Err()
}
func (rc *redisCache) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}
