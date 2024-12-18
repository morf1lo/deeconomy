package redisrepo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type levelRepo struct {
	logger *zap.Logger
	rdb *redis.Client
}

func newLevelRepo(logger *zap.Logger, rdb *redis.Client) Level {
	return &levelRepo{
		logger: logger,
		rdb: rdb,
	}
}

func (r *levelRepo) Set(ctx context.Context, userID string, guildID string, level *model.Level, ttl time.Duration) error {
	levelJSON, err := json.Marshal(level)
	if err != nil {
		r.logger.Sugar().Errorf("failed to marshal level info to JSON with ID(%d): %s", level.ID.Hex(), err.Error())
		return err
	}

	return r.rdb.Set(ctx, LevelKey(userID, guildID), levelJSON, ttl).Err()
}

func (r *levelRepo) Get(ctx context.Context, userID string, guildID string) (*model.Level, error) {
	levelCache, err := r.rdb.Get(ctx, LevelKey(userID, guildID)).Result()
	if err != nil {
		return nil, err
	}

	var level model.Level
	if err := json.Unmarshal([]byte(levelCache), &level); err != nil {
		r.logger.Sugar().Errorf("failed to unmarshal cached level info: %s", err.Error())
		return nil, err
	}

	return &level, nil
}

func (r *levelRepo) Del(ctx context.Context, userID string, guildID string) error {
	return r.rdb.Del(ctx, LevelKey(userID, guildID)).Err()
}
