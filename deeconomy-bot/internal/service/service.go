package service

import (
	"context"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"github.com/morf1lo/deeconomy-bot/internal/repository"
	"go.uber.org/zap"
)

type Wallet interface {
	Create(ctx context.Context, wallet *model.Wallet) error
	FindByUserIDAndGuildID(ctx context.Context, userID string, guildID string) (*model.Wallet, error)
	Update(ctx context.Context, wallet *model.Wallet) error
	Delete(ctx context.Context, userID string, guildID string) error
	IncrBy(ctx context.Context, userID string, guildID string, num int64) error
	DecrBy(ctx context.Context, userID string, guildID string, num int64) error
}

type Level interface {
	Create(ctx context.Context, level *model.Level) error
	Update(ctx context.Context, level *model.Level) error
	FindByUserIDAndGuildID(ctx context.Context, userID string, guildID string) (*model.Level, error)
	Delete(ctx context.Context, userID string, guildID string) error
}

type Transaction interface {
	New(ctx context.Context, transaction *model.Transaction) error
	FindUserTransactionsInGuild(ctx context.Context, userID string, guildID string, scope string) ([]*model.Transaction, error)
}

type Guild interface {
	FindByGuildID(ctx context.Context, guildID string) (*model.Guild, error)
}

type Service struct {
	Wallet
	Level
	Transaction
	Guild
}

func New(logger *zap.Logger, repo *repository.Repository) *Service {
	return &Service{
		Wallet: newWalletService(logger, repo),
		Level: newLevelService(logger, repo),
		Transaction: newTransactionService(logger, repo),
		Guild: newGuildService(repo),
	}
}
