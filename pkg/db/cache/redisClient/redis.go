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

	// 1. Пытаемся получить настройки из переменных окружения (для Docker)
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	// 2. Если переменных нет, используем дефолтные значения или те, что в конфиге (для локального запуска)
	if cfg.Redis.Host == "" {
		host = "localhost"
	}
	if cfg.Redis.Port == "" {
		port = "6379"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	// Парсим остальные настройки
	db := cfg.Redis.DefaultDB

	// Безопасное чтение таймаура
	timeoutInt, err := strconv.Atoi(cfg.Redis.Timeout)

	if timeoutInt == 0 {
		timeoutInt = 5 // Дефолтный таймаут 5 секунд, если в конфиге пусто
	}

	fmt.Printf("Connecting to Redis at: %s\n", addr) // Лог для отладки

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
