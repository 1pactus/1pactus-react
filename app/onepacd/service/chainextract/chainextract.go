package chainextract

import (
	"fmt"
	"time"

	"github.com/1pactus/1pactus-react/app/onepacd/service/gather"
	"github.com/1pactus/1pactus-react/app/onepacd/store"
	"github.com/1pactus/1pactus-react/lifecycle"
	"github.com/1pactus/1pactus-react/log"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type ChainExtractService struct {
	*lifecycle.ServiceLifeCycle
	log         log.ILogger
	config      *Config
	grpc        *gather.GrpcClient
	kafkaEnable bool
}

func NewChainExtractService(appLifeCycle *lifecycle.AppLifeCycle, config *Config, kafkaEnable bool) *ChainExtractService {
	return &ChainExtractService{
		ServiceLifeCycle: lifecycle.NewServiceLifeCycle(appLifeCycle),
		log:              log.WithKv("service", "chainextract"),
		config:           config,
		grpc:             gather.NewGrpcClient(time.Second*10, config.GrpcServers),
		kafkaEnable:      kafkaEnable,
	}
}

func (s *ChainExtractService) Run() {
	defer s.LifeCycleDead(true)
	defer s.log.Info("BYE")

	s.log.Info("HI")

	/*
		if err := s.fetchBlockchain(); err != nil {
			s.log.Errorf("fetchBlockchain failed: %v", err.Error())
			return
		}*/

	if err := s.fetchBlockFromKafka(); err != nil {
		s.log.Errorf("fetchBlockFromKafka failed: %v", err.Error())
		return
	}

	<-s.Done()
}

func (s *ChainExtractService) fetchBlockchain() error {
	err := s.grpc.Connect()

	if err != nil {
		return fmt.Errorf("failed to connect grpc servers: %w", err)
	}

	var height int64
	var lastBlockHeight int64

	blockchainInfo, err := s.grpc.GetBlockchainInfo()

	if err != nil {
		return fmt.Errorf("getBlockchainInfo failed: %w", err)
	}

	lastBlockHeight = int64(blockchainInfo.LastBlockHeight)
	//lastBlockHeight = 10_0000

	startTime := time.Now()

	s.log.Warnf("start")

	for {
		select {
		case <-s.Done():
			return nil
		default:
			height++

			if height >= lastBlockHeight {
				// top reached, wait for new blocks
				height--
				//time.Sleep(time.Second * 10)
				s.log.Info("done")
				return nil
			}

			block, err := s.grpc.GetBlock(uint32(height), pactus.BlockVerbosity_BLOCK_VERBOSITY_TRANSACTIONS)

			if err != nil {
				s.log.Errorf("getBlock failed: %v", err.Error())
				return err
			}

			err = store.Kafka.SendBlock(block)

			if err != nil {
				s.log.Errorf("send block to kafka failed: %v", err.Error())
				return err
			}

			if height%1000 == 0 {
				timeElapsed := time.Since(startTime)
				s.log.Infof("block height %d processed, %v", height, timeElapsed)
			}
		}
	}
}

func (s *ChainExtractService) fetchBlockFromKafka() error {
	height, err := store.Kafka.GetLastBlockHeight()

	if err != nil {
		return fmt.Errorf("GetLastBlockHeight failed: %w", err)
	}

	s.log.Infof("starting from block height %d", height+1)

	startTime := time.Now()

	// Implement fetching block from Kafka
	err = store.Kafka.ConsumeMessages(store.KafkaTopicBlocks, func(message kafka.Message) error {
		var block pactus.GetBlockResponse
		if err := proto.Unmarshal(message.Value, &block); err != nil {
			s.log.Errorf("failed to unmarshal block: %v", err)
			return err
		}

		if block.Height%1000 == 0 {
			timeElapsed := time.Since(startTime)
			s.log.Infof("block height %d consumed, %v", block.Height, timeElapsed)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("consumeMessages failed: %w", err)
	}

	return nil
}
