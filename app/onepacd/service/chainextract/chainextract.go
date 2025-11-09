package chainextract

import (
	"sync/atomic"
	"time"

	"github.com/1pactus/1pactus-react/app/onepacd/service/chainextract/chainreader"
	gather "github.com/1pactus/1pactus-react/app/onepacd/service/chainextract/chainreader"
	"github.com/1pactus/1pactus-react/lifecycle"
	"github.com/1pactus/1pactus-react/log"
)

type ChainExtractService struct {
	*lifecycle.ServiceLifeCycle
	log         log.ILogger
	config      *Config
	grpc        *gather.GrpcClient
	kafkaEnable bool
	mainReader  atomic.Value // stores chainreader.BlockchainReader
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
	reader := s.mainReader.Load()
	if reader == nil {
		return nil
	}
	return reader.(chainreader.BlockchainReader)
}

func (s *ChainExtractService) Run() {
	defer s.LifeCycleDead(true)
	defer s.log.Info("Chain Extract Service stopped")
	s.log.Infof("Chain Extract Service is starting...")

	s.log.Infof("start to connect grpc servers: %v", s.grpc.GetServers())

	err := s.grpc.Connect()

	if err != nil {
		s.log.Errorf("failed to connect grpc servers: %v", err.Error())
		return
	}

	s.log.Infof("kafka enable: %v", s.kafkaEnable)

	var reader chainreader.BlockchainReader
	if s.kafkaEnable {
		reader, err = chainreader.NewBlockchainKafkaReader(s.ServiceLifeCycle.Context(), s.grpc, s.log)
	} else {
		reader, err = chainreader.NewBlockchainGrpcReader(s.ServiceLifeCycle.Context(), s.grpc, s.log)
	}

	if err != nil {
		s.log.Errorf("NewBlockchainKafkaReader failed: %v", err.Error())
		return
	}

	s.mainReader.Store(reader)
	defer reader.Close()

	<-s.Done()
}