package storedriver

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/frimin/1pactus-react/config"
	"github.com/frimin/1pactus-react/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type IMongoStore interface {
	Init(store Mongo, conf *config.MongoConfig)
	Indexes() map[*mongo.Collection][]mongo.IndexModel
}

type Mongo interface {
	GetDatabase() *mongo.Database
	GetTimeout() time.Duration
}

type mongoImpl struct {
	conf     *config.MongoConfig
	client   *mongo.Client
	database *mongo.Database
	stores   []IMongoStore
	timeout  time.Duration
	log      log.ILogger
}

func (db *mongoImpl) Close() {
	if db.client != nil {
		db.client.Disconnect(context.Background())
	}
}

func MongoStart(name string, conf *config.MongoConfig, stores []IMongoStore) error {
	m := &mongoImpl{
		conf:    conf,
		stores:  stores,
		timeout: time.Second * 10, // Default timeout
		log:     log.WithKv("module", "store").WithKv("mongo", name),
	}

	var err error

	maxRetry := 10

	for {
		if err = m.connect(); err == nil {
			m.initStores()

			m.log.Infof("mongo connect and initialized success")

			if err := m.ensureIndexes(); err != nil {
				return fmt.Errorf("mongo [%s] ensure indexes failed: %v", name, err)
			}
			go m.monitorConnection()
			return nil
		} else {
			maxRetry -= 1

			if maxRetry <= 0 {
				return fmt.Errorf("mongo [%s] connect failed after retries: %v", name, err)
			}

			m.log.Errorf("mongo [%s] connect failed: %v", name, err)

			time.Sleep(time.Second * 5)
		}
	}
}

func (db *mongoImpl) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	journal := true
	cliWC := &writeconcern.WriteConcern{
		W:       "majority",
		Journal: &journal,
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.conf.Uri).
		SetWriteConcern(cliWC).
		SetRetryWrites(true))
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	db.client = client
	db.database = client.Database(db.conf.Database)
	return nil
}

func (db *mongoImpl) initStores() {
	for _, store := range db.stores {
		store.Init(db, db.conf)
	}
}

func (db *mongoImpl) ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	for _, store := range db.stores {
		ensureIndexes := store.Indexes()
		for collection, indexes := range ensureIndexes {
			for _, index := range indexes {
				if err := db.ensureIndex(ctx, collection, index); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (db *mongoImpl) ensureIndex(ctx context.Context, collection *mongo.Collection, indexModel mongo.IndexModel) error {
	indexView := collection.Indexes()

	// List existing indexes
	cursor, err := indexView.List(ctx)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var existingIndexes []bson.M
	if err := cursor.All(ctx, &existingIndexes); err != nil {
		return err
	}
	// Check if the same index already exists
	for _, index := range existingIndexes {
		keys := index["key"].(bson.M)
		targetKeys := indexModel.Keys.(bson.D)

		if reflect.DeepEqual(keys, targetKeys) {
			db.log.Warnf("Index already exists")
			return nil
		}
	}
	// Create new index if not exists
	_, err = indexView.CreateOne(ctx, indexModel)
	return err
}

func (db *mongoImpl) monitorConnection() {
	healthcheck := time.Duration(db.conf.Healthcheck)

	for {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			defer cancel()

			if err := db.client.Ping(ctx, readpref.Primary()); err != nil {
				db.log.Errorf("lost mongodb connection, retrying: %v", err)

				for {
					if err := db.connect(); err == nil {
						db.initStores()
						db.log.Info("successfully reconnected to mongodb")
						break
					}
					db.log.Errorf("error connecting to mongodb: %v", err)
					time.Sleep(5 * time.Second)
				}
			}

		}()

		time.Sleep(healthcheck * time.Second)
	}
}

func (db *mongoImpl) GetDatabase() *mongo.Database {
	return db.database
}

func (db *mongoImpl) GetTimeout() time.Duration {
	return db.timeout
}
