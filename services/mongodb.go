package services

import (
	"budget-tracker-api/observability"
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.opentelemetry.io/otel/attribute"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	// MongodbURI will define a URI to be used across packages
	MongodbURI = "mongodb://localhost:27017/"
	// MongodbDatabase will define a database to be used across packages
	MongodbDatabase = "budget-tracker"
	// MongodbUserCollection will define a URI to be used across packages
	MongodbUserCollection = "users"
	// MongodbCardsCollection will define a User collection
	MongodbCardsCollection = "cards"
	// MongodbBalanceCollection will define a Cards collection
	MongodbBalanceCollection = "balance"
	// MongodbSpendsCollection will define a Spend collection
	MongodbSpendsCollection = "spends"
)

var (
	// MongoClient will define a mongoDB client to be initialized and used across pkgs
	MongoClient *mongo.Client
)

// MongoCfg satisfies DataManager and Monger Interfaces
type MongoCfg struct {
	// Database URI for mongodb. Example: "mongodb://localhost:2701"
	URI string
	// Database name for MongoDB
	Database string
	// Database Collection name for MongoDB
	Colletion string
}

// setIndexes will set all indexes necessary for MongoDB client
func setIndexes(ctx context.Context, c *mongo.Client) error {
	_, err := setIndex(
		ctx,
		c,
		MongodbDatabase,
		MongodbUserCollection,
		bsonx.Doc{{Key: "login", Value: bsonx.Int32(1)}},
		options.Index().SetUnique(true),
	)
	if err != nil {
		return err
	}

	_, err = setIndex(
		ctx,
		c,
		MongodbDatabase,
		MongodbCardsCollection,
		bsonx.Doc{
			{Key: "owner_id", Value: bsonx.Int32(1)},
			{Key: "last_digits", Value: bsonx.Int32(1)},
		},
		options.Index().SetUnique(true),
	)
	if err != nil {
		return err
	}

	_, err = setIndex(
		ctx,
		c,
		MongodbDatabase,
		MongodbBalanceCollection,
		bsonx.Doc{
			{Key: "owner_id", Value: bsonx.Int32(1)},
			{Key: "month", Value: bsonx.Int32(1)},
			{Key: "year", Value: bsonx.Int32(1)},
		},
		options.Index().SetUnique(true),
	)
	if err != nil {
		return err
	}

	return nil
}

// setIndex will performance a single Indexes().CreateOne operation based on a index name
func setIndex(ctx context.Context, c *mongo.Client, database string, collection string, keys bsonx.Doc, opts *options.IndexOptions) (index string, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("index.collection.name").String(collection),
	}

	ctx, span := observability.Span(ctx, "mongodb", "CreateIndex", spanTags)
	defer span.End()

	ind := c.Database(database).Collection(collection).Indexes()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	i, err := ind.CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    keys,
			Options: opts,
		},
	)

	log.Debugln("created index", i)
	span.SetAttributes(attribute.Key("index.name").String(i))

	defer cancel()

	return i, nil
}

// InitDatabaseWithURI will return a database client for usage
func InitDatabaseWithURI(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.Monitor = otelmongo.NewMonitor("mongodb")
	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cancel()
		return &mongo.Client{}, err
	}

	err = setIndexes(ctx, c)
	if err != nil {
		cancel()
		return &mongo.Client{}, err
	}

	defer cancel()

	return c, err
}

// Health will return if the mongoDB is healthy
func (m MongoCfg) Health() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	err = MongoClient.Ping(ctx, &readpref.ReadPref{})
	if err != nil {
		cancel()
		return err
	}

	defer cancel()
	return nil
}

// Get will perform a mongoDB FindOne operation
func (m MongoCfg) Get(ctx context.Context, filter interface{}) (r *mongo.SingleResult, err error) {
	col := MongoClient.Database(m.Database).Collection(m.Colletion)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	r = col.FindOne(ctx, filter)
	if r.Err() != nil {
		cancel()
		return &mongo.SingleResult{}, r.Err()
	}

	defer cancel()
	return r, nil
}

// GetAll will perform a mongoDB Find operation
func (m MongoCfg) GetAll(ctx context.Context, filter interface{}) (r *mongo.Cursor, err error) {
	col := MongoClient.Database(m.Database).Collection(m.Colletion)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	r, err = col.Find(ctx, filter)
	if err != nil {
		cancel()
		return r, err
	}

	defer cancel()

	return r, nil
}

// Create will perform a mongoDB InsertOne operation
func (m MongoCfg) Create(ctx context.Context, filter interface{}) (r *mongo.InsertOneResult, err error) {
	col := MongoClient.Database(m.Database).Collection(m.Colletion)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	r, err = col.InsertOne(ctx, filter)
	if err != nil {
		cancel()
		return r, err
	}
	defer cancel()

	return r, nil
}

// Delete will perform a mongoDB DeleteOne operation
func (m MongoCfg) Delete(ctx context.Context, filter interface{}) (r *mongo.DeleteResult, err error) {

	col := MongoClient.Database(m.Database).Collection(m.Colletion)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	r, err = col.DeleteOne(ctx, filter)
	if err != nil {
		cancel()
		return r, err
	}
	defer cancel()

	return r, nil
}
