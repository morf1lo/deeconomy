package mongorepo

import (
	"context"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type walletRepo struct {
	walletsCollection *mongo.Collection
}

func newWalletRepo(db *mongo.Database) Wallet {
	return &walletRepo{
		walletsCollection: db.Collection("wallets"),
	}
}

func (r *walletRepo) Create(ctx context.Context, wallet *model.Wallet) error {
	wallet.ID = primitive.NewObjectID()
	_, err := r.walletsCollection.InsertOne(ctx, wallet)
	return err
}

func (r *walletRepo) FindByUserIDAndGuildID(ctx context.Context, userID string, guildID string) (*model.Wallet, error) {
	filter := bson.M{
		"userId": userID,
		"guildId": guildID,
	}

	var wallet model.Wallet
	if err := r.walletsCollection.FindOne(ctx, filter).Decode(&wallet); err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepo) Update(ctx context.Context, wallet *model.Wallet) error {
	filter := bson.M{
		"userId": wallet.UserID,
		"guildId": wallet.GuildID,
	}

	update := bson.M{
		"$set": bson.M{
			"userId": wallet.UserID,
			"guildId": wallet.GuildID,
			"balance": wallet.Balance,
			"lastDailyCollected": wallet.LastDailyCollected,
		},
	}

	_, err := r.walletsCollection.UpdateOne(ctx, filter, update)
	return err
}

func (r *walletRepo) Delete(ctx context.Context, userID string, guildID string) error {
	filter := bson.M{
		"userId": userID,
		"guildId": guildID,
	}

	_, err := r.walletsCollection.DeleteOne(ctx, filter)
	return err
}

func (r *walletRepo) IncrBy(ctx context.Context, userID string, guildID string, num int64) error {
	if num < 1 {
		return errInvalidNumber
	}

	filter := bson.M{
		"userId": userID,
		"guildId": guildID,
	}

	update := bson.M{
		"$inc": bson.M{
			"balance": num,
		},
	}

	_, err := r.walletsCollection.UpdateOne(ctx, filter, update)
	return err
}

func (r *walletRepo) DecrBy(ctx context.Context, userID string, guildID string, num int64) error {
	if num < 1 {
		return errInvalidNumber
	}

	filter := bson.M{
		"userId": userID,
		"guildId": guildID,
	}

	update := bson.M{
		"$inc": bson.M{
			"balance": -num,
		},
	}

	_, err := r.walletsCollection.UpdateOne(ctx, filter, update)
	return err
}
