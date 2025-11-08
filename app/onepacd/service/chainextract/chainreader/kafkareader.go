package chainreader

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/1pactus/1pactus-react/app/onepacd/store"
	"github.com/1pactus/1pactus-react/log"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type blockchainKafkaReaderImpl struct {
	grpcReader          BlockchainReader
	log                 log.ILogger
	ctx                 context.Context
	cancel              context.CancelFunc
	lastError           error
	producerRunOnce     sync.Once
	consumerSyncMap     sync.Map
	closeOnce           sync.Once
	readGrpcStartHeight int64
}

type blockchainKafkaReaderConsumer struct {
	reader  *blockchainKafkaReaderImpl
	groupID string
	log     log.ILogger

	lastError   error
	ctx         context.Context
	cancel      context.CancelFunc
	beginHeight int64
	blockChan   chan *pactus.GetBlockResponse
	runOnce     sync.Once
	closeOnce   sync.Once
}

func NewBlockchainKafkaReader(parentCtx context.Context, grpc *GrpcClient, parentLogger log.ILogger) (BlockchainReader, error) {
	reader := &blockchainKafkaReaderImpl{
		log: parentLogger.WithKv("reader", "kafka"),
	}

	reader.ctx, reader.cancel = context.WithCancel(parentCtx)

	if err := reader.initGrpcReader(grpc); err != nil {
		return nil, err
	}

	return reader, nil
}

func (r *blockchainKafkaReaderImpl) initGrpcReader(grpc *GrpcClient) error {
	height, err := store.Kafka.GetLastBlockHeight()
	if err != nil {
		if err == store.ErrorKafkaTopicEmpty {
			r.log.Infof("kafka topic is empty, starting from block height %d", height)
		} else {
			return fmt.Errorf("GetLastBlockHeight failed: %w", err)
		}
	} else {
		r.log.Infof("last block height in kafka is %d", height)
	}

	height++

	r.readGrpcStartHeight = height

	grpcReader, err := NewBlockchainGrpcReader(r.ctx, grpc, r.log.WithField("reader_to", "kafka"))

	if err != nil {
		return fmt.Errorf("newBlockchainGrpcReader failed: %w", err)
	}

	r.grpcReader = grpcReader

	return nil
}

func (r *blockchainKafkaReaderImpl) GetBlockchainInfo() (*pactus.GetBlockchainInfoResponse, error) {
	return r.grpcReader.GetBlockchainInfo()
}

func (r *blockchainKafkaReaderImpl) CreateGroup(beginHeight int64, consumerGroupID string) (BlockchainReaderGroup, bool) {
	consumer, exists := r.consumerSyncMap.LoadOrStore(consumerGroupID, &blockchainKafkaReaderConsumer{
		reader:      r,
		beginHeight: beginHeight,
		groupID:     consumerGroupID,
		blockChan:   make(chan *pactus.GetBlockResponse, 100),
		log:         r.log.WithKv("groupid", consumerGroupID),
	})

	if c, ok := consumer.(*blockchainKafkaReaderConsumer); ok {
		if !exists && ok {
			c.ctx, c.cancel = context.WithCancel(r.ctx)
		}
		return c, exists && ok
	} else {
		return nil, false
	}
}

func (r *blockchainKafkaReaderImpl) Close() {
	r.closeOnce.Do(func() {
		r.consumerSyncMap.Range(func(key, value any) bool {
			if consumer, ok := value.(*blockchainKafkaReaderConsumer); ok {
				consumer.Close()
			}
			return true
		})
		r.grpcReader.Close()
		r.cancel()
	})
}

func (r *blockchainKafkaReaderImpl) safeRunProducer() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				r.lastError = err.(error)
				r.log.Errorf("blockchainKafkaReaderImpl run panic: %v", r.lastError.Error())
			}

			r.Close()
		}()

		if err := r.runProducer(); err != nil {
			r.lastError = err
			r.log.Errorf("blockchainKafkaReaderImpl run failed: %v", err.Error())
			return
		}
	}()
}

func (r *blockchainKafkaReaderImpl) runProducer() error {
	group, _ := r.grpcReader.CreateGroup(r.readGrpcStartHeight, "kafka_producer")

	defer group.Close()

	for block := range group.Read() {
		err := store.Kafka.SendBlock(block)
		if err != nil {
			return err
		}
	}

	return nil
}

//////////// reader group impl

func (g *blockchainKafkaReaderConsumer) Read() <-chan *pactus.GetBlockResponse {
	g.runOnce.Do(g.safeRunConsumer)
	g.reader.producerRunOnce.Do(g.reader.safeRunProducer)

	return g.blockChan
}

func (g *blockchainKafkaReaderConsumer) Close() {
	g.closeOnce.Do(func() {
		g.cancel()
		g.reader.consumerSyncMap.Delete(g.groupID)
	})
}

func (g *blockchainKafkaReaderConsumer) IsSlowMode() bool {
	return false
}

func (r *blockchainKafkaReaderConsumer) safeRunConsumer() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				r.lastError = err.(error)
				r.log.Errorf("blockchainKafkaReaderImpl run panic: %v", r.lastError.Error())
			}

			r.Close()
		}()

		if err := r.runConsumer(); err != nil {
			r.lastError = err
			r.log.Errorf("blockchainKafkaReaderImpl run failed: %v", err.Error())
			return
		}
	}()
}

func (r *blockchainKafkaReaderConsumer) runConsumer() error {
	topicOffset, err := store.Kafka.GetBlockHeightOffset(r.beginHeight)

	if err != nil {
		for {
			r.log.Warnf("GetBlockHeightOffset beginHeight=%v failed, kafka data maybe unavailable, read from grpc now, retrying in 1 minute: %v", r.beginHeight, err)

			for i := 0; i < 60; i++ {
				select {
				case <-r.ctx.Done():
				default:
					time.Sleep(1 * time.Second)
				}
			}

			topicOffset, err = store.Kafka.GetBlockHeightOffset(r.beginHeight)

			if err == nil {
				r.log.Infof("GetBlockHeightOffset beginHeight=%v succeeded", r.beginHeight)
				break
			}
		}
	}

	defer close(r.blockChan)

	err = store.Kafka.ConsumeBlocks(r.ctx, r.groupID, topicOffset, r.blockChan)

	if err != nil {
		if err == context.Canceled {
			r.log.Infof("read canceled")
			return nil
		}
		return fmt.Errorf("ConsumeBlocks failed: %w", err)
	}

	return nil
}
