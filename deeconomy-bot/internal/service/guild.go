package service

import (
	"context"
	"time"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"github.com/morf1lo/deeconomy-bot/internal/repository"
	"github.com/redis/go-redis/v9"
)

type guildService struct {
	repo *repository.Repository
}

func newGuildService(repo *repository.Repository) Guild {
	return &guildService{
		repo: repo,
	}
}

func (s *guildService) FindByGuildID(ctx context.Context, guildID string) (*model.Guild, error) {
	guildCache, err := s.repo.Redis.Guild.Get(ctx, guildID)
	if err == nil {
		return guildCache, nil
	}

	if err != redis.Nil {
		return nil, err
	}

	guildDB, err := s.repo.Mongo.Guild.FindByGuildID(ctx, guildID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Redis.Guild.Set(ctx, guildID, guildDB, time.Hour * 4); err != nil {
		return nil, err
	}

	return guildDB, nil
}
