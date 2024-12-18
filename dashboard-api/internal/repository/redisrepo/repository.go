package redisrepo

import (
	"context"
	"time"

	"github.com/morf1lo/deeconomy-bot-api/internal/model"
	"github.com/redis/go-redis/v9"
)

type Default interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type Guild interface {
	Set(ctx context.Context, guildID string, guild *model.Guild, ttl time.Duration) error
	Get(ctx context.Context, guildID string) (*model.Guild, error)
}

type Redis struct {
	Default
	Guild
}

func New(rdb *redis.Client) *Redis {
	return &Redis{
		Default: newDefaultRepo(rdb),
		Guild: newGuildRepo(rdb),
	}
}
