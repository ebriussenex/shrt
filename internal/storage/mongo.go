package shorten_store

import (
	"context"
	"fmt"
	"time"

	"github.com/ebriussenex/shrt/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoStore struct {
	db *mongo.Database
}

func NewMongoDB(client *mongo.Client) *mongoStore {
	return &mongoStore{db: client.Database("shrt")}
}

func (m *mongoStore) collection() *mongo.Collection {
	return m.db.Collection("shrt_urls")
}

func (m *mongoStore) Put(ctx context.Context, shortened model.Shortened) (*model.Shortened, error) {
	const op = "Putting shortened to mongoDb"

	shortened.CreatedAt = time.Now().UTC()

	count, err := m.collection().CountDocuments(ctx, bson.M{"_id": shortened.Id})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if count > 0 {
		return nil, fmt.Errorf("%s: %w", op, model.ErrIdAlreadyExists)
	}

	_, err = m.collection().InsertOne(ctx, model.EntToMongo(shortened))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &shortened, nil
}

func (m *mongoStore) Get(ctx context.Context, id string) (*model.Shortened, error) {
	const op = "Getting shortened from mongo"
	var shorten model.MongoShortenUrl

	if err := m.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&shorten); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%s: %w", op, model.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return model.MongoToEnt(shorten), nil
}

func (m *mongoStore) IncrVisits(ctx context.Context, id string) error {
	const op = "Updating shortened url visit count"
	
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{"visitCount": 1},
		"$set": bson.M{"updatedAt": time.Now().UTC()},
	}

	_, err := m.collection().UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
