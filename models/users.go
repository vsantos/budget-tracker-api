package models

import (
	"budget-tracker-api/crypt"
	"budget-tracker-api/observability"
	"budget-tracker-api/services"
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// GetUsers will return all users from database
func GetUsers(parentCtx context.Context) (users []SanitizedUser, err error) {
	ctx, span := observability.Span(parentCtx, "mongodb", "getUsers", []attribute.KeyValue{})
	defer span.End()

	dbClient, err := services.InitDatabase()
	if err != nil {
		return []SanitizedUser{}, err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbUserCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		cancel()
		return []SanitizedUser{}, err
	}

	defer cursor.Close(ctx)
	defer cancel()

	for cursor.Next(ctx) {
		var user SanitizedUser
		cursor.Decode(&user)
		users = append(users, user)
		defer cancel()
	}

	if err := cursor.Err(); err != nil {
		cancel()
		return []SanitizedUser{}, err
	}

	return users, nil
}

// GetUser will return a user from database based on ID
func GetUser(parentCtx context.Context, id string) (u *User, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("user.id").String(id),
	}
	ctx, span := observability.Span(parentCtx, "mongodb", "getUser", spanTags)
	defer span.End()

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &User{}, err
	}

	dbClient, err := services.InitDatabase()
	if err != nil {
		return &User{}, err
	}

	var user User

	col := dbClient.Database(mongodbDatabase).Collection(mongodbUserCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	err = col.FindOne(ctx, bson.M{"_id": pid}).Decode(&user)
	if err != nil {
		cancel()
		return &User{}, err
	}

	span.SetAttributes(attribute.Key("user.login").String(user.Login))
	defer cancel()
	return &user, nil
}

// GetUserByFilter will return a user from database based on key,pair BSON
func GetUserByFilter(bsonKey string, bsonValue string) (u *User, err error) {
	dbClient, err := services.InitDatabase()
	if err != nil {
		return &User{}, err
	}

	var user User

	col := dbClient.Database(mongodbDatabase).Collection(mongodbUserCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	err = col.FindOne(ctx, bson.M{bsonKey: bsonValue}).Decode(&user)
	if err != nil {
		cancel()
		return &User{}, err
	}

	defer cancel()
	return &user, nil
}

// CreateUser creates an user based on request body payload
func CreateUser(parentCtx context.Context, u User) (id string, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("user.id").String(u.ID.String()),
		attribute.Key("user.login").String(u.Login),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "CreateUser", spanTags)
	defer span.End()

	dbClient, err := services.InitDatabase()
	if err != nil {
		return "", err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbUserCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	_, err = col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{Key: "login", Value: bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)

	// adding timestamp to creationDate
	t := time.Now()
	u.CreatedAt = primitive.NewDateTimeFromTime(t)

	// adding salted password for user
	if u.SaltedPassword == "" {
		cancel()
		return "", errors.New("empty password input")
	}

	u.SaltedPassword, err = crypt.GenerateSaltedPassword(u.SaltedPassword)
	if err != nil {
		cancel()
		return "", err
	}

	r, err := col.InsertOne(ctx, u)
	if err != nil {
		cancel()
		return "", err
	}

	defer cancel()

	log.Infoln("created user", u.Login)
	return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

// DeleteUser creates an user based on request body payload
func DeleteUser(parentCtx context.Context, id string) (err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("user.id").String(id),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "DeleteUser", spanTags)
	defer span.End()

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	dbClient, err := services.InitDatabase()
	if err != nil {
		return err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbUserCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	log.Infoln("deleting user", id)
	result, err := col.DeleteOne(ctx, bson.M{"_id": pid})
	if err != nil {
		cancel()
		return err
	}

	log.Infoln("number of users deleted:", result.DeletedCount)

	if result.DeletedCount == 0 {
		cancel()
		return errors.New("non existent user")
	}

	defer cancel()

	log.Infoln("deleted user", id)
	return nil
}
