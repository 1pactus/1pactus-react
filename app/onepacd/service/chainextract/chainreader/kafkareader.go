package chainreader

import (
	"context"
	"fmt"
	"sync"

	"github.com/1pactus/1pactus-react/app/onepacd/service/gather"
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

	ctx         context.Context
	cancel      context.CancelFunc
	initOnce    sync.Once
	beginHeight int64
	blockChan   chan *pactus.GetBlockResponse
	closeOnce   sync.Once
}

func (b *blockchainKafkaReaderConsumer) Close() {
	b.closeOnce.Do(func() {
		close(b.blockChan)
		b.cancel()
		b.reader.consumerSyncMap.Delete(b.groupID)
	})
}

func NewBlockchainKafkaReader(parentCtx context.Context, grpc *gather.GrpcClient, parentLogger log.ILogger) (BlockchainReader, error) {
	reader := &blockchainKafkaReaderImpl{
		log: parentLogger.WithKv("reader", "kafka"),
	}

	reader.ctx, reader.cancel = context.WithCancel(parentCtx)

	if err := reader.initGrpcReader(grpc); err != nil {
		return nil, err
	}

	return reader, nil
}

func (r *blockchainKafkaReaderImpl) initGrpcReader(grpc *gather.GrpcClient) error {
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

func (r *blockchainKafkaReaderImpl) Read(beginHeight int64, consumerGroupID string) <-chan *pactus.GetBlockResponse {
	if consumerGroupID != "" {
		r.log.Warnf("consumerGroupID is empty in kafka reader")
		consumerGroupID = "default"
	}

	consumer, _ := r.consumerSyncMap.LoadOrStore(consumerGroupID, &blockchainKafkaReaderConsumer{
		reader:      r,
		beginHeight: beginHeight,
		groupID:     consumerGroupID,
		blockChan:   make(chan *pactus.GetBlockResponse, 100),
	})

	var c <-chan *pactus.GetBlockResponse

	if cons, ok := consumer.(*blockchainKafkaReaderConsumer); ok {
		cons.ctx, cons.cancel = context.WithCancel(r.ctx)
		cons.initOnce.Do(func() {
			r.safeRunConsumer(cons)
		})
		c = cons.blockChan
	}

	r.producerRunOnce.Do(r.safeRunProducer)

	return c
}

func (r *blockchainKafkaReaderImpl) Close() {
	r.closeOnce.Do(r.oneceClose)
}

func (r *blockchainKafkaReaderImpl) oneceClose() {
	r.consumerSyncMap.Range(func(key, value any) bool {
		if consumer, ok := value.(*blockchainKafkaReaderConsumer); ok {
			consumer.Close()
		}
		return true
	})
	r.grpcReader.Close()
	r.cancel()
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

func (r *blockchainKafkaReaderImpl) safeRunConsumer(consumer *blockchainKafkaReaderConsumer) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				r.lastError = err.(error)
				r.log.Errorf("blockchainKafkaReaderImpl run panic: %v", r.lastError.Error())
			}

			r.Close()
		}()

		if err := r.runConsumer(consumer); err != nil {
			r.lastError = err
			r.log.Errorf("blockchainKafkaReaderImpl run failed: %v", err.Error())
			return
		}
	}()
}

func (r *blockchainKafkaReaderImpl) IsSlowMode() bool {
	return r.grpcReader.IsSlowMode()
}

func (r *blockchainKafkaReaderImpl) runProducer() error {
	for block := range r.grpcReader.Read(r.readGrpcStartHeight, "") {
		err := store.Kafka.SendBlock(block)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *blockchainKafkaReaderImpl) runConsumer(consumer *blockchainKafkaReaderConsumer) error {
	topicOffset, err := store.Kafka.GetBlockHeightOffset(consumer.beginHeight)

	if err != nil {
		return fmt.Errorf("GetBlockHeightOffset failed: %w", err)
	}

	err = store.Kafka.ConsumeBlocks(consumer.groupID, topicOffset, func(block *pactus.GetBlockResponse) (bool, error) {
		consumer.blockChan <- block
		return true, nil
	})

	if err != nil {
		return fmt.Errorf("ConsumeBlocks failed: %w", err)
	}

	return nil
}
