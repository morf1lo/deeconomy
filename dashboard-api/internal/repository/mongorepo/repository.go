package mongorepo

import (
	"context"

	"github.com/morf1lo/deeconomy-bot-api/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type Guild interface {
	Create(ctx context.Context, guild *model.Guild) error
	FindByGuildID(ctx context.Context, guildID string) (*model.Guild, error)
}

type Mongo struct {
	Guild
}

func New(db *mongo.Database) *Mongo {
	return &Mongo{
		Guild: newGuildRepo(db),
	}
}
