package chainextract

import (
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

	<-s.Done()
}
