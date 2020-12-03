package models

import (
	"budget-tracker/services"
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateSpend creates a card for a given owner_id
func CreateSpend(s Spend) (id string, err error) {
	dbClient, err := services.InitDatabase()
	if err != nil {
		return "", err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbSpendsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// adding timestamp to creationDate
	t := time.Now()
	s.CreatedAt = primitive.NewDateTimeFromTime(t)

	r, err := col.InsertOne(ctx, s)
	if err != nil {
		cancel()
		return "", err
	}

	defer cancel()

	log.Infoln("created spend", r.InsertedID.(primitive.ObjectID).Hex())
	return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetSpends will return all spends from a specific owner_id
func GetSpends(ownerID string) (spends []Spend, err error) {
	dbClient, err := services.InitDatabase()
	if err != nil {
		return []Spend{}, err
	}

	pid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return []Spend{}, err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbSpendsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := col.Find(ctx, bson.M{"owner_id": pid})
	if err != nil {
		cancel()
		return []Spend{}, err
	}

	defer cursor.Close(ctx)
	defer cancel()

	for cursor.Next(ctx) {
		var spend Spend
		cursor.Decode(&spend)
		spends = append(spends, spend)
		defer cancel()
	}

	if err := cursor.Err(); err != nil {
		cancel()
		return []Spend{}, err
	}

	return spends, nil
}
