package store

import (
	"github.com/frimin/1pactus-react/app/onepacd/store/data"
	"github.com/frimin/1pactus-react/stored/storedriver"
)

type IMongo interface {
	storedriver.IMongoStore

	GetDBAdapter() *DbClient

	FetchNetworkGlobalStats(count int64) ([]data.GlobalStateData, error)
}

var (
	Mongo IMongo = &mongoStore{}
)
