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

// NewCardRepository will return a CardRepository interface based on a struct
func NewCardRepository(c CardRepository) CardRepository {
	return c
}

// NewBalanceRepository will return a BalanceRepository interface based on a struct
func NewBalanceRepository(b BalanceRepository) BalanceRepository {
	return b
}

// NewSpendRepository will return a SpendRepository interface based on a struct
func NewSpendRepository(s SpendRepository) SpendRepository {
	return s
}

// DatabaseManagerRepository will define instance operations
type DatabaseManagerRepository interface {
	Health() (err error)
	// CreateIndex() (err error)
}

// UserRepository defines a User
type UserRepository interface {
	Get(ctx context.Context, id string) (SanitizedUser, error)
	GetAll(ctx context.Context) ([]SanitizedUser, error)
	Create(ctx context.Context, d User) (id string, err error)
	Delete(ctx context.Context, id string) error
}

// CardRepository defines a Card
type CardRepository interface {
	Get(ctx context.Context, ownerID string) ([]CreditCard, error)
	GetAll(ctx context.Context) ([]CreditCard, error)
	Create(ctx context.Context, c CreditCard) (id string, err error)
	Delete(ctx context.Context, id string) error
}

// SpendRepository defines a Spend
type SpendRepository interface {
	Get(ctx context.Context, ownerID string) ([]Spend, error)
	GetAll(ctx context.Context) ([]Spend, error)
	Create(ctx context.Context, s Spend) (id string, err error)
	Delete(ctx context.Context, id string) error
}

// BalanceRepository defines a Balance
type BalanceRepository interface {
	Get(ctx context.Context, ownerID string, month int64, year int64) (Balance, error)
	GetAll(ctx context.Context) ([]Balance, error)
	Create(ctx context.Context, b Balance) (id string, err error)
	Delete(ctx context.Context, id string) error
}
