package repository

import (
	"github.com/morf1lo/deeconomy-bot-api/internal/repository/mongorepo"
	"github.com/morf1lo/deeconomy-bot-api/internal/repository/redisrepo"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	*mongorepo.Mongo
	*redisrepo.Redis
}

func New(db *mongo.Database, rdb *redis.Client) *Repository {
	return &Repository{
		Mongo: mongorepo.New(db),
		Redis: redisrepo.New(rdb),
	}
}
