package models

import (
	"budget-tracker-api/observability"
	"budget-tracker-api/repository"
	"budget-tracker-api/services"
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/attribute"
)

// CreateBalance creates a balance for a given owner_id
func CreateBalance(parentCtx context.Context, b repository.Balance) (id string, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("balance.owner.id").String(b.OwnerID.String()),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "CreateBalance", spanTags)
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	// adding timestamp to creationDate
	t := time.Now()
	b.CreatedAt = primitive.NewDateTimeFromTime(t)
	b.UpdatedAt = primitive.NewDateTimeFromTime(t)
	b.SpendableAmount = b.Income.NetIncome
	b.Historic = []repository.Spend{}

	repo := repository.NewBalanceRepository(&repository.BalanceRepositoryMongoDB{
		Client: services.MongoClient,
		Config: services.MongoCfg{
			URI:       services.MongodbURI,
			Database:  services.MongodbDatabase,
			Colletion: services.MongodbBalanceCollection,
		},
	})

	id, err = repo.Create(ctx, b)
	if err != nil {
		cancel()
		return "", err
	}

	span.SetAttributes(attribute.Key("balance.id").String(id))
	defer cancel()

	observability.Metrics.Balances.BalancesCreated.Inc()
	log.Infoln("created balance", id)
	return id, nil
}

// GetBalance will return a balance from an owner_id based on month and year
func GetBalance(parentCtx context.Context, ownerID string, month int64, year int64) (*repository.Balance, error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("balance.owner.id").String(ownerID),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "GetBalance", spanTags)
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	repo := repository.NewBalanceRepository(&repository.BalanceRepositoryMongoDB{
		Client: services.MongoClient,
		Config: services.MongoCfg{
			URI:       services.MongodbURI,
			Database:  services.MongodbDatabase,
			Colletion: services.MongodbBalanceCollection,
		},
	})

	b, err := repo.Get(ctx, ownerID, month, year)
	if err != nil {
		cancel()
		return &repository.Balance{}, err
	}

	defer cancel()
	return &b, nil
}

// GetAllBalances will return all balances from an owner_id
func GetAllBalances(parentCtx context.Context, ownerID string) ([]repository.Balance, error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("balance.owner.id").String(ownerID),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "GetAllBalances", spanTags)
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	repo := repository.NewBalanceRepository(&repository.BalanceRepositoryMongoDB{
		Client: services.MongoClient,
		Config: services.MongoCfg{
			URI:       services.MongodbURI,
			Database:  services.MongodbDatabase,
			Colletion: services.MongodbBalanceCollection,
		},
	})

	b, err := repo.GetAll(ctx)
	if err != nil {
		cancel()
		return []repository.Balance{}, err
	}

	defer cancel()

	return b, nil
}
