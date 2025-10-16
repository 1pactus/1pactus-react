package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlockData struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Height    uint32              `json:"height,omitempty" bson:"height,omitempty"`
	TimeIndex uint32              `json:"time_index,omitempty" bson:"time_index,omitempty"`
}

func CreateBlockDataIndex(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "height", Value: 1}},
		Options: options.Index().
			SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err
	}

	return nil
}
