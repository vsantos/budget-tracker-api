package services

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoCfg satisfies DataManager Interface
type MongoCfg struct {
	URI string
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
