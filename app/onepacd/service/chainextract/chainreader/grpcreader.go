package chainreader

import (
	"context"
	"sync"
	"time"

	"github.com/1pactus/1pactus-react/log"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type blockchainGrpcReaderImpl struct {
	consumerSyncMap sync.Map
	grpc            *GrpcClient
	log             log.ILogger
	ctx             context.Context
	cancel          context.CancelFunc
	closeOnce       sync.Once
}

type blockchainGrpcReaderGroupImpl struct {
	reader *blockchainGrpcReaderImpl
	log    log.ILogger

	slowMode  bool
	groupID   string
	lastError error

	ctx    context.Context
	cancel context.CancelFunc

	height          int64
	lastBlockHeight int64

	blockChan chan *pactus.GetBlockResponse

	runOnce   sync.Once
	closeOnce sync.Once
}

func NewBlockchainGrpcReader(parentCtx context.Context, grpc *GrpcClient, parentLogger log.ILogger) (BlockchainReader, error) {
	reader := &blockchainGrpcReaderImpl{
		grpc: grpc,
		log:  parentLogger.WithKv("reader", "grpc"),
	}

	reader.ctx, reader.cancel = context.WithCancel(parentCtx)

	return reader, nil
}

func (r *blockchainGrpcReaderImpl) Close() {
	r.closeOnce.Do(func() {
		r.consumerSyncMap.Range(func(key, value any) bool {
			if consumer, ok := value.(*blockchainGrpcReaderGroupImpl); ok {
				consumer.Close()
			}
			return true
		})
		r.cancel()
	})
}

func (r *blockchainGrpcReaderImpl) GetBlockchainInfo() (*pactus.GetBlockchainInfoResponse, error) {
	return r.grpc.GetBlockchainInfo()
}

func (r *blockchainGrpcReaderImpl) CreateGroup(beginHeight int64, consumerGroupID string) (BlockchainReaderGroup, bool) {
	consumer, exists := r.consumerSyncMap.LoadOrStore(consumerGroupID, &blockchainGrpcReaderGroupImpl{
		reader:    r,
		height:    beginHeight,
		groupID:   consumerGroupID,
		blockChan: make(chan *pactus.GetBlockResponse, 100),
		log:       r.log.WithKv("groupid", consumerGroupID),
	})

	if c, ok := consumer.(*blockchainGrpcReaderGroupImpl); ok {
		if !exists && ok {
			c.ctx, c.cancel = context.WithCancel(r.ctx)
		}
		return c, exists && ok
	} else {
		return nil, false
	}
}

//////////// reader group impl

func (g *blockchainGrpcReaderGroupImpl) Read() <-chan *pactus.GetBlockResponse {
	g.runOnce.Do(g.safeRun)

	return g.blockChan
}

func (g *blockchainGrpcReaderGroupImpl) Close() {
	g.closeOnce.Do(func() {
		g.cancel()
		g.reader.consumerSyncMap.Delete(g.groupID)
	})
}

func (g *blockchainGrpcReaderGroupImpl) safeRun() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				g.lastError = err.(error)
				g.log.Errorf("blockchainGrpcReaderImpl run panic: %v", g.lastError.Error())
			}

			g.Close()
		}()

		if err := g.run(); err != nil {
			g.lastError = err
			g.log.Errorf("blockchainGrpcReaderImpl run failed: %v", err.Error())
			return
		}
	}()
}

func (g *blockchainGrpcReaderGroupImpl) IsSlowMode() bool {
	return g.slowMode
}

func (g *blockchainGrpcReaderGroupImpl) run() error {
	blockchainInfo, err := g.reader.grpc.GetBlockchainInfo()

	if err != nil {
		return err
	}

	defer close(g.blockChan)

	g.lastBlockHeight = int64(blockchainInfo.LastBlockHeight)

	startTime := time.Now()

	firstGet := true

	for {
		select {
		case <-g.ctx.Done():
			return nil
		default:
			if g.height > g.lastBlockHeight {
				blockchainInfo, err := g.reader.grpc.GetBlockchainInfo()

				if err != nil {
					g.log.Errorf("getBlockchainInfo failed, try later: %w", err)
					time.Sleep(5 * time.Second)
					continue
				}

				if int64(blockchainInfo.LastBlockHeight) > g.lastBlockHeight {
					g.lastBlockHeight = int64(blockchainInfo.LastBlockHeight)
					g.slowMode = true
				} else {
					// wait new block
					time.Sleep(time.Second * 5)
					continue
				}
			}

			block, err := g.reader.grpc.GetBlock(uint32(g.height), pactus.BlockVerbosity_BLOCK_VERBOSITY_TRANSACTIONS)

			if err != nil {
				g.log.Errorf("getBlock failed, try later: %v", err.Error())
				time.Sleep(5 * time.Second)
				continue
			}

			if firstGet {
				firstGet = false
				g.log.Infof("starting get block from grpc at height %d", g.height)
			}

			g.blockChan <- block

			if g.slowMode {
				g.log.Infof("read block %d", block.Height)
			} else {
				if g.height%8640 == 0 || g.height+1 > g.lastBlockHeight {
					timeElapsed := time.Since(startTime)
					g.log.Infof("read block %d/%d (%.2f%%) (%v)", g.height, g.lastBlockHeight, float64(g.height)/float64(g.lastBlockHeight)*100, timeElapsed)
				}
			}

			g.height++
		}
	}
}
