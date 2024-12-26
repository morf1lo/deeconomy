package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/morf1lo/deeconomy-bot-api/internal/model"
	"github.com/morf1lo/deeconomy-bot-api/internal/repository"
	"go.uber.org/zap"
)

type User interface {
	Authorize(ctx context.Context, discordCode string) (string, string, error)
}

type Guild interface {
	Create(ctx context.Context, guild *model.Guild) (*model.Guild, error)
	FindByGuildID(ctx context.Context, guildID string) (*model.Guild, error)
	FindUserGuilds(ctx context.Context, discordID string, accessToken string) ([]*discordgo.Guild, error)
}

type Service struct {
	User
	Guild
}

func New(logger *zap.Logger, repo *repository.Repository) *Service {
	return &Service{
		User: newUserService(logger, repo),
		Guild: newGuildService(logger, repo),
	}
}
