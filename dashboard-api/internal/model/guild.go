package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Guild struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
	GuildID string `bson:"guildId" json:"guildId"`
	LevelRewardingRoles map[int]string `bson:"levelRewardingRoles" json:"levelRewardingRoles"`
	ShopSettings map[string]int64 `bson:"shopSettings" json:"shopSettings"`
	DailyAmount int64 `bson:"dailyAmount" json:"dailyAmount"`
}
