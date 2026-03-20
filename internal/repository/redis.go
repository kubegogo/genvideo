package repository

import (
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
