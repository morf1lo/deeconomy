package mongorepo

import (
	"context"

	"github.com/morf1lo/deeconomy-bot-api/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongo(ctx context.Context, cfg *config.MongoConfig) (*mongo.Database, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	db := client.Database(cfg.DBName)

	if _, err := db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "discordId", Value: 1},
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
