package cache

import (
	"context"
	"time"
)

type Cache[T any] interface {
	Set(ctx context.Context, key string, value T) error
	Get(ctx context.Context, key string) (value *T, err error)
	Delete(ctx context.Context, key string) error
	Exist(ctx context.Context, key string) bool
	ExpirationTime() time.Duration
}
