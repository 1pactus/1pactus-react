package chainreader

import pactus "github.com/pactus-project/pactus/www/grpc/gen/go"

type BlockchainReader interface {
	Read(beginHeight int64, consumerGroupID string) <-chan *pactus.GetBlockResponse
	Close()
	IsSlowMode() bool
}

type BlockchainReaderGroup interface {
	Read() <-chan *pactus.GetBlockResponse
	Close()
	IsSlowMode() bool
}

const (
	DefaultBlockchainReaderChanSize = 100
)
