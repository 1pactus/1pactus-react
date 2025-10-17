package store

import (
	"context"
	"slices"

	"github.com/frimin/1pactus-react/app/onepacd/store/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *mongoStore) FetchNetworkGlobalStats(count int64) ([]data.GlobalStateData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MONGO_DB_TIMEOUT)
	defer cancel()

	var rets []data.GlobalStateData

	opts := options.Find().
		SetSort(bson.D{{Key: "time_index", Value: -1}})

	if count > 0 {
		opts = opts.SetLimit(count)
	}

	cursor, err := s.collection.global_state_index.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &rets); err != nil {
		return nil, err
	}

	slices.Reverse(rets)

	return rets, nil
}
