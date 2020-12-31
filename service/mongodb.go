package service

import (
	"card-keeper-api/config"
	logger "card-keeper-api/log"
	"card-keeper-api/model"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoStore struct {
	db *mongo.Database
}

var mongoLogger = logger.NewLogger()

// MongoDB returns the MongoDB service.
func MongoDB(configs config.DBConfiguration) (Repository, error) {
	uri := buildConnectionURI(configs)

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		mongoLogger.LogErrorWithFields(
			logger.Fields{
				"Error": err,
			}, "Could not connect to mongodb!")
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		mongoLogger.LogErrorWithFields(
			logger.Fields{
				"Error": err,
			}, "Could not connect to mongodb!")
	}

	mongodb := client.Database(configs.Database)

	createIndexForCardsCollection(mongodb)

	return &mongoStore{
		db: mongodb,
	}, err
}

func (r *mongoStore) GetAll() (*[]model.Card, error) {
	cardsCollection := r.db.Collection("cards")

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
	cardsCollection := r.db.Collection("cards")
	var serviceError error = nil

	insert, err := cardsCollection.InsertOne(context.TODO(), card)

	if err != nil {
		mongoLogger.LogErrorWithFields(
			logger.Fields{
				"Error": err,
				"Card":  card,
			}, "not able to add card to collection")

		serviceError = wrapMongoDBError(err.(mongo.WriteException))
	} else {
		mongoLogger.LogInfoWithFields(
			logger.Fields{
				"id":   insert.InsertedID,
				"Card": card,
			}, "card inserted to collection")
	}

	return serviceError
}

func buildConnectionURI(dbConfig config.DBConfiguration) string {
	var auth string
	var replicaSet string

	if dbConfig.User != "" && dbConfig.Password != "" {
		auth = fmt.Sprintf("%s:%s@", dbConfig.User, dbConfig.Password)
	} else {
		mongoLogger.LogInfo("auth-less mongodb connection")
	}

	if dbConfig.ReplicaSet != "" {
		replicaSet = fmt.Sprintf("?replicaSet=%s", dbConfig.ReplicaSet)
	} else {
		mongoLogger.LogInfo("no replica set configured")
	}

	hosts := dbConfig.Host
	uri := fmt.Sprintf("mongodb://%s%s/%s", auth, hosts, replicaSet)

	return uri
}

func createIndexForCardsCollection(client *mongo.Database) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"base.year": 1, "base.set": 1, "base.make": 1, "base.player": 1},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cardsCollection := client.Collection("cards")

	_, err := cardsCollection.Indexes().CreateOne(ctx, mod)

	if err != nil {
		mongoLogger.LogErrorWithFields(
			logger.Fields{
				"Error": err,
			}, "not able to create the index for cards")
	}
}

func wrapMongoDBError(err mongo.WriteException) error {
	var wrappedError error

	switch mongoErrorCode := err.WriteErrors[0].Code; mongoErrorCode {
	case 11000:
		wrappedError = &DuplicateError{}
	default:
		wrappedError = &UnknownError{
			Message: err.WriteConcernError.Message,
		}
	}

	return wrappedError
}
