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
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// GetUsers will return all users from database
func GetUsers(parentCtx context.Context) (users []SanitizedUser, err error) {
	ctx, span := observability.Span(parentCtx, "mongodb", "getUsers", []attribute.KeyValue{})
	defer span.End()

	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbUserCollection,
	}

	cursor, err := m.GetAll(ctx, bson.M{})
	if err != nil {
		return []SanitizedUser{}, err
	}
	if err != nil {
		return []SanitizedUser{}, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user SanitizedUser
		cursor.Decode(&user)
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
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

	var user User
	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbUserCollection,
	}

	d, err := m.Get(ctx, bson.M{"_id": pid})
	if err != nil {
		return &User{}, err
	}

	err = d.Decode(&user)
	if err != nil {
		return &User{}, err
	}

	span.SetAttributes(attribute.Key("user.login").String(user.Login))
	return &user, nil
}

// GetUserByFilter will return a user from database based on key,pair BSON
func GetUserByFilter(parentCtx context.Context, bsonKey string, bsonValue string) (u *User, err error) {
	ctx, span := observability.Span(parentCtx, "mongodb", "getUser", []attribute.KeyValue{})
	defer span.End()

	var user User
	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbUserCollection,
	}

	d, err := m.Get(ctx, bson.M{bsonKey: bsonValue})
	if err != nil {
		return &User{}, err
	}

	err = d.Decode(&user)
	if err != nil {
		return &User{}, err
	}

	span.SetAttributes(attribute.Key("user.id").String(user.ID.String()))
	span.SetAttributes(attribute.Key("user.login").String(user.Login))
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

	// adding timestamp to creationDate
	t := time.Now()
	u.CreatedAt = primitive.NewDateTimeFromTime(t)

	// adding salted password for user
	if u.SaltedPassword == "" {
		return "", errors.New("empty password input")
	}

	u.SaltedPassword, err = crypt.GenerateSaltedPassword(u.SaltedPassword)
	if err != nil {
		return "", err
	}

	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbUserCollection,
	}

	_, err = m.SetIndex(
		ctx,
		bsonx.Doc{{Key: "login", Value: bsonx.Int32(1)}},
		options.Index().SetUnique(true),
	)
	if err != nil {
		return "", err
	}

	r, err := m.Create(ctx, u)
	if err != nil {
		return "", err
	}

	observability.Metrics.Users.UsersCreated.Inc()
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

	log.Infoln("deleting user", id)
	var m services.Monger
	m = services.MongoCfg{
		URI:       services.DatabaseURI,
		Database:  services.MongodbDatabase,
		Colletion: services.MongodbUserCollection,
	}

	r, err := m.Delete(ctx, bson.M{"_id": pid})
	if err != nil {
		return err
	}

	log.Infoln("number of users deleted:", r.DeletedCount)

	if r.DeletedCount == 0 {
		return errors.New("non existent user")
	}

	log.Infoln("deleted user", id)
	return nil
}
