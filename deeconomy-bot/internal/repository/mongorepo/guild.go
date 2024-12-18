package mongorepo

import (
	"context"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type guildRepo struct {
	guildsCollection *mongo.Collection
}

func newGuildRepo(db *mongo.Database) Guild {
	return &guildRepo{
		guildsCollection: db.Collection("guilds"),
	}
}

func (r *guildRepo) FindByGuildID(ctx context.Context, guildID string) (*model.Guild, error) {
	filter := bson.M{
		"guildId": guildID,
	}

	var guild model.Guild
	if err := r.guildsCollection.FindOne(ctx, filter).Decode(&guild); err != nil {
		return nil, err
	}

	return &guild, nil
}
