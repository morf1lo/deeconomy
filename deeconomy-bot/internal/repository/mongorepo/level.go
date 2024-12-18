package mongorepo

import (
	"context"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type levelRepo struct {
	levelsCollection *mongo.Collection
}

func newLevelRepo(db *mongo.Database) Level {
	return &levelRepo{
		levelsCollection: db.Collection("levels"),
	}
}

func (r *levelRepo) Create(ctx context.Context, level *model.Level) error {
	level.ID = primitive.NewObjectID()
	_, err := r.levelsCollection.InsertOne(ctx, level)
	return err
}

func (r *levelRepo) Update(ctx context.Context, level *model.Level) error {
	filter := bson.M{
		"userId": level.UserID,
		"guildId": level.GuildID,
	}

	update := bson.M{
		"$set": bson.M{
			"lvl": level.Lvl,
			"xp": level.XP,
		},
	}

	_, err := r.levelsCollection.UpdateOne(ctx, filter, update, options.Update())
	return err
}

func (r *levelRepo) FindByUserIDAndGuildID(ctx context.Context, userID string, guildID string) (*model.Level, error) {
	filter := bson.M{
		"userId": userID,
		"guildId": guildID,
	}

	var level model.Level
	if err := r.levelsCollection.FindOne(ctx, filter).Decode(&level); err != nil {
		return nil, err
	}
	return &level, nil
}

func (r *levelRepo) Delete(ctx context.Context, userID string, guildID string) error {
	filter := bson.M{
		"userId": userID,
		"guildId": guildID,
	}

	_, err := r.levelsCollection.DeleteOne(ctx, filter)
	return err
}
