package redisClient

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"shb/pkg/configs"
	"strconv"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisClient() (*RedisCache, error) {
	db, _ := strconv.Atoi(configs.RedisDefaultDB)
	timeout, _ := strconv.Atoi(configs.RedisTimeout)

	client := redis.NewClient(&redis.Options{
		Addr: configs.RedisHost,
		DB:   db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &RedisCache{client: client}, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Increment(ctx context.Context, key string) error {
	return r.client.Incr(ctx, key).Err()
}
