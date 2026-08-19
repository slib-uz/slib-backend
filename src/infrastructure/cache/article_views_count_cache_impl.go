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

	//"strconv"
	"time"
)

const (
	viewsCountPrefix = "views_count:"
	viewsPrefix      = "view:"
)

type ArticleViewsCountCacheImpl struct {
	client *redis.Client
	env    *config.Config
}

// @inject
func NewArticleViewsCountCache(client *redis.Client, env *config.Config) cache.ArticleViewsCountCache {
	return &ArticleViewsCountCacheImpl{client: client, env: env}
}

func (this *ArticleViewsCountCacheImpl) Add(userKey string, articleID uint) (int64, error) {

	added, err := this.client.SetNX(
		context.Background(),
		this.viewsKey(userKey, articleID),
		articleID,
		time.Duration(this.env.ViewsCountLifetimeMinute)*time.Minute,
	).Result()

	if err != nil {
		return 0, err
	}

	if added {
		return this.inc(articleID)
	}

	return this.Get(articleID)
}

func (this *ArticleViewsCountCacheImpl) inc(articleID uint) (int64, error) {
	// Increment the views count for the article
	inc := this.client.Incr(context.Background(), this.viewsCountKey(articleID))
	if err := inc.Err(); err != nil {
		return 0, fmt.Errorf("failed to increment views count for article %d: %w", articleID, err)
	}
	return inc.Result()

}

func (this *ArticleViewsCountCacheImpl) Get(articleID uint) (int64, error) {
	val, err := this.client.Get(context.Background(), this.viewsCountKey(articleID)).Result()
	if errors.Is(err, redis.Nil) {
		return this.inc(articleID)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get views count for article %d: %w", articleID, err)
	}

	count, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse views count for article %d: %w", articleID, err)
	}

	return count, nil
}

func (this *ArticleViewsCountCacheImpl) viewsKey(userKey string, articleID uint) string {
	return fmt.Sprintf("%s%d:%x", viewsPrefix, articleID, sha1.Sum([]byte(userKey)))
}

func (this *ArticleViewsCountCacheImpl) viewsCountKey(articleID uint) string {
	return fmt.Sprintf("%s%d", viewsCountPrefix, articleID)
}

func (this *ArticleViewsCountCacheImpl) GetAll() (map[uint]int64, error) {
	const pattern = viewsCountPrefix + "*" // "views_count:*"
	const scanCount = int64(100)           // SCAN har iteratsiyada necha kalit ko‘rsin
	ctx := context.Background()

	out := make(map[uint]int64)

	iter := this.client.Scan(ctx, 0, pattern, scanCount).Iterator()
	for iter.Next(ctx) {
		key := iter.Val() // masalan: "views_count:123"
		val, err := this.client.GetDel(ctx, key).Int64()
		if err != nil {
			// kalit bu orada o‘chib ketgan bo‘lishi mumkin – davom etamiz
			if errors.Is(err, redis.Nil) {
				continue
			}
			return nil, err
		}
		if val == 0 {
			continue
		}

		idStr := strings.TrimPrefix(key, viewsCountPrefix)
		id64, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			// noto‘g‘ri kalit nomi – log qilishingiz mumkin
			continue
		}
		id := uint(id64)

		// Ba’zan bir article uchun bir nechta kalit bo‘lishi mumkin:
		out[id] += val
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
