package redisrepo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type transactionRepo struct {
	logger *zap.Logger
	rdb *redis.Client
}

func newTransactionRepo(logger *zap.Logger, rdb *redis.Client) Transaction {
	return &transactionRepo{
		logger: logger,
		rdb: rdb,
	}
}

func (r *transactionRepo) SetMany(ctx context.Context, userID string, guildID string, scope string, transactions []*model.Transaction, ttl time.Duration) error {
	transactionsJSON, err := json.Marshal(transactions)
	if err != nil {
		r.logger.Sugar().Errorf("failed to marshal user(%s) transactions in guild(%s) to JSON: %s", userID, guildID, err.Error())
		return err
	}

	return r.rdb.Set(ctx, TransactionsKey(userID, guildID, scope), transactionsJSON, ttl).Err()
}

func (r *transactionRepo) GetMany(ctx context.Context, userID string, guildID string, scope string) ([]*model.Transaction, error) {
	transactionsCache, err := r.rdb.Get(ctx, TransactionsKey(userID, guildID, scope)).Result()
	if err != nil {
		return nil, err
	}

	var transactions []*model.Transaction
	if err := json.Unmarshal([]byte(transactionsCache), &transactions); err != nil {
		r.logger.Sugar().Errorf("failed to unmarshal cached user(%s) transactions in guild(%s) from JSON: %s", userID, guildID, err.Error())
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepo) Del(ctx context.Context, userID string, guildID string, scope string) error {
	return r.rdb.Del(ctx, TransactionsKey(userID, guildID, scope)).Err()
}
