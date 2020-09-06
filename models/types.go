package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	mongodbDatabase          = "budget-tracker"
	mongodbUserCollection    = "users"
	mongodbCardsCollection   = "cards"
	mongodbBalanceCollection = "balance"
	mongodbSpendsCollection  = "spends"
)

// Database creates a Database client
type Database struct {
	client *mongo.Client
}

// User struct defines a user
type User struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Login          string             `json:"login" bson:"login"`
	Firstname      string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname       string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email          string             `json:"email,omitempty" bson:"email,omitempty"`
	SaltedPassword string             `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt      primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// JWTUser defines a user to generate JWT tokens
type JWTUser struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Login    string             `json:"login" bson:"login"`
	Password string             `json:"password" bson:"password"`
}

// SanitizedUser defines a sanited user to GET purposes
type SanitizedUser struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Login     string             `json:"login" bson:"login"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
}

// CreditCard defines a user credit card
type CreditCard struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerID    primitive.ObjectID `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	Alias      string             `json:"alias" bson:"alias"`
	Network    string             `json:"network" bson:"network"`
	LastDigits int32              `json:"last_digits" bson:"last_digits"`
	CreatedAt  primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// Income defines an user outcome for a certain month
type Income struct {
	GrossIncome float64 `json:"gross" bson:"gross"`
	NetIncome   float64 `json:"net" bson:"net"`
}

// Outcome defines an user outcome for a certain month
type Outcome struct {
	FixedOutcome   float64 `json:"fixed" bson:"fixed"`
	DynamicOutcome float64 `json:"dynamic" bson:"dynamic"`
}

// PaymentMethod defines which payment method was used for a certain spend
type PaymentMethod struct {
	Credit      CreditCard `json:"credit,omitempty" bson:"credit,omitempty"`
	Debit       bool       `json:"debit,omitempty" bson:"debit,omitempty"`
	PaymentSlip bool       `json:"payment_slip,omitempty" bson:"payment_slip,omitempty"`
}

// Spend defines a user spend to be added to Balance
type Spend struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerID       primitive.ObjectID `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	Type          string             `json:"type" bson:"type"`
	Description   string             `json:"description" bson:"description"`
	PaymentMethod PaymentMethod      `json:"payment_method,omitempty" bson:"payment_method,omitempty"`
	Categories    []string           `json:"category,omitempty" bson:"category,omitempty"`
	CreatedAt     primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// Balance defines an user balance
type Balance struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerID         primitive.ObjectID `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	Income          Income             `json:"income,omitempty" bson:"income,omitempty"`
	Outcome         Outcome            `json:"outcome" bson:"outcome"`
	SpendableAmount float64            `json:"spendable_amount" bson:"spendable_amount"`
	Historic        []Spend            `json:"historic" bson:"historic"`
	Currency        string             `json:"currency" bson:"currency"`
	Month           int64              `json:"month" bson:"month"`
	Year            int64              `json:"year" bson:"year"`
	CreatedAt       primitive.DateTime `json:"created_at" bson:"created_at"`
}
