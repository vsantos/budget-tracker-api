package models

import (
	"budget-tracker/crypt"
	"budget-tracker/services"
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	mongodbDatabase   = "budget-tracker"
	mongodbCollection = "users"
)

// GetUsers will return all users from database
func GetUsers() (us *[]SanitizedUser, err error) {
	dbClient, err := services.InitDatabase()
	if err != nil {
		return &[]SanitizedUser{}, err
	}

	var users []SanitizedUser

	col := dbClient.Database(mongodbDatabase).Collection(mongodbCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		cancel()
		return &[]SanitizedUser{}, err
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
		return &[]SanitizedUser{}, err
	}

	return &users, nil
}

// GetUser will return a user from database
func GetUser(login string) (u *User, err error) {
	dbClient, err := services.InitDatabase()
	if err != nil {
		return &User{}, err
	}

	var user User

	col := dbClient.Database(mongodbDatabase).Collection(mongodbCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	err = col.FindOne(ctx, User{Login: login}).Decode(&user)
	if err != nil {
		cancel()
		return &User{}, err
	}

	defer cancel()
	return &user, nil
}

// CreateUser creates an user based on request body payload
func CreateUser(u *User) (err error) {
	dbClient, err := services.InitDatabase()
	if err != nil {
		return err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

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
		return errors.New("empty password input")
	}

	u.SaltedPassword, err = crypt.GenerateSaltedPassword(u.SaltedPassword)
	if err != nil {
		cancel()
		return err
	}

	_, err = col.InsertOne(ctx, u)
	if err != nil {
		cancel()
		return err
	}

	defer cancel()

	log.Infoln("created user", u.Login)
	return nil
}

// DeleteUser creates an user based on request body payload
func DeleteUser(login string) (err error) {
	dbClient, err := services.InitDatabase()
	if err != nil {
		return err
	}

	col := dbClient.Database(mongodbDatabase).Collection(mongodbCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	log.Infoln("deleting user", login)
	result, err := col.DeleteOne(ctx, bson.M{"login": login})
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

	log.Infoln("deleted user", login)
	return nil
}
