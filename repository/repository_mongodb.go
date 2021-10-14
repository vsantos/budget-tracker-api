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

// GetAll will return all Users (in a sanitized way)
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

// Create will create a user based
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

// Delete will delete a user based on it's ID
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

// Get will return all cards from a given owner ID
func (c *CardRepositoryMongoDB) Get(ctx context.Context, ownerID string) ([]CreditCard, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		cancel()
		return []CreditCard{}, err
	}

	cursor, err := c.Config.GetAll(ctx, bson.M{"owner_id": oid})
	if err != nil {
		cancel()
		return []CreditCard{}, err
	}

	var cards []CreditCard
	for cursor.Next(ctx) {
		var card CreditCard
		cursor.Decode(&card)
		cards = append(cards, card)
		defer cancel()
	}

	if err := cursor.Err(); err != nil {
		cancel()
		return []CreditCard{}, err
	}

	if len(cards) == 0 {
		return []CreditCard{}, errors.New("could not find any cards")
	}

	return cards, nil
}

// GetAll will return literally all cards from the database
func (c *CardRepositoryMongoDB) GetAll(ctx context.Context) ([]CreditCard, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := c.Config.GetAll(ctx, bson.M{})
	if err != nil {
		cancel()
		return []CreditCard{}, err
	}

	var cards []CreditCard
	for cursor.Next(ctx) {
		var card CreditCard
		cursor.Decode(&card)
		cards = append(cards, card)
		defer cancel()
	}

	if err := cursor.Err(); err != nil {
		cancel()
		return []CreditCard{}, err
	}

	return cards, nil
}

// Create will create a card
func (c *CardRepositoryMongoDB) Create(ctx context.Context, card CreditCard) (id string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	r, err := c.Config.Create(ctx, card)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			cancel()
			return "", errors.New("card already exists")
		}

		cancel()
		return "", err
	}

	return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Delete will delete a card based on it's ID
func (c *CardRepositoryMongoDB) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		cancel()
		return err
	}

	r, err := c.Config.Delete(ctx, bson.M{"_id": pid})
	if err != nil {
		cancel()
		return err
	}

	if r.DeletedCount == 0 {
		cancel()
		return errors.New("non existent card")
	}

	return nil
}

// Get will
func (b *BalanceRepositoryMongoDB) Get(ctx context.Context, ownerID string, month int64, year int64) (Balance, error) {
	var balance Balance

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		cancel()
		return Balance{}, err
	}

	r, err := b.Config.Get(ctx, bson.M{
		"owner_id": oid,
		"month":    month,
		"year":     year,
	})

	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			cancel()
			return Balance{}, errors.New("could not find balance")
		}
		cancel()
		return Balance{}, err
	}

	r.Decode(&balance)

	return balance, nil
}

// GetAll will
func (b *BalanceRepositoryMongoDB) GetAll(ctx context.Context) ([]Balance, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := b.Config.GetAll(ctx, bson.M{})
	if err != nil {
		cancel()
		return []Balance{}, err
	}

	var balances []Balance
	for cursor.Next(ctx) {
		var balance Balance
		cursor.Decode(&balance)
		balances = append(balances, balance)
		defer cancel()
	}

	if err := cursor.Err(); err != nil {
		cancel()
		return []Balance{}, err
	}

	if len(balances) == 0 {
		return []Balance{}, errors.New("could not find any balances")
	}

	return balances, nil
}

// Create will
func (b *BalanceRepositoryMongoDB) Create(ctx context.Context, balance Balance) (id string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	r, err := b.Config.Create(ctx, balance)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			cancel()
			return "", errors.New("balance already exists")
		}

		cancel()
		return "", err
	}

	return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Delete will
func (b *BalanceRepositoryMongoDB) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		cancel()
		return err
	}

	r, err := b.Config.Delete(ctx, bson.M{"_id": pid})
	if err != nil {
		cancel()
		return err
	}

	if r.DeletedCount == 0 {
		cancel()
		return errors.New("non existent balance")
	}

	return nil
}

// Get will return a list of spends from a given ownerID
func (s *SpendRepositoryMongoDB) Get(ctx context.Context, ownerID string) ([]Spend, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		cancel()
		return []Spend{}, err
	}

	cursor, err := s.Config.GetAll(ctx, bson.M{"owner_id": oid})
	if err != nil {
		cancel()
		return []Spend{}, err
	}

	var spends []Spend
	for cursor.Next(ctx) {
		var spend Spend
		cursor.Decode(&spend)
		spends = append(spends, spend)
		defer cancel()
	}

	if err := cursor.Err(); err != nil {
		cancel()
		return []Spend{}, err
	}

	if len(spends) == 0 {
		return []Spend{}, errors.New("could not find any spends")
	}

	return spends, nil
}

// GetAll will all spends in database but currently is not supported
func (s *SpendRepositoryMongoDB) GetAll(ctx context.Context) ([]Spend, error) {
	return []Spend{}, nil
}

// Create will
func (s *SpendRepositoryMongoDB) Create(ctx context.Context, spend Spend) (id string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	r, err := s.Config.Create(ctx, spend)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			cancel()
			return "", errors.New("balance already exists")
		}

		cancel()
		return "", err
	}

	return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Delete will
func (s *SpendRepositoryMongoDB) Delete(ctx context.Context, id string) error {
	return nil
}
