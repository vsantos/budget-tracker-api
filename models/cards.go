package models

import (
	"budget-tracker-api/services"
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// CreateCard creates a card for a given owner_id
func CreateCard(c CreditCard) (id string, err error) {
	dbClient, err := services.InitDatabase()
	if err != nil {
		return "", err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbCardsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	_, err = col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{Key: "owner_id", Value: bsonx.Int32(1)}, {Key: "last_digits", Value: bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)

	// adding timestamp to creationDate
	t := time.Now()
	c.CreatedAt = primitive.NewDateTimeFromTime(t)

	r, err := col.InsertOne(ctx, c)
	if err != nil {
		cancel()
		return "", err
	}

	defer cancel()

	log.Infoln("created card", c.Alias)
	return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetAllCards will return a list of all cards from the database
func GetAllCards() (cards []CreditCard, err error) {
	dbClient, err := services.InitDatabase()
	if err != nil {
		return []CreditCard{}, err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbCardsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		cancel()
		return []CreditCard{}, err
	}

	defer cursor.Close(ctx)
	defer cancel()

	for cursor.Next(ctx) {
		var card CreditCard
		cursor.Decode(&card)
		cards = append(cards, card)
		defer cancel()
	}

	if err := cursor.Err(); err != nil {
		cancel()
		return []CreditCard{}, err
	}

	return cards, nil
}

// GetCards will return a list of cards from a owner_id
func GetCards(ownerID string) (cards []CreditCard, err error) {
	dbClient, err := services.InitDatabase()
	if err != nil {
		return []CreditCard{}, err
	}

	pid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return []CreditCard{}, err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbCardsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := col.Find(ctx, bson.M{"owner_id": pid})
	if err != nil {
		cancel()
		return []CreditCard{}, err
	}

	defer cursor.Close(ctx)
	defer cancel()

	for cursor.Next(ctx) {
		var card CreditCard
		cursor.Decode(&card)
		cards = append(cards, card)
		defer cancel()
	}

	if err := cursor.Err(); err != nil {
		cancel()
		return []CreditCard{}, err
	}

	return cards, nil
}

// DeleteCard creates an user based on request body payload
func DeleteCard(id string) (err error) {
	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	dbClient, err := services.InitDatabase()
	if err != nil {
		return err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbCardsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	log.Infoln("deleting card", id)
	result, err := col.DeleteOne(ctx, bson.M{"_id": pid})
	if err != nil {
		cancel()
		return err
	}

	log.Infoln("number of cards deleted:", result.DeletedCount)

	if result.DeletedCount == 0 {
		cancel()
		return errors.New("non existent card")
	}

	defer cancel()

	log.Infoln("deleted card", id)
	return nil
}
