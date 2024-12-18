package service

import (
	"context"
	"time"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"github.com/morf1lo/deeconomy-bot/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type transactionService struct {
	logger *zap.Logger
	repo *repository.Repository
}

func newTransactionService(logger *zap.Logger, repo *repository.Repository) Transaction {
	return &transactionService{
		logger: logger,
		repo: repo,
	}
}

func (s *transactionService) New(ctx context.Context, transaction *model.Transaction) error {
	go func(transaction *model.Transaction) {
		if err := s.clearAllCache(ctx, transaction); err != nil {
			s.logger.Sugar().Errorf("failed to clear cache for transactions: %s", err.Error())
		}
	}(transaction)

	return s.repo.Mongo.Transaction.New(ctx, transaction)
}

func (s *transactionService) clearAllCache(ctx context.Context, transaction *model.Transaction) error {
	if err := s.repo.Redis.Transaction.Del(ctx, transaction.SenderID, transaction.GuildID, DefaultScope); err != nil {
		return err
	}

	if err := s.repo.Redis.Transaction.Del(ctx, transaction.ReceiverID, transaction.GuildID, DefaultScope); err != nil {
		return err
	}

	if err := s.repo.Redis.Transaction.Del(ctx, transaction.SenderID, transaction.GuildID, SenderOnlyScope); err != nil {
		return err
	}

	if err := s.repo.Redis.Transaction.Del(ctx, transaction.ReceiverID, transaction.GuildID, SenderOnlyScope); err != nil {
		return err
	}

	if err := s.repo.Redis.Transaction.Del(ctx, transaction.SenderID, transaction.GuildID, ReceiverOnlyScope); err != nil {
		return err
	}

	if err := s.repo.Redis.Transaction.Del(ctx, transaction.ReceiverID, transaction.GuildID, ReceiverOnlyScope); err != nil {
		return err
	}

	return nil
}

func (s *transactionService) FindUserTransactionsInGuild(ctx context.Context, userID string, guildID string, scope string) ([]*model.Transaction, error) {
	transactionsCache, err := s.repo.Redis.Transaction.GetMany(ctx, userID, guildID, scope)
	if err == nil {
		return transactionsCache, nil
	}
	if err != redis.Nil {
		return nil, err
	}

	filter := bson.M{
		"guildId": guildID,
	}

	if scope == SenderOnlyScope {
		filter["senderId"] = userID
	} else if scope == ReceiverOnlyScope {
		filter["receiverId"] = userID
	} else if scope == DefaultScope {
		filter["$or"] = []bson.M{
			{
				"senderId": userID,
			},
			{
				"receiverId": userID,
			},
		}
	} else {
		return nil, errScopeMustBeProvided
	}

	transactions, err := s.repo.Mongo.Transaction.FindUserTransactionsInGuild(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Redis.Transaction.SetMany(ctx, userID, guildID, scope, transactions, time.Hour * 3); err != nil {
		return nil, err
	}

	return transactions, nil
}
