package store

import (
	"github.com/1pactus/1pactus-react/app/onepacd/store/data"
	"github.com/1pactus/1pactus-react/app/onepacd/store/model"
	"github.com/1pactus/1pactus-react/store/storedriver"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
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

type IKafka interface {
	storedriver.IKafkaStore

	SendBlock(block *pactus.GetBlockResponse) error
	ConsumeBlocks(groupID string, offset int64, handler func(*pactus.GetBlockResponse) (bool, error)) error
	GetLastBlockHeight() (int64, error)
	GetBlockHeightOffset(height int64) (int64, error)
}

var (
	Mongo    IMongo    = &mongoStore{}
	Postgres IPostgres = &postgresStore{}
	Kafka    IKafka    = &kafkaStore{}
)
