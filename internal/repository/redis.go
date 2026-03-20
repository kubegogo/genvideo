package repository

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/kubegogo/genvideo/internal/config"
)

func NewRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})
}

// CacheGet retrieves a value from cache
func CacheGet(ctx context.Context, rdb *redis.Client, key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

// CacheSet stores a value in cache with expiration
func CacheSet(ctx context.Context, rdb *redis.Client, key string, value interface{}, expiration int) error {
	return rdb.Set(ctx, key, value, 0).Err()
}
