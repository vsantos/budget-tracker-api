package models

import (
	"budget-tracker-api/observability"
	"budget-tracker-api/repository"
	"budget-tracker-api/services"
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/attribute"
)

// CreateSpend creates a card for a given owner_id
func CreateSpend(parentCtx context.Context, s repository.Spend) (id string, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("spend.owner.id").String(s.OwnerID.String()),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "CreateSpend", spanTags)
	defer span.End()

	// adding timestamp to creationDate
	t := time.Now()
	s.CreatedAt = primitive.NewDateTimeFromTime(t)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	repo := repository.NewSpendRepository(&repository.SpendRepositoryMongoDB{
		Client: services.MongoClient,
		Config: services.MongoCfg{
			URI:       services.MongodbURI,
			Database:  services.MongodbDatabase,
			Colletion: services.MongodbSpendsCollection,
		},
	})

	id, err = repo.Create(ctx, s)
	fmt.Println(id)
	if err != nil {
		cancel()
		return "", err
	}
	span.SetAttributes(attribute.Key("spend.id").String(id))
	defer cancel()

	observability.Metrics.Spends.SpendsCreated.Inc()
	log.Infoln("created spend", id)
	return id, nil
}

// GetSpends will return all spends from a specific owner_id
func GetSpends(parentCtx context.Context, ownerID string) ([]repository.Spend, error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("spend.owner.id").String(ownerID),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "GetSpends", spanTags)
	defer span.End()

	repo := repository.NewSpendRepository(&repository.SpendRepositoryMongoDB{
		Client: services.MongoClient,
		Config: services.MongoCfg{
			URI:       services.MongodbURI,
			Database:  services.MongodbDatabase,
			Colletion: services.MongodbSpendsCollection,
		},
	})

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	spends, err := repo.Get(ctx, ownerID)
	if err != nil {
		cancel()
		return []repository.Spend{}, err
	}

	return spends, nil
}
