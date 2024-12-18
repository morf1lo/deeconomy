package redisrepo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type walletRepo struct {
	logger *zap.Logger
	rdb *redis.Client
}

func newWalletRepo(logger *zap.Logger, rdb *redis.Client) Wallet {
	return &walletRepo{
		logger: logger,
		rdb: rdb,
	}
}

func (r *walletRepo) Set(ctx context.Context, userID string, guildID string, wallet *model.Wallet, ttl time.Duration) error {
	walletJSON, err := json.Marshal(wallet)
	if err != nil {
		r.logger.Sugar().Errorf("failed to marshal wallet info to JSON with ID(%d): %s", wallet.ID.Hex(), err.Error())
		return err
	}

	return r.rdb.Set(ctx, WalletKey(userID, guildID), walletJSON, ttl).Err()
}

func (r *walletRepo) Get(ctx context.Context, userID string, guildID string) (*model.Wallet, error) {
	walletCache, err := r.rdb.Get(ctx, WalletKey(userID, guildID)).Result()
	if err != nil {
		return nil, err
	}

	var wallet model.Wallet
	if err := json.Unmarshal([]byte(walletCache), &wallet); err != nil {
		r.logger.Sugar().Errorf("failed to unmarshal cached wallet info: %s", err.Error())
		return nil, err
	}

	return &wallet, nil
}

func (r *walletRepo) Del(ctx context.Context, userID string, guildID string) error {
	return r.rdb.Del(ctx, WalletKey(userID, guildID)).Err()
}
