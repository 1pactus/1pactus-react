package gather

import (
	"fmt"
	"time"

	"github.com/1pactus/1pactus-react/lifecycle"
	"github.com/1pactus/1pactus-react/log"
	"github.com/robfig/cron/v3"
)

type DataGatherService struct {
	*lifecycle.ServiceLifeCycle
	log    log.ILogger
	config *Config
	cron   *cron.Cron
}

func NewGatherService(appLifeCycle *lifecycle.AppLifeCycle, config *Config) *DataGatherService {
	return &DataGatherService{
		ServiceLifeCycle: lifecycle.NewServiceLifeCycle(appLifeCycle),
		log:              log.WithKv("service", "gather"),
		config:           config,
		cron:             cron.New(cron.WithLocation(time.UTC)),
	}
}

func (s *DataGatherService) Run() {
	defer s.LifeCycleDead(true)
	defer s.log.Info("BYE")

	gatherChan := make(chan struct{}, 2)

	defer close(gatherChan)

	s.log.Info("HI")

	_, err := s.cron.AddFunc("10 0 * * *", func() {
		s.log.Info("starting scheduled task at UTC 00:10")
		gatherChan <- struct{}{}
	})
	if err != nil {
		s.log.Errorf("failed to add cron job: %v", err)
		return
	}

	s.cron.Start()
	defer s.cron.Stop()

	gatherChan <- struct{}{}

	go s.runGatherWorker(gatherChan)

	<-s.Done()
	s.log.Info("data collect received done signal")
}

func (s *DataGatherService) runGatherWorker(gatherChan <-chan struct{}) {
	s.log.Infof("gather waiting started")
	defer s.log.Infof("gather waiting stopped")
	for range gatherChan {
		err := s.startGatherBlockchain(s.Done())
		if err != nil {
			s.log.Errorf("%v", err)
		}
	}
}

func (s *DataGatherService) startGatherBlockchain(dieChan <-chan struct{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("startGather panic: %v", r)
		}
	}()

	s.log.Infof("blockchain gather started")
	defer s.log.Infof("blockchain gather stopped")

	/*
		cg := NewChainGather(s.log, s.config.GrpcServers)

		if err := cg.Connect(); err != nil {
			return fmt.Errorf("failed to connect to grpc servers: %w", err)
		}

		if err := cg.FetchBlockchain(context.Background()); err != nil {
			return fmt.Errorf("failed to fetch blockchain: %w", err)
		}*/

	cg := NewPgChainGather(s.log, s.config.GrpcServers)

	if err := cg.Connect(); err != nil {
		return fmt.Errorf("failed to connect to grpc servers: %w", err)
	}

	if err := cg.FetchBlockchain(dieChan); err != nil {
		return fmt.Errorf("failed to fetch blockchain: %w", err)
	}

	/*if timeIndexes, err := store.Mongo.GetLastDaysTimeIndex(30); err == nil {
		for _, timeIndex := range timeIndexes {
			s.log.Infof("fetching global state for time index %d", timeIndex)
		}
	}*/

	/*
		if data, err := store.Mongo.GetUnbond(30); err == nil {
			for k, v := range data {
				s.log.Infof("unbond at time index %d: %d", k, v)
			}
		}*/

	return nil
}
