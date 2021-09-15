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
		Colletion: services.MongodbBalanceCollection,
	}

	_, err := m.SetIndex(
		context.Background(),
		bsonx.Doc{
			{Key: "owner_id", Value: bsonx.Int32(1)},
			{Key: "month", Value: bsonx.Int32(1)},
			{Key: "year", Value: bsonx.Int32(1)},
		},
		options.Index().SetUnique(true),
	)
	if err != nil {
		log.Error(err)
	}
}

// CreateBalance creates a balance for a given owner_id
func CreateBalance(parentCtx context.Context, b Balance) (id string, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("balance.owner.id").String(b.OwnerID.String()),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "CreateBalance", spanTags)
	defer span.End()

	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbBalanceCollection,
	}

	// adding timestamp to creationDate
	t := time.Now()
	b.CreatedAt = primitive.NewDateTimeFromTime(t)
	b.UpdatedAt = primitive.NewDateTimeFromTime(t)
	b.SpendableAmount = b.Income.NetIncome
	b.Historic = []Spend{}

	r, err := m.Create(ctx, b)
	if err != nil {
		return "", err
	}

	span.SetAttributes(attribute.Key("balance.id").String(r.InsertedID.(primitive.ObjectID).Hex()))

	observability.Metrics.Balances.BalancesCreated.Inc()
	log.Infoln("created balance", r.InsertedID.(primitive.ObjectID).Hex())
	return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetBalance will return a balance from an owner_id based on month and year
func GetBalance(parentCtx context.Context, ownerID string, month int64, year int64) (b *Balance, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("balance.owner.id").String(ownerID),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "GetBalance", spanTags)
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return &Balance{}, err
	}

	var balance Balance
	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbBalanceCollection,
	}

	d, err := m.Get(ctx, bson.M{
		"owner_id": oid,
		"month":    month,
		"year":     year,
	})
	if err != nil {
		return &Balance{}, err
	}

	err = d.Decode(&balance)
	if err != nil {
		return &Balance{}, err
	}

	span.SetAttributes(attribute.Key("balance.id").String(balance.ID.String()))
	return &balance, nil
}

// GetAllBalances will return all balances from an owner_id
func GetAllBalances(parentCtx context.Context, ownerID string) (balances []Balance, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("balance.owner.id").String(ownerID),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "GetAllBalances", spanTags)
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return []Balance{}, err
	}

	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbBalanceCollection,
	}

	cursor, err := m.GetAll(ctx, bson.M{"owner_id": oid})
	if err != nil {
		return []Balance{}, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var balance Balance
		cursor.Decode(&balance)
		balances = append(balances, balance)
	}

	if err := cursor.Err(); err != nil {
		return []Balance{}, err
	}

	return balances, nil
}
