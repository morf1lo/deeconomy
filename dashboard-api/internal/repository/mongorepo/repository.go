package mongorepo

import (
	"context"

	"github.com/morf1lo/deeconomy-bot-api/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
	FindByDiscordID(ctx context.Context, discordID string) (*model.User, error)
}

type Guild interface {
	Create(ctx context.Context, guild *model.Guild) (*model.Guild, error)
	FindByGuildID(ctx context.Context, guildID string) (*model.Guild, error)
}

type Mongo struct {
	User
	Guild
}

func New(db *mongo.Database) *Mongo {
	return &Mongo{
		User: newUserRepo(db),
		Guild: newGuildRepo(db),
	}
}
