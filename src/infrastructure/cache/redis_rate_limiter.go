package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"slib.uz/src/core/domain/ports/ratelimit"
)

type RedisRateLimiter struct {
	client *redis.Client
}

// @inject
func NewRedisRateLimiter(client *redis.Client) ratelimit.RateLimiter {
	return &RedisRateLimiter{client: client}
}

// Hit kalitni atomik oshiradi. Natija 1 bo'lsa (yangi kalit), window muddatiga
// TTL o'rnatiladi. INCR atomik, shuning uchun ikki parallel urinish ham to'g'ri
// sanoq oladi.
func (this *RedisRateLimiter) Hit(ctx context.Context, key string, window time.Duration) (int64, error) {
	count, err := this.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		if err := this.client.Expire(ctx, key, window).Err(); err != nil {
			return count, err
		}
	}
	return count, nil
}

func (this *RedisRateLimiter) Reset(ctx context.Context, key string) error {
	return this.client.Del(ctx, key).Err()
}
