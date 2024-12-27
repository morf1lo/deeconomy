package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/morf1lo/deeconomy-bot-api/internal/model"
	"github.com/morf1lo/deeconomy-bot-api/internal/repository"
	"github.com/morf1lo/deeconomy-bot-api/internal/repository/redisrepo"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type guildService struct {
	logger *zap.Logger
	repo *repository.Repository
	httpClient *http.Client
}

func newGuildService(logger *zap.Logger, repo *repository.Repository) Guild {
	return &guildService{
		logger: logger,
		repo: repo,
		httpClient: &http.Client{},
	}
}

func (s *guildService) Create(ctx context.Context, guild *model.Guild) (*model.Guild, error) {
	result, err := s.repo.Mongo.Guild.Create(ctx, guild)
	if err != nil {
		s.logger.Sugar().Errorf("failed to CREATE guild: %s", err.Error())
	}
	return result, err
}

func (s *guildService) FindByGuildID(ctx context.Context, guildID string) (*model.Guild, error) {
	guildCache, err := s.repo.Redis.Guild.Get(ctx, guildID)
	if err == nil {
		return guildCache, nil
	}
	if err != redis.Nil {
		s.logger.Sugar().Errorf("failed to GET CACHED value from Redis: %s", err.Error())
		return nil, ErrInternal
	}

	guild, err := s.repo.Mongo.Guild.FindByGuildID(ctx, guildID)
	if err != nil {
		s.logger.Sugar().Errorf("failed to find guild(%s): %s", guildID, err.Error())
		return nil, ErrInternal
	}

	if err := s.repo.Redis.Guild.Set(ctx, guild, time.Hour); err != nil {
		s.logger.Sugar().Errorf("failed to SET guild(%s) in Redis: %s", guildID, err.Error())
		return nil, ErrInternal
	}

	return guild, nil
}

func (s *guildService) FindUserGuilds(ctx context.Context, discordID string, accessToken string) ([]*discordgo.Guild, error) {
	guildsCache, err := redisrepo.GetMany[discordgo.Guild](s.repo.Redis.Default, ctx, redisrepo.UserGuildsKey(discordID))
	if err == nil {
		return guildsCache, nil
	}
	if err != redis.Nil {
		s.logger.Sugar().Errorf("failed to get cached user(%s) guilds: %s", discordID, err.Error())
		return nil, ErrInternal
	}

	url := fmt.Sprintf("%s/users/@me/guilds", DISCORD_HOST)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		s.logger.Sugar().Errorf("failed to CREATE request: %s", err.Error())
		return nil, ErrInternal
	}

	req.Header.Add("Authotization", "Bearer " + accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Sugar().Errorf("failed to DO request to find user guilds: %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	var guilds []*discordgo.Guild
	if err := json.NewDecoder(resp.Body).Decode(&guilds); err != nil {
		s.logger.Sugar().Errorf("failed to decode response body: %s", err.Error())
		return nil, ErrInternal
	}

	if err := s.repo.Redis.Default.SetJSON(ctx, redisrepo.UserGuildsKey(discordID), guilds, time.Hour * 3); err != nil {
		s.logger.Sugar().Errorf("failed to set user(%s) guilds in Redis: %s", discordID, err.Error())
		return nil, ErrInternal
	}

	if len(guilds) == 0 {
		return nil, nil
	}

	return guilds, nil
}
