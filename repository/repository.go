package repository

import "context"

// NewDatabaseManagerRepository will return a UserRepository interface based on a struct
func NewDatabaseManagerRepository(d DatabaseManagerRepository) DatabaseManagerRepository {
	return d
}

// NewUserRepository will return a UserRepository interface based on a struct
func NewUserRepository(u UserRepository) UserRepository {
	return u
}

// DatabaseManagerRepository will define instance operations
type DatabaseManagerRepository interface {
	Health() (err error)
	// CreateIndex() (err error)
}

// UserRepository defines a User
type UserRepository interface {
	Get(ctx context.Context, id string) (User, error)
	GetAll(ctx context.Context) ([]SanitizedUser, error)
	Create(ctx context.Context, d User) (id string, err error)
	Delete(ctx context.Context, id string) error
}

// CardRepository defines a User
type CardRepository interface {
	Get(ctx context.Context, id string) (CreditCard, error)
	GetAll(ctx context.Context) ([]CreditCard, error)
	Create(ctx context.Context, c CreditCard) error
	Delete(ctx context.Context, id string) error
}

// SpendRepository defines a User
type SpendRepository interface {
	Get(ctx context.Context, id string) (Spend, error)
	GetAll(ctx context.Context) ([]Spend, error)
	Create(ctx context.Context, s Spend) error
	Delete(ctx context.Context, id string) error
}

// BalanceRepository defines a User
type BalanceRepository interface {
	Get(ctx context.Context, id string) (Balance, error)
	GetAll(ctx context.Context) ([]Balance, error)
	Create(ctx context.Context, b Balance) error
	Delete(ctx context.Context, id string) error
}
