package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Level struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
	UserID string `bson:"userId" json:"userId"`
	GuildID string `bson:"guildId" json:"guildId"`
	Lvl int `bson:"lvl" json:"lvl"`
	XP int `bson:"xp" json:"xp"`
}
