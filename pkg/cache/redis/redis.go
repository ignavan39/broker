package cache

import (
	"context"
	"time"

	redisCache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

const defaultExpirationTime = time.Duration(24 * time.Hour)

type RedisCache[T any] struct {
	expirationTime time.Duration
	cache          *redisCache.Cache
	prefix         string
}

func NewRedisCache[T any](
	redisClient *redis.Client,
	ttl time.Duration,
	prefix string,
	size int,
) *RedisCache[T] {
	return &RedisCache[T]{
		cache: redisCache.New(&redisCache.Options{
			Redis:      redisClient,
			LocalCache: redisCache.NewTinyLFU(size, ttl),
		}),
		prefix: prefix,
		expirationTime: defaultExpirationTime,
	}
}

func (w *RedisCache[T]) WithExpirationTime(expirationTime time.Duration) *RedisCache[T] {
	w.expirationTime = expirationTime
	return w
}

func (w *RedisCache[T]) Set(ctx context.Context, key string, value T) error {
	return w.cache.Set(&redisCache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL: w.expirationTime,
	})
}

func (w *RedisCache[T]) Get(ctx context.Context, key string) (*T,error) {
	var value T
	err := w.cache.Get(ctx,key,&value)

	return &value, err
}

func (w *RedisCache[T]) Delete(ctx context.Context, key string) error {
	return w.cache.Delete(ctx,key)
}

func (w *RedisCache[T]) Exist(ctx context.Context, key string) bool {
	return w.cache.Exists(ctx,key)
}

func (w *RedisCache[T]) ExpirationTime() time.Duration {
	return w.expirationTime
}
