package store

import (
	"github.com/1pactus/1pactus-react/config"
	"github.com/1pactus/1pactus-react/store/storedriver"
)

func Init(config *config.ConfigBase) error {
	if err := setupMongo(config.Mongo); err != nil {
		return err
	}

	if err := setupPostgres(config.Postgres); err != nil {
		return err
	}

	return nil
}

func setupMongo(conf *config.MongoConfig) (err error) {
	err = storedriver.MongoStart("base", conf, []storedriver.IMongoStore{
		Mongo,
	})

	return
}

func setupPostgres(conf *config.PostgresConfig) error {
	if err := storedriver.PostgresGormStart("base", conf, []storedriver.IPostgresGormStore{
		Postgres,
	}); err != nil {
		return err
	}

	return nil
}
