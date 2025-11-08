package chainreader

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/1pactus/1pactus-react/app/onepacd/service/gather"
	"github.com/1pactus/1pactus-react/log"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type blockchainGrpcReaderImpl struct {
	height          int64
	lastBlockHeight int64
	grpc            *gather.GrpcClient
	log             log.ILogger
	blockChan       chan *pactus.GetBlockResponse
	ctx             context.Context
	cancel          context.CancelFunc
	lastError       error
	runOnce         sync.Once
	closeOnce       sync.Once
	slowMode        bool
}

func NewBlockchainGrpcReader(parentCtx context.Context, grpc *gather.GrpcClient, parentLogger log.ILogger) (BlockchainReader, error) {
	reader := &blockchainGrpcReaderImpl{
		grpc:      grpc,
		log:       parentLogger.WithKv("reader", "grpc"),
		blockChan: make(chan *pactus.GetBlockResponse, DefaultBlockchainReaderChanSize),
		slowMode:  false,
	}

	reader.ctx, reader.cancel = context.WithCancel(parentCtx)

	return reader, nil
}

func (r *blockchainGrpcReaderImpl) Read(beginHeight int64, consumerGroupID string) <-chan *pactus.GetBlockResponse {
	r.height = beginHeight

	r.runOnce.Do(r.safeRun)

	if consumerGroupID != "" {
		r.log.Warnf("consumerGroupID is ignored in grpc reader")
	}

	return r.blockChan
}

func (r *blockchainGrpcReaderImpl) Close() {
	r.closeOnce.Do(r.oneceClose)
}

func (r *blockchainGrpcReaderImpl) oneceClose() {
	close(r.blockChan)
	r.cancel()
}

func (r *blockchainGrpcReaderImpl) safeRun() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				r.lastError = err.(error)
				r.log.Errorf("blockchainGrpcReaderImpl run panic: %v", r.lastError.Error())
			}

			r.Close()
		}()

		if err := r.run(); err != nil {
			r.lastError = err
			r.log.Errorf("blockchainGrpcReaderImpl run failed: %v", err.Error())
			return
		}
	}()
}

func (r *blockchainGrpcReaderImpl) IsSlowMode() bool {
	return r.slowMode
}

func (r *blockchainGrpcReaderImpl) run() error {
	blockchainInfo, err := r.grpc.GetBlockchainInfo()

	if err != nil {
		return err
	}

	r.lastBlockHeight = int64(blockchainInfo.LastBlockHeight)

	startTime := time.Now()

	firstGet := true

	for {
		select {
		case <-r.ctx.Done():
			return nil
		default:
			if r.height > r.lastBlockHeight {
				blockchainInfo, err := r.grpc.GetBlockchainInfo()

				if err != nil {
					return fmt.Errorf("getBlockchainInfo failed: %w", err)
				}

				if int64(blockchainInfo.LastBlockHeight) > r.lastBlockHeight {
					r.lastBlockHeight = int64(blockchainInfo.LastBlockHeight)
					r.slowMode = true
				} else {
					// wait new block
					time.Sleep(time.Second * 5)
					continue
				}
			}

			block, err := r.grpc.GetBlock(uint32(r.height), pactus.BlockVerbosity_BLOCK_VERBOSITY_TRANSACTIONS)

			if err != nil {
				r.log.Errorf("getBlock failed: %v", err.Error())
				return err
			}

			if firstGet {
				firstGet = false
				r.log.Infof("starting get block from grpc at height %d", r.height)
			}

			r.blockChan <- block

			if r.slowMode {
				r.log.Infof("read block %d", block.Height)
			} else {
				if r.height%8640 == 0 || r.height+1 > r.lastBlockHeight {
					timeElapsed := time.Since(startTime)
					r.log.Infof("read block %d/%d (%.2f%%) (%v)", r.height, r.lastBlockHeight, float64(r.height)/float64(r.lastBlockHeight)*100, timeElapsed)
				}
			}

			r.height++
		}
	}
}
