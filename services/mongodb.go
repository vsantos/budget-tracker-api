package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// InitDatabase will return a database client for usage
func InitDatabase() (c *mongo.Client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017/")
	dbClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cancel()
		return &mongo.Client{}, err
	}

	defer cancel()

	return dbClient, err
}

// DatabaseHealth will execute a test connection
func DatabaseHealth() (err error) {
	dbClient, err := InitDatabase()
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
