package models

import (
	"budget-tracker-api/observability"
	"budget-tracker-api/services"
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.opentelemetry.io/otel/attribute"
)

func init() {
	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbCardsCollection,
	}

	_, err := m.SetIndex(
		context.Background(),
		bsonx.Doc{
			{Key: "owner_id", Value: bsonx.Int32(1)},
			{Key: "last_digits", Value: bsonx.Int32(1)},
		},
		options.Index().SetUnique(true),
	)
	if err != nil {
		log.Error(err)
	}
}

// CreateCard creates a card for a given owner_id
func CreateCard(parentCtx context.Context, c CreditCard) (id string, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("card.owner.id").String(c.OwnerID.String()),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "CreateCard", spanTags)
	defer span.End()

	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbCardsCollection,
	}

	// adding timestamp to creationDate
	t := time.Now()
	c.CreatedAt = primitive.NewDateTimeFromTime(t)

	r, err := m.Create(ctx, c)
	if err != nil {
		return "", err
	}

	span.SetAttributes(attribute.Key("card.id").String(r.InsertedID.(primitive.ObjectID).Hex()))

	observability.Metrics.Cards.CardsCreated.Inc()
	log.Infoln("created card", c.Alias)
	return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetAllCards will return a list of all cards from the database
func GetAllCards(parentCtx context.Context) (cards []CreditCard, err error) {
	ctx, span := observability.Span(parentCtx, "mongodb", "GetAllCards", []attribute.KeyValue{})
	defer span.End()

	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbCardsCollection,
	}

	cursor, err := m.GetAll(ctx, bson.M{})
	if err != nil {
		return []CreditCard{}, err
	}
	if err != nil {
		return []CreditCard{}, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var card CreditCard
		cursor.Decode(&card)
		cards = append(cards, card)
	}

	if err := cursor.Err(); err != nil {
		return []CreditCard{}, err
	}

	return cards, nil
}

// GetCards will return a list of cards from a owner_id
func GetCards(parentCtx context.Context, ownerID string) (cards []CreditCard, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("card.owner.id").String(ownerID),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "GetUserCards", spanTags)
	defer span.End()

	pid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return []CreditCard{}, err
	}

	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbCardsCollection,
	}

	cursor, err := m.GetAll(ctx, bson.M{"owner_id": pid})
	if err != nil {
		return []CreditCard{}, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var card CreditCard
		cursor.Decode(&card)
		cards = append(cards, card)
	}

	if err := cursor.Err(); err != nil {
		return []CreditCard{}, err
	}

	return cards, nil
}

// DeleteCard creates an user based on request body payload
func DeleteCard(parentCtx context.Context, id string) (err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("card.id").String(id),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "DeleteUserCard", spanTags)
	defer span.End()

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	log.Infoln("deleting card", id)
	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbCardsCollection,
	}

	r, err := m.Delete(ctx, bson.M{"_id": pid})
	if err != nil {
		return err
	}

	log.Infoln("number of cards deleted:", r.DeletedCount)
	if r.DeletedCount == 0 {
		return errors.New("non existent card")
	}

	log.Infoln("deleted card", id)
	return nil
}
