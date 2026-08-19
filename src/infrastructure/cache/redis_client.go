package cache

import (
	"github.com/redis/go-redis/v9"
	"slib.uz/src/infrastructure/config"
)

// @inject
func NewRedisClient(env *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: env.RedisAddress,
		DB:   1,
	})
}
