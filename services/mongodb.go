package services

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// DatabaseURI will set where to store data
	DatabaseURI = "mongodb://localhost:27017/"
	// MongodbDatabase will define a default mongo database
	MongodbDatabase = "budget-tracker"
	// MongodbUserCollection will define a default users collection
	MongodbUserCollection = "users"
	// MongodbCardsCollection will define a default cards collection
	MongodbCardsCollection = "cards"
	// MongodbBalanceCollection will define a default balance collection
	MongodbBalanceCollection = "balance"
	// MongodbSpendsCollection will define a default spends collection
	MongodbSpendsCollection = "spends"
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

// Monger will
type Monger interface {
	SetIndex(ctx context.Context, keys interface{}, opts *options.IndexOptions) (index string, err error)
	Get(ctx context.Context, document primitive.M) (s *mongo.SingleResult, err error)
	GetAll(ctx context.Context, document primitive.M) (s *mongo.Cursor, err error)
	Create(ctx context.Context, document interface{}) (r *mongo.InsertOneResult, err error)
	Delete(ctx context.Context, document primitive.M) (r *mongo.DeleteResult, err error)
}

// InitDatabase will return a database client for usage
func InitDatabase(URI string) (c *mongo.Client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	if URI == "" {
		URI = "mongodb://mongodb:27017"
	}

	clientOptions := options.Client().ApplyURI(URI)
	clientOptions.Monitor = otelmongo.NewMonitor("mongodb")
	dbClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cancel()
		return &mongo.Client{}, err
	}

	defer cancel()

	return dbClient, err
}

// Health will return if the mongoDB is healthy
func (m MongoCfg) Health() (err error) {
	dbClient, err := InitDatabase(m.URI)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	err = dbClient.Ping(ctx, &readpref.ReadPref{})
	if err != nil {
		cancel()
		return err
	}

	defer cancel()
	return nil
}

// SetIndex will performance a single Indexes().CreateOne operation based on a index name
func (m MongoCfg) SetIndex(ctx context.Context, keys interface{}, opts *options.IndexOptions) (index string, err error) {
	dbClient, err := InitDatabase(m.URI)
	if err != nil {
		return "", err
	}

	col := dbClient.Database(m.Database).Collection(m.Colletion)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	i, err := col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    keys,
			Options: opts,
		},
	)

	defer cancel()

	return i, nil
}

// Get will perform a mongoDB FindOne operation
func (m MongoCfg) Get(ctx context.Context, document primitive.M) (r *mongo.SingleResult, err error) {
	dbClient, err := InitDatabase(m.URI)
	if err != nil {
		return r, err
	}

	col := dbClient.Database(m.Database).Collection(m.Colletion)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	r = col.FindOne(ctx, document)
	defer cancel()

	return r, nil
}

// GetAll will perform a mongoDB Find operation
func (m MongoCfg) GetAll(ctx context.Context, document primitive.M) (r *mongo.Cursor, err error) {
	dbClient, err := InitDatabase(m.URI)
	if err != nil {
		return r, err
	}

	col := dbClient.Database(m.Database).Collection(m.Colletion)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	r, err = col.Find(ctx, document)
	if err != nil {
		cancel()
		return r, err
	}

	defer cancel()

	return r, nil
}

// Create will perform a mongoDB InsertOne operation
func (m MongoCfg) Create(ctx context.Context, document interface{}) (r *mongo.InsertOneResult, err error) {
	dbClient, err := InitDatabase(m.URI)
	if err != nil {
		return r, err
	}

	col := dbClient.Database(m.Database).Collection(m.Colletion)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	r, err = col.InsertOne(ctx, document)
	if err != nil {
		cancel()
		return r, err
	}
	defer cancel()

	return r, nil
}

// Delete will perform a mongoDB DeleteOne operation
func (m MongoCfg) Delete(ctx context.Context, document primitive.M) (r *mongo.DeleteResult, err error) {
	dbClient, err := InitDatabase(m.URI)
	if err != nil {
		return r, err
	}

	col := dbClient.Database(m.Database).Collection(m.Colletion)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	r, err = col.DeleteOne(ctx, document)
	if err != nil {
		cancel()
		return r, err
	}
	defer cancel()

	return r, nil
}
