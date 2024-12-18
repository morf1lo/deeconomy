package service

import (
	"context"
	"time"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"github.com/morf1lo/deeconomy-bot/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type levelService struct {
	logger *zap.Logger
	repo *repository.Repository
}

func newLevelService(logger *zap.Logger, repo *repository.Repository) Level {
	return &levelService{
		logger: logger,
		repo: repo,
	}
}

func (s *levelService) Create(ctx context.Context, level *model.Level) error {
	return s.repo.Mongo.Level.Create(ctx, level)
}

func (s *levelService) Update(ctx context.Context, level *model.Level) error {
	if err := s.repo.Mongo.Level.Update(ctx, level); err != nil {
		return err
	}

	err := s.repo.Redis.Level.Del(ctx, level.UserID, level.GuildID)
	return err
}

func (s *levelService) FindByUserIDAndGuildID(ctx context.Context, userID string, guildID string) (*model.Level, error) {
	levelCache, err := s.repo.Redis.Level.Get(ctx, userID, guildID)
	if err == nil {
		return levelCache, nil
	}
	if err != redis.Nil {
		return nil, err
	}

	level, err := s.repo.Mongo.Level.FindByUserIDAndGuildID(ctx, userID, guildID)
	if err == mongo.ErrNoDocuments {
		level = &model.Level{
			UserID: userID,
			GuildID: guildID,
			Lvl: 1,
			XP: 0,
		}
		if err := s.repo.Mongo.Level.Create(ctx, level); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if err := s.repo.Redis.Level.Set(ctx, userID, guildID, level, time.Hour); err != nil {
		return nil, err
	}

	return level, nil
}

func (s *levelService) Delete(ctx context.Context, userID string, guildID string) error {
	if err := s.repo.Mongo.Level.Delete(ctx, userID, guildID); err != nil {
		return err
	}

	err := s.repo.Redis.Level.Del(ctx, userID, guildID)
	return err
}
