package mongorepo

import (
	"context"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	FindUserTransactionsInGuild(ctx context.Context, customFilter bson.M) ([]*model.Transaction, error)
}

type Guild interface {
	FindByGuildID(ctx context.Context, guildID string) (*model.Guild, error)
}

type Mongo struct {
	Wallet
	Level
	Transaction
	Guild
}

func New(db *mongo.Database) *Mongo {
	return &Mongo{
		Wallet: newWalletRepo(db),
		Level: newLevelRepo(db),
		Transaction: newTransactionRepo(db),
		Guild: newGuildRepo(db),
	}
}
