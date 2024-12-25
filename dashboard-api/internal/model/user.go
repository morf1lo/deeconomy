package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
	DiscordID string `bson:"discordId" json:"discordId"`
	DiscordAccessToken string `bson:"discordAccessToken" json:"discordAccessToken"`
	DiscordRefreshToken string `bson:"discordRefreshToken" json:"discordRefreshToken"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
