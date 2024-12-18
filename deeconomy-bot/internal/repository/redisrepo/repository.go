package redisrepo

import (
	"context"
	"time"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Default interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type Wallet interface {
	Set(ctx context.Context, userID string, guildID string, wallet *model.Wallet, ttl time.Duration) error
	Get(ctx context.Context, userID string, guildID string) (*model.Wallet, error)
	Del(ctx context.Context, userID string, guildID string) error
}

type Level interface {
	Set(ctx context.Context, userID string, guildID string, level *model.Level, ttl time.Duration) error
	Get(ctx context.Context, userID string, guildID string) (*model.Level, error)
	Del(ctx context.Context, userID string, guildID string) error
}

type Transaction interface {
	SetMany(ctx context.Context, userID string, guildID string, scope string, transactions []*model.Transaction, ttl time.Duration) error
	GetMany(ctx context.Context, userID string, guildID string, scope string) ([]*model.Transaction, error)
	Del(ctx context.Context, userID string, guildID string, scope string) error
}

type Guild interface {
	Set(ctx context.Context, guildID string, guild *model.Guild, ttl time.Duration) error
	Get(ctx context.Context, guildID string) (*model.Guild, error)
}

type Redis struct {
	Default
	Wallet
	Level
	Transaction
	Guild
}

func New(logger *zap.Logger, rdb *redis.Client) *Redis {
	return &Redis{
		Default: newDefaultRepo(rdb),
		Wallet: newWalletRepo(logger, rdb),
		Level: newLevelRepo(logger, rdb),
		Transaction: newTransactionRepo(logger, rdb),
		Guild: newGuildRepo(logger, rdb),
	}
}
