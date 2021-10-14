package models

import (
	"budget-tracker-api/crypt"
	"budget-tracker-api/observability"
	"budget-tracker-api/repository"
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
func GetUsers(parentCtx context.Context) ([]repository.SanitizedUser, error) {
	ctx, span := observability.Span(parentCtx, "mongodb", "getUsers", []attribute.KeyValue{})
	defer span.End()

	repo := repository.NewUserRepository(&repository.UserRepositoryMongoDB{
		Client: services.MongoClient,
		Config: services.MongoCfg{
			URI:       services.MongodbURI,
			Database:  services.MongodbDatabase,
			Colletion: services.MongodbUserCollection,
		},
	})

	users, err := repo.GetAll(ctx)
	if err != nil {
		return []repository.SanitizedUser{}, err
	}

	return users, nil
}

// GetUser will return a user from database based on ID
func GetUser(parentCtx context.Context, id string) (*repository.User, error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("user.id").String(id),
	}
	ctx, span := observability.Span(parentCtx, "mongodb", "getUser", spanTags)
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	repo := repository.NewUserRepository(&repository.UserRepositoryMongoDB{
		Client: services.MongoClient,
		Config: services.MongoCfg{
			URI:       services.MongodbURI,
			Database:  services.MongodbDatabase,
			Colletion: services.MongodbUserCollection,
		},
	})

	u, err := repo.Get(ctx, id)
	if err != nil {
		cancel()
		return &repository.User{}, err
	}

	span.SetAttributes(attribute.Key("user.login").String(u.Login))
	defer cancel()
	return &u, nil
}

// GetUserByFilter will return a user from database based on key,pair BSON
func GetUserByFilter(parentCtx context.Context, bsonKey string, bsonValue string) (u *repository.User, err error) {
	ctx, span := observability.Span(parentCtx, "mongodb", "getUser", []attribute.KeyValue{})
	defer span.End()

	var user repository.User

	col := services.MongoClient.Database(mongodbDatabase).Collection(mongodbUserCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	err = col.FindOne(ctx, bson.M{bsonKey: bsonValue}).Decode(&user)
	if err != nil {
		cancel()
		return &repository.User{}, err
	}

	span.SetAttributes(attribute.Key("user.id").String(user.ID.String()))
	span.SetAttributes(attribute.Key("user.login").String(user.Login))
	defer cancel()
	return &user, nil
}

// CreateUser creates an user based on request body payload
func CreateUser(parentCtx context.Context, u repository.User) (id string, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("user.id").String(u.ID.String()),
		attribute.Key("user.login").String(u.Login),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "CreateUser", spanTags)
	defer span.End()

	col := services.MongoClient.Database(mongodbDatabase).Collection(mongodbUserCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	_, err = col.Indexes().CreateOne(
		ctx,
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

	repo := repository.NewUserRepository(&repository.UserRepositoryMongoDB{
		Client: services.MongoClient,
		Config: services.MongoCfg{
			URI:       services.MongodbURI,
			Database:  services.MongodbDatabase,
			Colletion: services.MongodbUserCollection,
		},
	})

	id, err = repo.Create(ctx, u)
	if err != nil {
		cancel()
		return "", err
	}

	defer cancel()

	observability.Metrics.Users.UsersCreated.Inc()
	log.Infoln("created user", u.Login)
	return id, nil
}

// DeleteUser creates an user based on request body payload
func DeleteUser(parentCtx context.Context, id string) (err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("user.id").String(id),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "DeleteUser", spanTags)
	defer span.End()

	repo := repository.NewUserRepository(&repository.UserRepositoryMongoDB{
		Client: services.MongoClient,
		Config: services.MongoCfg{
			URI:       services.MongodbURI,
			Database:  services.MongodbDatabase,
			Colletion: services.MongodbUserCollection,
		},
	})

	err = repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	log.Infoln("deleted user", id)
	return nil
}
