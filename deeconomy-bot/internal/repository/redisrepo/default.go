package redisrepo

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type defaultRepo struct {
	rdb *redis.Client
}

func newDefaultRepo(rdb *redis.Client) Default {
	return &defaultRepo{
		rdb: rdb,
	}
}

func (r *defaultRepo) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.rdb.Set(ctx, key, value, ttl).Err()
}

func (r *defaultRepo) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.rdb.Get(ctx, key)
}

func (r *defaultRepo) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.rdb.Del(ctx, keys...)
}
