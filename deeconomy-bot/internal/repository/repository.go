package repository

import (
	"github.com/morf1lo/deeconomy-bot/internal/repository/mongorepo"
	"github.com/morf1lo/deeconomy-bot/internal/repository/redisrepo"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Repository struct {
	*mongorepo.Mongo
	*redisrepo.Redis
}

func New(logger *zap.Logger, db *mongo.Database, rdb *redis.Client) *Repository {
	return &Repository{
		Mongo: mongorepo.New(db),
		Redis: redisrepo.New(logger, rdb),
	}
}
