package gather

import (
	"context"
	"fmt"

	"github.com/frimin/1pactus-react/backend/lifecycle"
	"github.com/frimin/1pactus-react/backend/log"
)

type DataGatherService struct {
	*lifecycle.ServiceLifeCycle
	log    log.ILogger
	config *Config
}

func NewGatherService(appLifeCycle *lifecycle.AppLifeCycle, config *Config) *DataGatherService {
	return &DataGatherService{
		ServiceLifeCycle: lifecycle.NewServiceLifeCycle(appLifeCycle),
		log:              log.WithKv("service", "gather"),
		config:           config,
	}
}

func (s *DataGatherService) Run() {
	defer s.LifeCycleDead(true)
	defer s.log.Info("BYE")

	s.log.Info("HI")

	if err := s.StartGather(); err != nil {
		s.log.Errorf("failed to start gather: %v", err)
		return
	}

	select {
	case <-s.Done():
		s.log.Info("data collect received done signal")
		return
	}
}

func (s *DataGatherService) StartGather() error {
	s.log.Infof("start gather")
	defer s.log.Infof("gather stopped")
	cg := NewChainGather(s.config.GrpcServers)

	if err := cg.Connect(); err != nil {
		return fmt.Errorf("failed to connect to grpc servers: %w", err)
	}

	if err := cg.FetchBlockchain(context.Background()); err != nil {
		return fmt.Errorf("failed to fetch blockchain: %w", err)
	}

	return nil
}
