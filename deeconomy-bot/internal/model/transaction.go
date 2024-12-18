package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
	GuildID string `bson:"guildId" json:"guildId"`
	SenderID string `bson:"senderId" json:"senderId"`
	ReceiverID string `bson:"receiverId" json:"receiverId"`
	Amount int64 `bson:"amount" json:"amount"`
	ReducedAmount int64 `bson:"reducedAmount" json:"reducedAmount"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

func (t *Transaction) SetInsertTimestamp(now time.Time) {
	t.CreatedAt = now
}
