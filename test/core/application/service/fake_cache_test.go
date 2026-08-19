package service_test

import (
	"context"
	"time"
)

// fakeCache — CacheProvider portining xotiradagi soxta amalga oshirilishi.
// getErr / setErr o'rnatilsa, Redis nosozligini taqlid qiladi.
type fakeCache struct {
	data   map[string]string
	ttls   map[string]time.Duration
	getErr error
	setErr error
}

func newFakeCache() *fakeCache {
	return &fakeCache{
		data: make(map[string]string),
		ttls: make(map[string]time.Duration),
	}
}

func (f *fakeCache) GetByKey(ctx context.Context, key string) (string, error) {
	if f.getErr != nil {
		return "", f.getErr
	}
	return f.data[key], nil
}

func (f *fakeCache) Set(ctx context.Context, key string, value string, expSeconds time.Duration) error {
	if f.setErr != nil {
		return f.setErr
	}
	f.data[key] = value
	f.ttls[key] = expSeconds
	return nil
}
