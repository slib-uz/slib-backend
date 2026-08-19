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
	journalViewsCountPrefix = "journal_views_count:"
	journalViewsPrefix      = "journal_view:"
)

type JournalViewsCountCacheImpl struct {
	client *redis.Client
	env    *config.Config
}

// @inject
func NewJournalViewsCountCache(client *redis.Client, env *config.Config) cache.JournalViewsCountCache {
	return &JournalViewsCountCacheImpl{client: client, env: env}
}

func (this *JournalViewsCountCacheImpl) Add(userKey string, journalID uint) (int64, error) {

	added, err := this.client.SetNX(
		context.Background(),
		this.viewsKey(userKey, journalID),
		journalID,
		time.Duration(this.env.ViewsCountLifetimeMinute)*time.Minute,
	).Result()

	if err != nil {
		return 0, err
	}

	if added {
		return this.inc(journalID)
	}

	return this.Get(journalID)
}

func (this *JournalViewsCountCacheImpl) inc(journalID uint) (int64, error) {
	inc := this.client.Incr(context.Background(), this.viewsCountKey(journalID))
	if err := inc.Err(); err != nil {
		return 0, fmt.Errorf("failed to increment views count for journal %d: %w", journalID, err)
	}
	return inc.Result()

}

func (this *JournalViewsCountCacheImpl) Get(journalID uint) (int64, error) {
	val, err := this.client.Get(context.Background(), this.viewsCountKey(journalID)).Result()
	if errors.Is(err, redis.Nil) {
		return this.inc(journalID)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get views count for journal %d: %w", journalID, err)
	}

	count, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse views count for journal %d: %w", journalID, err)
	}

	return count, nil
}

func (this *JournalViewsCountCacheImpl) viewsKey(userKey string, journalID uint) string {
	return fmt.Sprintf("%s%d:%x", journalViewsPrefix, journalID, sha1.Sum([]byte(userKey)))
}

func (this *JournalViewsCountCacheImpl) viewsCountKey(journalID uint) string {
	return fmt.Sprintf("%s%d", journalViewsCountPrefix, journalID)
}

func (this *JournalViewsCountCacheImpl) GetAll() (map[uint]int64, error) {
	const pattern = journalViewsCountPrefix + "*"
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

		idStr := strings.TrimPrefix(key, journalViewsCountPrefix)
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
