package models

import (
	"budget-tracker-api/observability"
	"budget-tracker-api/repository"
	"budget-tracker-api/services"
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.opentelemetry.io/otel/attribute"
)

// CreateCard creates a card for a given owner_id
func CreateCard(parentCtx context.Context, c repository.CreditCard) (id string, err error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("card.owner.id").String(c.OwnerID.String()),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "CreateCard", spanTags)
	defer span.End()

	col := services.MongoClient.Database(mongodbDatabase).Collection(mongodbCardsCollection)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	_, err = col.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bsonx.Doc{{Key: "owner_id", Value: bsonx.Int32(1)}, {Key: "last_digits", Value: bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)

	// adding timestamp to creationDate
	t := time.Now()
	c.CreatedAt = primitive.NewDateTimeFromTime(t)

	repo := repository.NewCardRepository(&repository.CardRepositoryMongoDB{
		Client: services.MongoClient,
		Config: services.MongoCfg{
			URI:       services.MongodbURI,
			Database:  services.MongodbDatabase,
			Colletion: services.MongodbCardsCollection,
		},
	})

	id, err = repo.Create(ctx, c)
	if err != nil {
		cancel()
		return "", err
	}

	span.SetAttributes(attribute.Key("card.id").String(id))
	defer cancel()

	observability.Metrics.Cards.CardsCreated.Inc()
	log.Infoln("created card", c.Alias)
	return id, nil
}

// GetAllCards will return a list of all cards from the database
func GetAllCards(parentCtx context.Context) ([]repository.CreditCard, error) {
	ctx, span := observability.Span(parentCtx, "mongodb", "GetAllCards", []attribute.KeyValue{})
	defer span.End()

	repo := repository.NewCardRepository(&repository.CardRepositoryMongoDB{
		Client: services.MongoClient,
		Config: services.MongoCfg{
			URI:       services.MongodbURI,
			Database:  services.MongodbDatabase,
			Colletion: services.MongodbCardsCollection,
		},
	})

	cards, err := repo.GetAll(ctx)
	if err != nil {
		return []repository.CreditCard{}, err
	}

	return cards, nil
}

// GetCards will return a list of cards from a owner_id
func GetCards(parentCtx context.Context, ownerID string) ([]repository.CreditCard, error) {
	spanTags := []attribute.KeyValue{
		attribute.Key("card.owner.id").String(ownerID),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "GetUserCards", spanTags)
	defer span.End()

	repo := repository.NewCardRepository(&repository.CardRepositoryMongoDB{
		Client: services.MongoClient,
		Config: services.MongoCfg{
			URI:       services.MongodbURI,
			Database:  services.MongodbDatabase,
			Colletion: services.MongodbCardsCollection,
		},
	})

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	cards, err := repo.Get(ctx, ownerID)
	if err != nil {
		cancel()
		return []repository.CreditCard{}, err
	}

	defer cancel()

	return cards, nil
}

// DeleteCard creates an user based on request body payload
func DeleteCard(parentCtx context.Context, id string) error {
	spanTags := []attribute.KeyValue{
		attribute.Key("card.id").String(id),
	}

	ctx, span := observability.Span(parentCtx, "mongodb", "DeleteUserCard", spanTags)
	defer span.End()

	repo := repository.NewCardRepository(&repository.CardRepositoryMongoDB{
		Client: services.MongoClient,
		Config: services.MongoCfg{
			URI:       services.MongodbURI,
			Database:  services.MongodbDatabase,
			Colletion: services.MongodbCardsCollection,
		},
	})

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	log.Infoln("deleting card", id)

	err := repo.Delete(ctx, id)
	if err != nil {
		cancel()
		return err
	}

	defer cancel()

	log.Infoln("deleted card", id)
	return nil
}
