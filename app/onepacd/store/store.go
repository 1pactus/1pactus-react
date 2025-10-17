package store

import (
	"github.com/frimin/1pactus-react/config"
	"github.com/frimin/1pactus-react/store/storedriver"
)

func Init(config *config.ConfigBase) error {
	if err := setupMongo(config.Mongo); err != nil {
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
