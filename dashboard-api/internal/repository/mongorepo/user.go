package mongorepo

import (
	"context"
	"time"

	"github.com/morf1lo/deeconomy-bot-api/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	usersCollection *mongo.Collection
}

func newUserRepo(db *mongo.Database) User {
	return &userRepo{
		usersCollection: db.Collection("users"),
	}
}

func (r *userRepo) Create(ctx context.Context, user *model.User) (*model.User, error) {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	_, err := r.usersCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) UpdateByID(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	update := bson.M{
		"$set": bson.M{},
	}
	if value, exists := updates["discordAccessToken"]; exists {
		update["$set"].(bson.M)["discordAccessToken"] = value
	}
	if value, exists := updates["discordRefreshToken"]; exists {
		update["$set"].(bson.M)["discordRefreshToken"] = value
	}

	_, err := r.usersCollection.UpdateByID(ctx, id, update)
	return err
}

func (r *userRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	filter := bson.M{
		"_id": id,
	}

	var user model.User
	if err := r.usersCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) FindByDiscordID(ctx context.Context, discordID string) (*model.User, error) {
	filter := bson.M{
		"discordId": discordID,
	}

	var user model.User
	if err := r.usersCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
