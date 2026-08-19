package cache

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"slib.uz/src/core/domain/ports/cache"
	"slib.uz/src/infrastructure/config"

	"time"
)

const (
	newsViewsCountPrefix = "news_views_count:"
	newsViewsPrefix      = "news_view:"
)

type NewsViewsCountCacheImpl struct {
	client *redis.Client
	env    *config.Config
}

// @inject
func NewNewsViewsCountCache(client *redis.Client, env *config.Config) cache.NewsViewsCountCache {
	return &NewsViewsCountCacheImpl{client: client, env: env}
}

func (this *NewsViewsCountCacheImpl) Add(userKey string, newsID uint) (int64, error) {

	added, err := this.client.SetNX(
		context.Background(),
		this.viewsKey(userKey, newsID),
		newsID,
		time.Duration(this.env.ViewsCountLifetimeMinute)*time.Minute,
	).Result()

	if err != nil {
		return 0, err
	}

	if added {
		return this.inc(newsID)
	}

	return this.Get(newsID)
}

func (this *NewsViewsCountCacheImpl) inc(newsID uint) (int64, error) {
	inc := this.client.Incr(context.Background(), this.viewsCountKey(newsID))
	if err := inc.Err(); err != nil {
		return 0, fmt.Errorf("failed to increment views count for news %d: %w", newsID, err)
	}
	return inc.Result()

}

func (this *NewsViewsCountCacheImpl) Get(newsID uint) (int64, error) {
	val, err := this.client.Get(context.Background(), this.viewsCountKey(newsID)).Result()
	if errors.Is(err, redis.Nil) {
		return this.inc(newsID)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get views count for news %d: %w", newsID, err)
	}

	count, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse views count for news %d: %w", newsID, err)
	}

	return count, nil
}

func (this *NewsViewsCountCacheImpl) viewsKey(userKey string, newsID uint) string {
	return fmt.Sprintf("%s%d:%x", newsViewsPrefix, newsID, sha1.Sum([]byte(userKey)))
}

func (this *NewsViewsCountCacheImpl) viewsCountKey(newsID uint) string {
	return fmt.Sprintf("%s%d", newsViewsCountPrefix, newsID)
}

func (this *NewsViewsCountCacheImpl) GetAll() (map[uint]int64, error) {
	const pattern = newsViewsCountPrefix + "*"
	const scanCount = int64(100)
	ctx := context.Background()

	out := make(map[uint]int64)

	iter := this.client.Scan(ctx, 0, pattern, scanCount).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := this.client.GetDel(ctx, key).Int64()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			return nil, err
		}
		if val == 0 {
			continue
		}

		idStr := strings.TrimPrefix(key, newsViewsCountPrefix)
		id64, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			continue
		}
		id := uint(id64)

		out[id] += val
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
