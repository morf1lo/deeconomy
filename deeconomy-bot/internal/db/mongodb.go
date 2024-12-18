package db

import (
	"context"

	"github.com/morf1lo/deeconomy-bot/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB(ctx context.Context, cfg *config.MongoDBConfig) (*mongo.Database, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	db := client.Database(cfg.DBName)

	if _, err := db.Collection("wallets").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "userId", Value: 1},
			{Key: "guildId", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, err
	}

	if _, err := db.Collection("levels").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "userId", Value: 1},
			{Key: "guildId", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, err
	}

	if _, err := db.Collection("guilds").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "guildId", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, err
	}

	return db, nil
}