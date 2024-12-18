package service

import (
	"context"

	"github.com/morf1lo/deeconomy-bot-api/internal/model"
	"github.com/morf1lo/deeconomy-bot-api/internal/repository"
)

type User interface {
}

type Guild interface {
	Create(ctx context.Context, guild *model.Guild) (*model.Guild, error)
	FindByGuildID(ctx context.Context, guildID string) (*model.Guild, error)
	FindUserGuilds(ctx context.Context, discordID string) ([]*model.Guild, error)
}

type Service struct {
	User
	Guild
}

func New(repo *repository.Repository) *Service {
	return &Service{
		
	}
}
