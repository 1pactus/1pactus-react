package store

import (
	"github.com/frimin/1pactus-react/app/onepacd/store/data"
	"github.com/frimin/1pactus-react/app/onepacd/store/model"
	"github.com/frimin/1pactus-react/store/storedriver"
)

type IMongo interface {
	storedriver.IMongoStore

	GetDBAdapter() *DbClient

	GetNetworkGlobalStats(count int64) ([]data.GlobalStateData, error)
	GetUnbond(days int) (map[int64]int64, error)
}

type IPostgres interface {
	storedriver.IPostgresGormStore

	GetTopGlobalState() (*model.GlobalState, error)
	InsertGlobalState(state *model.GlobalState) error
	GetNetworkGlobalStats(count int64) ([]model.GlobalState, error)

	GetTopBlock() (*model.Block, error)
	Commit(commitContext PgCommitContext) error
}

var (
	Mongo    IMongo    = &mongoStore{}
	Postgres IPostgres = &postgresStore{}
)
