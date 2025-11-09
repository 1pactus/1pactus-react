package chainextract

import (
	"fmt"
	"time"

	"github.com/1pactus/1pactus-react/app/onepacd/service/chainextract/chainreader"
	gather "github.com/1pactus/1pactus-react/app/onepacd/service/chainextract/chainreader"
	"github.com/1pactus/1pactus-react/app/onepacd/store"
	"github.com/1pactus/1pactus-react/lifecycle"
	"github.com/1pactus/1pactus-react/log"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type ChainExtractService struct {
	*lifecycle.ServiceLifeCycle
	log         log.ILogger
	config      *Config
	grpc        *gather.GrpcClient
	kafkaEnable bool
	mainReader  chainreader.BlockchainReader
}

func NewChainExtractService(appLifeCycle *lifecycle.AppLifeCycle, config *Config, kafkaEnable bool) *ChainExtractService {
	return &ChainExtractService{
		ServiceLifeCycle: lifecycle.NewServiceLifeCycle(appLifeCycle),
		log:              log.WithKv("service", "chainextract"),
		config:           config,
		grpc:             gather.NewGrpcClient(time.Second*5, config.GrpcServers),
		kafkaEnable:      kafkaEnable,
	}
}

func (s *ChainExtractService) GetReader() chainreader.BlockchainReader {
	return s.mainReader
}

func (s *ChainExtractService) Run() {
	defer s.LifeCycleDead(true)
	defer s.log.Info("BYE")

	s.log.Infof("start to connect grpc servers: %v", s.grpc.GetServers())

	err := s.grpc.Connect()

	if err != nil {
		s.log.Errorf("failed to connect grpc servers: %v", err.Error())
		return
	}

	s.log.Infof("kafka enable: %v", s.kafkaEnable)

	if s.kafkaEnable {
		s.mainReader, err = chainreader.NewBlockchainKafkaReader(s.ServiceLifeCycle.Context(), s.grpc, s.log)
	} else {
		s.mainReader, err = chainreader.NewBlockchainGrpcReader(s.ServiceLifeCycle.Context(), s.grpc, s.log)
	}

	if err != nil {
		s.log.Errorf("NewBlockchainKafkaReader failed: %v", err.Error())
		return
	}

	defer s.mainReader.Close()

	/*
		defer reader.Close()

		go func() {
			defer s.log.Info("kafka reader goroutine exited")
			group, _ := reader.CreateGroup(1, "reader0")
			defer group.Close()

			isStart := true
			startTime := time.Now()

			for block := range group.Read() {
				if isStart {
					isStart = false
					s.log.Infof("started consuming blocks from kafka topic at block height %d", block.Height)
				} else {
					if block.Height%1000 == 0 {
						timeElapsed := time.Since(startTime)
						s.log.WithKv("groupid", "reader0").Infof("block height %d consumed, %v", block.Height, timeElapsed)
					}
				}
			}
		}()

		go func() {
			defer s.log.Info("kafka reader goroutine exited")
			group, _ := reader.CreateGroup(1000000, "reader1")
			defer group.Close()

			isStart := true
			startTime := time.Now()

			for block := range group.Read() {
				if isStart {
					isStart = false
					s.log.Infof("started consuming blocks from kafka topic at block height %d", block.Height)
				} else {
					if block.Height%1000 == 0 {
						timeElapsed := time.Since(startTime)
						s.log.WithKv("groupid", "reader1").Infof("block height %d consumed, %v", block.Height, timeElapsed)
					}
				}
			}
		}()*/

	/*
		if s.kafkaEnable {
			if err := s.fetchBlockchainToKafka(); err != nil {
				s.log.Errorf("fetchBlockchainToKafka failed: %v", err.Error())
				return
			}
		}*/

	/*
		if err := s.fetchBlockchain(); err != nil {
			s.log.Errorf("fetchBlockchain failed: %v", err.Error())
			return
		}*/

	/*if err := s.fetchBlockFromKafka(); err != nil {
		s.log.Errorf("fetchBlockFromKafka failed: %v", err.Error())
		return
	}*/

	<-s.Done()
}

func (s *ChainExtractService) fetchBlockchainToKafka() error {
	defer s.log.Infof("fetchBlockchainToKafka stopped")

	height, err := store.Kafka.GetLastBlockHeight()
	if err != nil {
		if err == store.ErrorKafkaTopicEmpty {
			s.log.Infof("kafka topic is empty, starting from block height %d", height)
		} else {
			return fmt.Errorf("GetLastBlockHeight failed: %w", err)
		}
	} else {
		s.log.Infof("last block height in kafka is %d", height)
	}

	height++

	var lastBlockHeight int64

	blockchainInfo, err := s.grpc.GetBlockchainInfo()

	if err != nil {
		return fmt.Errorf("getBlockchainInfo failed: %w", err)
	}

	lastBlockHeight = int64(blockchainInfo.LastBlockHeight)

	startTime := time.Now()

	firstGet := true
	slowMode := true

	for {
		select {
		case <-s.Done():
			return nil
		default:
			if height > lastBlockHeight {
				blockchainInfo, err := s.grpc.GetBlockchainInfo()

				if err != nil {
					return fmt.Errorf("getBlockchainInfo failed: %w", err)
				}

				if int64(blockchainInfo.LastBlockHeight) > lastBlockHeight {
					lastBlockHeight = int64(blockchainInfo.LastBlockHeight)
					slowMode = true
				} else {
					// wait new block
					time.Sleep(time.Second * 5)
					continue
				}
			}

			block, err := s.grpc.GetBlock(uint32(height), pactus.BlockVerbosity_BLOCK_VERBOSITY_TRANSACTIONS)

			if err != nil {
				s.log.Errorf("getBlock failed: %v", err.Error())
				return err
			}

			if firstGet {
				firstGet = false
				s.log.Infof("starting get block from grpc at height %d", height)
			}

			err = store.Kafka.SendBlock(block)

			if err != nil {
				s.log.Errorf("send block to kafka failed: %v", err.Error())
				return err
			}

			if slowMode {
				s.log.Infof("last block %d send to kafka", block.Height)
			} else {
				if height%8640 == 0 {
					timeElapsed := time.Since(startTime)
					s.log.Infof("block %d/%d (%.2f%%) send to kafka, (%v)", height, lastBlockHeight, float64(height)/float64(lastBlockHeight)*100, timeElapsed)
				}
			}

			height++
		}
	}
}

/*
func (s *ChainExtractService) fetchBlockFromKafka() error {


	height := int64(100_0000)

	topicOffset, err := store.Kafka.GetBlockHeightOffset(height)

	if err != nil {
		return fmt.Errorf("GetBlockHeightOffset failed: %w", err)
	}

	s.log.Infof("starting from block height %d", height)

	startTime := time.Now()

	isStart := true

	// Implement fetching block from Kafka
	err = store.Kafka.ConsumeBlocks("blocks_read", topicOffset, func(block *pactus.GetBlockResponse) (bool, error) {
		if isStart {
			isStart = false
			s.log.Infof("started consuming blocks from kafka topic at block height %d", block.Height)
		} else {
			if block.Height%1000 == 0 {
				timeElapsed := time.Since(startTime)
				s.log.Infof("block height %d consumed, %v", block.Height, timeElapsed)
			}
		}

		return true, nil
	})

	if err != nil {
		return fmt.Errorf("consumeMessages failed: %w", err)
	}

	return nil
}
*/
