package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
	UserID string `bson:"userId" json:"userId"`
	GuildID string `bson:"guildId" json:"guildId"`
	Balance int64 `bson:"balance" json:"balance"`
	LastDailyCollected time.Time `bson:"lastDailyCollected" json:"lastDailyCollected"`
}
