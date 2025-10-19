package store

import (
	"log"
	"time"

	"github.com/1pactus/1pactus-react/config"
	"github.com/1pactus/1pactus-react/store/storedriver"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MONGO_DB_TIMEOUT = 10 * time.Second
)

type mongoStore struct {
	storedriver.Mongo

	collection *DbCollection
	db         *DbClient
}

func (s *mongoStore) Init(store storedriver.Mongo, conf *config.MongoConfig) {
	s.Mongo = store

	s.db = NewDBClient()
	// Adapting old database connection code
	if err := s.db.Connect(s.GetDatabase()); err != nil {
		log.Fatalf("error to connect db")
	}

	s.collection = s.db.collection
}

func (s *mongoStore) Indexes() map[*mongo.Collection][]mongo.IndexModel {
	return map[*mongo.Collection][]mongo.IndexModel{
		/*s.player: {
			{
				Keys:    bson.D{{Key: "user_id", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},*/
	}
}

func (s *mongoStore) GetDBAdapter() *DbClient {
	return s.db
}
