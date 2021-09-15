package models

import (
	"budget-tracker-api/observability"
	"budget-tracker-api/services"
	"context"
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
		Colletion: services.MongodbSpendsCollection,
	}

	_, err := m.SetIndex(
		context.Background(),
		bsonx.Doc{
			{Key: "owner_id", Value: bsonx.Int32(1)},
		},
		options.Index().SetUnique(true),
	)
	if err != nil {
		log.Error(err)
	}
}

// CreateSpend creates a card for a given owner_id
func CreateSpend(parentCtx context.Context, s Spend) (id string, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("spend.owner.id").String(s.OwnerID.String()),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "CreateSpend", spanTags)
	defer span.End()

	// adding timestamp to creationDate
	t := time.Now()
	s.CreatedAt = primitive.NewDateTimeFromTime(t)

	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbSpendsCollection,
	}

	r, err := m.Create(ctx, s)
	if err != nil {
		return "", err
	}

	span.SetAttributes(attribute.Key("spend.id").String(r.InsertedID.(primitive.ObjectID).Hex()))

	observability.Metrics.Spends.SpendsCreated.Inc()
	log.Infoln("created spend", r.InsertedID.(primitive.ObjectID).Hex())
	return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetSpends will return all spends from a specific owner_id
func GetSpends(parentCtx context.Context, ownerID string) (spends []Spend, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("spend.owner.id").String(ownerID),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "GetSpends", spanTags)
	defer span.End()

	pid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return []Spend{}, err
	}

	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbSpendsCollection,
	}

	cursor, err := m.GetAll(ctx, bson.M{"owner_id": pid})
	if err != nil {
		return []Spend{}, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var spend Spend
		cursor.Decode(&spend)
		spends = append(spends, spend)
	}

	if err := cursor.Err(); err != nil {
		return []Spend{}, err
	}

	return spends, nil
}
