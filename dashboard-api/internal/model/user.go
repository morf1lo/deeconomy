package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	DiscordID string             `bson:"discordId" json:"discordId"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
