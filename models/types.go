package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User struct defines a user
type User struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Login          string             `json:"login" bson:"login"`
	Firstname      string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname       string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email          string             `json:"email,omitempty" bson:"email,omitempty"`
	SaltedPassword string             `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt      primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// SanitizedUser defines a sanited user to GET purposes
type SanitizedUser struct {
	Login     string `json:"login" bson:"login"`
	Firstname string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email     string `json:"email,omitempty" bson:"email,omitempty"`
}

// CreditCard defines a user credit card
type CreditCard struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OwnerID    primitive.ObjectID `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	Alias      string             `json:"alias" bson:"alias"`
	Network    string             `json:"network" bson:"network"`
	LastDigits int32              `json:"last_digits" bson:"last_digits"`
}
