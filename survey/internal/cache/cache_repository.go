package cache

import "context"

type CacheRepository interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key string, value string, ttlSeconds int) error
	Del(ctx context.Context, key string) error
}
