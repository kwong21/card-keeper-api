package service

import (
	logger "card-keeper-api/log"
	"card-keeper-api/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoStore struct {
	db *mongo.Client
}

var mongoLogger = logger.NewLogger()

// MongoDB returns the MongoDB service.
func MongoDB() (Repository, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		mongoLogger.LogErrorWithFields(
			logger.Fields{
				"Error": err,
			}, "Could not connect to mongodb!")
	}
	createIndexForCardsCollection(client)
	return &mongoStore{
		db: client,
	}, err
}

func createIndexForCardsCollection(client *mongo.Client) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"base.year": 1, "base.set": 1, "base.make": 1, "base.player": 1},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cardsCollection := client.Database("card-keeper").Collection("cards")

	_, err := cardsCollection.Indexes().CreateOne(ctx, mod)

	if err != nil {
		mongoLogger.LogErrorWithFields(
			logger.Fields{
				"Error": err,
			}, "not able to create the index for cards")
	}

}

func (r *mongoStore) GetAll() (*[]model.Card, error) {
	cardsCollection := r.db.Database("card-keeper").Collection("cards")

	var cards []model.Card

	cursor, err := cardsCollection.Find(context.TODO(), bson.D{{}})

	if err != nil {
		mongoLogger.LogErrorWithFields(
			logger.Fields{
				"Error": err,
			}, "not able to retrieve cards from collection")
	}

	if err = cursor.All(context.TODO(), &cards); err != nil {
		if err != nil {
			mongoLogger.LogErrorWithFields(
				logger.Fields{
					"Error": err,
				}, "not able to retrieve cards from collection")
		}
	}

	return &cards, err
}

func (r *mongoStore) AddCard(card model.Card) error {
	cardsCollection := r.db.Database("card-keeper").Collection("cards")

	opts := options.Update().SetUpsert(true)

	insert, err := cardsCollection.UpdateOne(context.TODO(), nil, nil, opts)

	if err != nil {
		mongoLogger.LogErrorWithFields(
			logger.Fields{
				"Error": err,
				"Card":  card,
			}, "not able to add card to collection")
	} else {
		mongoLogger.LogInfoWithFields(
			logger.Fields{
				"id":   insert.UpsertedID,
				"Card": card,
			}, "card inserted to collection")
	}

	return err
}
