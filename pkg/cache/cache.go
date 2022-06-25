package cache

import (
	cache "broker/pkg/cache/redis"
	"context"
	"time"
)

type Cache[T any] interface {
	Set(ctx context.Context, key string, value *T) error
	Get(ctx context.Context, key string) (value *T, err error)
	Delete(ctx context.Context, key string) error
	Exist() bool
	ExparationTime() time.Duration
}

type CacheWrapper[T any] struct {
	cache  *Cache[T]
	prefix string
}

func newCache[T any](r *cache.RedisWrapper, prefix string) *CacheWrapper[T] {
	return &CacheWrapper[T]{
		cache:  r,
		prefix: prefix,
	}
}
