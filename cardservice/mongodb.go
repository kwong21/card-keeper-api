package cardservice

import (
	config "card-keeper-api/internal/configs"
	logger "card-keeper-api/internal/logging"
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
	mongodb := client.Database(configs.Database)

	return &mongoStore{
		db: mongodb,
	}, err
}

func (r *mongoStore) GetAllCardsInCollection(collection string) ([]Card, error) {
	cardsCollection := r.getCollectionFromMongoDB(collection)

	var cards []Card

	cursor, err := cardsCollection.Find(context.TODO(), bson.D{{}})
	err = cursor.All(context.TODO(), &cards)

	return cards, err
}

func (r *mongoStore) AddCardToCollection(card Card, collection string) error {
	cardsCollection := r.getCollectionFromMongoDB(collection)
	var serviceError error = nil

	card.setCardID()

	insert, err := cardsCollection.InsertOne(context.TODO(), card)

	if err != nil {
		serviceError = wrapMongoDBError(err.(mongo.WriteException))
	} else {
		mongoLogger.LogInfoWithFields(
			logger.LogFields{
				"id":   insert.InsertedID,
				"Card": card,
			}, "card inserted to collection")
	}

	return serviceError
}

func (r *mongoStore) getCollectionFromMongoDB(collection string) *mongo.Collection {
	mod := mongo.IndexModel{
		Keys: bson.M{
			"cardID": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cardsCollection := r.db.Collection(collection)

	name, err := cardsCollection.Indexes().CreateOne(ctx, mod)

	if err != nil {
		mongoLogger.LogErrorWithFields(
			logger.LogFields{
				"Error": err,
			}, "not able to create the index for cards")
	}

	if name != "" {
		mongoLogger.LogInfo(fmt.Sprintf("Created index %s", name))
	}

	return cardsCollection
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

func wrapMongoDBError(err mongo.WriteException) error {
	var wrappedError error

	switch mongoErrorCode := err.WriteErrors[0].Code; mongoErrorCode {
	case 11000:
		wrappedError = &DuplicateError{
			Message: err.Error(),
		}
	default:
		wrappedError = &UnknownError{
			Message: err.WriteConcernError.Message,
		}
	}

	return wrappedError
}
