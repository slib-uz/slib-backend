package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"slib.uz/src/core/domain/ports/cache"
	"slib.uz/src/infrastructure/config"
)

type RedisCache struct {
	client *redis.Client
	env    *config.Config
}

// @inject
func NewRedisCache(client *redis.Client, env *config.Config) cache.CacheProvider {
	return &RedisCache{client: client, env: env}
}

func (this *RedisCache) GetByKey(ctx context.Context, key string) (string, error) {
	val, err := this.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("[RedisCache] GetByKey error while getting value from Redis: %w", err)
	}
	return val, nil
}

func (this *RedisCache) Set(ctx context.Context, key string, value string, expSeconds time.Duration) error {
	return this.client.Set(ctx, key, value, expSeconds).Err()
}
