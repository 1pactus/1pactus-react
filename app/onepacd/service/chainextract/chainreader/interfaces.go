package chainreader

import (
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type BlockchainReader interface {
	CreateGroup(beginHeight int64, consumerGroupID string) (BlockchainReaderGroup, bool)
	Close()

	GetBlockchainInfo() (*pactus.GetBlockchainInfoResponse, error)
}

type BlockchainReaderGroup interface {
	Read() <-chan *pactus.GetBlockResponse
	Close()
	IsSlowMode() bool
}

const (
	DefaultBlockchainReaderChanSize = 100
)
