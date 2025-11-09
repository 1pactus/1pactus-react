package store

/*

import (
	"context"
	"fmt"

	"github.com/1pactus/1pactus-react/app/onepacd/store/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *mongoStore) GetNetworkGlobalStats(count int64) ([]data.GlobalStateData, error) {
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

	return rets, nil
}

func (s *mongoStore) GetLastDaysTimeIndex(ctx context.Context, days int) ([]int64, error) {
	if days < 1 {
		return nil, fmt.Errorf("days must be greater than 0")
	}

	var timeIndexes []int64

	opts := options.Find().SetSort(bson.D{{Key: "time_index", Value: -1}}).SetLimit(int64(days))

	cursor, err := s.collection.block.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result struct {
			TimeIndex int64 `bson:"time_index"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		timeIndexes = append(timeIndexes, result.TimeIndex)
	}

	return timeIndexes, nil
}

func (s *mongoStore) GetUnbond(days int) (map[int64]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MONGO_DB_TIMEOUT)

	defer cancel()

	timeIndexes, err := s.GetLastDaysTimeIndex(ctx, days)

	if err != nil {
		return nil, err
	}

	cursor0, err := s.collection.address_unbond_index.Find(ctx, bson.D{
		{Key: "time_index", Value: bson.D{
			{Key: "$in", Value: timeIndexes},
		}},
	})

	if err != nil {
		return nil, err
	}
	defer cursor0.Close(ctx)

	addresses := make([]string, 0)

	if cursor0.Next(ctx) {
		var result struct {
			Address string `bson:"address"`
		}
		if err := cursor0.Decode(&result); err != nil {
			return nil, err
		}
		addresses = append(addresses, result.Address)
	}

	cursor1, err := s.collection.validator_stake.Find(ctx, bson.D{
		{
			Key: "address", Value: bson.D{
				{Key: "$in", Value: addresses},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	defer cursor1.Close(ctx)

	unboundWithTimeIndex := make(map[int64]int64)

	for cursor1.Next(ctx) {
		var result struct {
			Address         string `bson:"address"`
			StakeMax        int64  `bson:"stake_max"`
			UnbondTimeIndex int64  `bson:"unbond_time_index"`
		}

		if err := cursor1.Decode(&result); err != nil {
			return nil, err
		}

		unboundWithTimeIndex[result.UnbondTimeIndex] += result.StakeMax
	}

	return unboundWithTimeIndex, nil
}
*/
