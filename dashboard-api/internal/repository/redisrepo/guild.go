package redisrepo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/morf1lo/deeconomy-bot-api/internal/model"
	"github.com/redis/go-redis/v9"
)

type guildRepo struct {
	rdb *redis.Client
}

func newGuildRepo(rdb *redis.Client) Guild {
	return &guildRepo{
		rdb: rdb,
	}
}

func (r *guildRepo) Set(ctx context.Context, guildID string, guild *model.Guild, ttl time.Duration) error {
	guildJSON, err := json.Marshal(guild)
	if err != nil {
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
		return nil, err
	}

	return &guild, nil
}
