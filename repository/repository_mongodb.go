package repository

import (
	"budget-tracker-api/services"
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

// DatabaseRepositoryMongoDB defines a struct for mongoDB database operations
type DatabaseRepositoryMongoDB struct {
	Client *mongo.Client
	Config services.MongoCfg
}

// UserRepositoryMongoDB defines a struct for mongoDB User operations
type UserRepositoryMongoDB struct {
	Client *mongo.Client
	Config services.MongoCfg
}

// CardRepositoryMongoDB defines a struct for mongoDB Card operations
type CardRepositoryMongoDB struct {
	Client *mongo.Client
	Config services.MongoCfg
}

// BalanceRepositoryMongoDB defines a struct for mongoDB Balance operations
type BalanceRepositoryMongoDB struct {
	Client *mongo.Client
	Config services.MongoCfg
}

// SpendRepositoryMongoDB defines a struct for mongoDB Spend operations
type SpendRepositoryMongoDB struct {
	Client *mongo.Client
	Config services.MongoCfg
}

// Health will define a mongoDB healthcheck
func (d *DatabaseRepositoryMongoDB) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	err := d.Client.Ping(ctx, &readpref.ReadPref{})
	if err != nil {
		cancel()
		return err
	}

	defer cancel()
	return nil
}

// Get will
func (u *UserRepositoryMongoDB) Get(ctx context.Context, id string) (User, error) {
	var user User

	// col := u.Client.Database(u.Config.Database).Collection(u.Config.Colletion)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		cancel()
		return User{}, err
	}

	r, err := u.Config.Get(ctx, bson.M{"_id": pid})
	if err != nil {
		cancel()
		if strings.Contains(err.Error(), "no documents in result") {
			return User{}, errors.New("could not find user")
		}

		return User{}, err
	}
	r.Decode(&user)

	defer cancel()

	return user, nil
}

// GetAll will
func (u *UserRepositoryMongoDB) GetAll(ctx context.Context) ([]SanitizedUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	cursor, err := u.Config.GetAll(ctx, bson.M{})
	if err != nil {
		cancel()
		return []SanitizedUser{}, err
	}

	defer cursor.Close(ctx)
	defer cancel()

	var users []SanitizedUser
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

// Create will
func (u *UserRepositoryMongoDB) Create(ctx context.Context, user User) (id string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	r, err := u.Config.Create(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			cancel()
			return "", errors.New("user already exists")
		}

		cancel()
		return "", err
	}

	return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Delete will
func (u *UserRepositoryMongoDB) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		cancel()
		return err
	}

	r, err := u.Config.Delete(ctx, bson.M{"_id": pid})
	if err != nil {
		cancel()
		return err
	}

	if r.DeletedCount == 0 {
		cancel()
		return errors.New("non existent user")
	}

	return nil

}

// Get will
func (r *BalanceRepositoryMongoDB) Get(ctx context.Context, id string) (CreditCard, error) {
	return CreditCard{}, nil
}

// GetAll will
func (r *BalanceRepositoryMongoDB) GetAll(ctx context.Context) ([]CreditCard, error) {
	return []CreditCard{}, nil
}

// Create will
func (r *BalanceRepositoryMongoDB) Create(ctx context.Context, card CreditCard) error {
	return nil
}

// Delete will
func (r *BalanceRepositoryMongoDB) Delete(ctx context.Context, id string) error {
	return nil
}

// Get will
func (r *SpendRepositoryMongoDB) Get(ctx context.Context, id string) (CreditCard, error) {
	return CreditCard{}, nil
}

// GetAll will
func (r *SpendRepositoryMongoDB) GetAll(ctx context.Context) ([]CreditCard, error) {
	return []CreditCard{}, nil
}

// Create will
func (r *SpendRepositoryMongoDB) Create(ctx context.Context, card CreditCard) error {
	return nil
}

// Delete will
func (r *SpendRepositoryMongoDB) Delete(ctx context.Context, id string) error {
	return nil
}
