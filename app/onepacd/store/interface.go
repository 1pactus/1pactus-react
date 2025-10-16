package store

import "github.com/frimin/1pactus-react/stored/storedriver"

type IMongo interface {
	storedriver.IMongoStore

	GetDBAdapter() *DbClient
}

var (
	Mongo IMongo = &mongoStore{}
)
