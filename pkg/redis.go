package pkg

import (
	"Auth/config"
	"Auth/helper"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

type RedisRepository struct {
	client *redis.Client
}

func InitializeRedis(cfg *config.Config) *RedisRepository {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		helper.ThrowError(err)
	}

	return &RedisRepository{client: RedisClient}
}

func Set(key string, value any, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

func Delete(key string) error {
	return RedisClient.Del(ctx, key).Err()
}
