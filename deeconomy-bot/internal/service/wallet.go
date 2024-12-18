package service

import (
	"context"
	"time"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"github.com/morf1lo/deeconomy-bot/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type walletService struct {
	logger *zap.Logger
	repo *repository.Repository
}

func newWalletService(logger *zap.Logger, repo *repository.Repository) Wallet {
	return &walletService{
		logger: logger,
		repo: repo,
	}
}

func (s *walletService) Create(ctx context.Context, wallet *model.Wallet) error {
	return s.repo.Mongo.Wallet.Create(ctx, wallet)
}

func (s *walletService) FindByUserIDAndGuildID(ctx context.Context, userID string, guildID string) (*model.Wallet, error) {
	walletCache, err := s.repo.Redis.Wallet.Get(ctx, userID, guildID)
	if err == nil {
		return walletCache, nil
	}
	if err != redis.Nil {
		return nil, err
	}

	wallet, err := s.repo.Mongo.Wallet.FindByUserIDAndGuildID(ctx, userID, guildID)
	if err == mongo.ErrNoDocuments {
		wallet = &model.Wallet{
			UserID: userID,
			GuildID: guildID,
			Balance: 0,
			LastDailyCollected: time.Time{},
		}
		if err := s.repo.Mongo.Wallet.Create(ctx, wallet); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if err := s.repo.Redis.Wallet.Set(ctx, userID, guildID, wallet, time.Hour); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *walletService) Update(ctx context.Context, wallet *model.Wallet) error {
	if err := s.repo.Redis.Wallet.Del(ctx, wallet.UserID, wallet.GuildID); err != nil {
		return err
	}

	return s.repo.Mongo.Wallet.Update(ctx, wallet)
}

func (s *walletService) Delete(ctx context.Context, userID string, guildID string) error {
	if err := s.repo.Redis.Wallet.Del(ctx, userID, guildID); err != nil {
		return err
	}

	err := s.repo.Mongo.Wallet.Delete(ctx, userID, guildID)
	return err
}

func (s *walletService) IncrBy(ctx context.Context, userID string, guildID string, num int64) error {
	if err := s.repo.Mongo.Wallet.IncrBy(ctx, userID, guildID, num); err != nil {
		return err
	}

	err := s.repo.Redis.Wallet.Del(ctx, userID, guildID)
	return err
}

func (s *walletService) DecrBy(ctx context.Context, userID string, guildID string, num int64) error {
	if err := s.repo.Mongo.Wallet.DecrBy(ctx, userID, guildID, num); err != nil {
		return err
	}

	err := s.repo.Redis.Wallet.Del(ctx, userID, guildID)
	return err
}
