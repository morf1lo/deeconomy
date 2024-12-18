package redisrepo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type guildRepo struct {
	logger *zap.Logger
	rdb *redis.Client
}

func newGuildRepo(logger *zap.Logger, rdb *redis.Client) Guild {
	return &guildRepo{
		logger: logger,
		rdb: rdb,
	}
}

func (r *guildRepo) Set(ctx context.Context, guildID string, guild *model.Guild, ttl time.Duration) error {
	guildJSON, err := json.Marshal(guild)
	if err != nil {
		r.logger.Sugar().Errorf("failed to marshal guild to JSON with ID(%s): %s", guild.ID.Hex(), err.Error())
		return err
	}

	return r.rdb.Set(ctx, GuildKey(guildID), guildJSON, ttl).Err()
}

func (r *guildRepo) Get(ctx context.Context, guildID string) (*model.Guild, error) {
	guildCache, err := r.rdb.Get(ctx, GuildKey(guildID)).Result()
	if err != nil {
		return nil, err
	}

	var guild model.Guild
	if err := json.Unmarshal([]byte(guildCache), &guild); err != nil {
		r.logger.Sugar().Errorf("failed to unmarshal cached guild with guildID(%s): %s", guildID, err.Error())
		return nil, err
	}

	return &guild, nil
}
