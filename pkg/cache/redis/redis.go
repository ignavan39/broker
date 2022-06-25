package cache

import (
	"context"
	"time"

	redisCache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type RedisWrapper[T any] struct {
	exparationTime time.Duration
	cache          *redisCache.Cache
	prefix         string
}

func newRedisCache[T any](
	redisClient *redis.Client,
	defaultExparationTime time.Duration,
	ttl time.Duration,
	prefix string) *RedisWrapper[T] {
	return &RedisWrapper[T]{
		cache: redisCache.New(&redisCache.Options{
			Redis:      redisClient,
			LocalCache: redisCache.NewTinyLFU(10000, ttl),
		}),
		exparationTime: defaultExparationTime,
		prefix:         prefix,
	}
}

func (w *RedisWrapper[T]) Set(ctx context.Context, key string, value *T) error {
	return w.cache.Get()
}

func (w *RedisWrapper[T]) Get(ctx context.Context, key string) (value *T, err error) {
	return nil
}

func (w *RedisWrapper[T]) Delete(ctx context.Context, key string) error {
	return nil
}

func (w *RedisWrapper[T]) Exist() bool {
	return false
}

func (w *RedisWrapper[T]) ExparationTime() time.Duration {
	return time.Since(1)
}
