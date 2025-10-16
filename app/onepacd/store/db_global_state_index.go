package store

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createTimeIndexIndex(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "time_index", Value: 1}},
		Options: options.Index().
			SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err
	}

	return nil
}

func (c *DbClient) createGlobalStateIndex() error {
	c.collection.global_state_index = c.database.Collection("global_state_index")

	if err := createTimeIndexIndex(c.collection.global_state_index); err != nil {
		return err
	}
	return nil
}
