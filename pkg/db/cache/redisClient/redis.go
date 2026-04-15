// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package redisClient

import (
	"context"
	"fmt"
	"os"
	"shb/internal/configs"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisClient() (*RedisCache, error) {
	cfg, err := configs.InitConfigs()
	if err != nil {
		return nil, err
	}

	// 1. Try to read connection settings from environment variables (Docker).
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	// 2. Fall back to config values if environment variables are not set (local run).
	if cfg.Redis.Host == "" {
		host = "localhost"
	}
	if cfg.Redis.Port == "" {
		port = "6379"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	// Parse remaining settings from config.
	db := cfg.Redis.DefaultDB

	// Safely read the timeout value.
	timeoutInt, err := strconv.Atoi(cfg.Redis.Timeout)

	if timeoutInt == 0 {
		timeoutInt = 5 // Default to 5 seconds if the config value is empty.
	}

	fmt.Printf("Connecting to Redis at: %s\n", addr)

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInt)*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis at %s: %w", addr, err)
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
