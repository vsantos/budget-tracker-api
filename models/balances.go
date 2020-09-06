package models

import (
	"budget-tracker/services"
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// CreateBalance creates a balance for a given owner_id
func CreateBalance(b Balance) (id string, err error) {
	dbClient, err := services.InitDatabase()
	if err != nil {
		return "", err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbBalanceCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	_, err = col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{
				{Key: "owner_id", Value: bsonx.Int32(1)},
				{Key: "month", Value: bsonx.Int32(1)},
				{Key: "year", Value: bsonx.Int32(1)},
			},
			Options: options.Index().SetUnique(true),
		},
	)

	// adding timestamp to creationDate
	t := time.Now()
	b.CreatedAt = primitive.NewDateTimeFromTime(t)
	b.SpendableAmount = b.Income.NetIncome
	b.Historic = []Spend{}

	r, err := col.InsertOne(ctx, b)
	if err != nil {
		cancel()
		return "", err
	}

	defer cancel()

	log.Infoln("created balance", r.InsertedID.(primitive.ObjectID).Hex())
	return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetBalance will return a balance from an owner_id based on month and year
func GetBalance(ownerID string, month int64, year int64) (balance *Balance, err error) {
	oid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return &Balance{}, err
	}

	dbClient, err := services.InitDatabase()
	if err != nil {
		return &Balance{}, err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbBalanceCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	err = col.FindOne(ctx, bson.M{
		"owner_id": oid,
		"month":    month,
		"year":     year,
	}).Decode(&balance)
	if err != nil {
		cancel()
		return &Balance{}, err
	}

	defer cancel()
	return balance, nil
}

// GetAllBalances will return all balances from an owner_id
func GetAllBalances(ownerID string) (balances []Balance, err error) {
	dbClient, err := services.InitDatabase()
	if err != nil {
		return []Balance{}, err
	}

	oid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return []Balance{}, err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbBalanceCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := col.Find(ctx, bson.M{"owner_id": oid})
	if err != nil {
		cancel()
		return []Balance{}, err
	}

	defer cursor.Close(ctx)
	defer cancel()

	for cursor.Next(ctx) {
		var balance Balance
		cursor.Decode(&balance)
		balances = append(balances, balance)
		defer cancel()
	}

	if err := cursor.Err(); err != nil {
		cancel()
		return []Balance{}, err
	}

	return balances, nil
}
