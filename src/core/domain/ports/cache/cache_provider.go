package cache

import (
	"context"
	"time"
)

type CacheProvider interface {
	GetByKey(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expSeconds time.Duration) error
}
