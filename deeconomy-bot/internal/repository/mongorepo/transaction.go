package mongorepo

import (
	"context"
	"time"

	"github.com/morf1lo/deeconomy-bot/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type transactionRepo struct {
	transactionsCollection *mongo.Collection
}

func newTransactionRepo(db *mongo.Database) Transaction {
	return &transactionRepo{
		transactionsCollection: db.Collection("transactions"),
	}
}

func (r *transactionRepo) New(ctx context.Context, transaction *model.Transaction) error {
	transaction.ID = primitive.NewObjectID()
	transaction.SetInsertTimestamp(time.Now())
	_, err := r.transactionsCollection.InsertOne(ctx, transaction)
	return err
}

func (r *transactionRepo) FindUserTransactionsInGuild(ctx context.Context, customFilter bson.M) ([]*model.Transaction, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{
		{
			Key: "createdAt",
			Value: -1,
		},
	})
	findOptions.SetLimit(3)
	
	cursor, err := r.transactionsCollection.Find(ctx, customFilter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*model.Transaction
	for cursor.Next(ctx) {
		var transaction model.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
