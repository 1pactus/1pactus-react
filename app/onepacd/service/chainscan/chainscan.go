package chainscan

import (
	"fmt"
	"time"

	"github.com/1pactus/1pactus-react/app/onepacd/service/chainextract/chainreader"
	"github.com/1pactus/1pactus-react/lifecycle"
	"github.com/1pactus/1pactus-react/log"
	"github.com/robfig/cron/v3"
)

type ReaderProvider interface {
	GetReader() chainreader.BlockchainReader
}

type ChainscanService struct {
	*lifecycle.ServiceLifeCycle
	log            log.ILogger
	config         *Config
	cron           *cron.Cron
	reader         chainreader.BlockchainReader
	readerProvider ReaderProvider
}

func NewChainscanService(appLifeCycle *lifecycle.AppLifeCycle, config *Config, readerProvider ReaderProvider) *ChainscanService {
	return &ChainscanService{
		ServiceLifeCycle: lifecycle.NewServiceLifeCycle(appLifeCycle),
		log:              log.WithKv("service", "chainscan"),
		config:           config,
		cron:             cron.New(cron.WithLocation(time.UTC)),
		readerProvider:   readerProvider,
	}
}

func (s *ChainscanService) Run() {
	defer s.LifeCycleDead(true)
	defer s.log.Info("Chain Scan Service stopped")
	s.log.Infof("Chain Scan Service is starting...")

	gatherChan := make(chan struct{}, 2)

	defer close(gatherChan)

	time.Sleep(1 * time.Second) // wait for chainextract to start

	for {
		s.reader = s.readerProvider.GetReader()
		if s.reader == nil {
			s.log.Warnf("blockchain reader is not initialized yet, retrying in 1 seconds...")
			time.Sleep(1 * time.Second)

			select {
			case <-s.Done(): // exit if service is stopping
				return
			default:
				continue
			}
		} else {
			break
		}
	}

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

	go s.runWorker(gatherChan)

	<-s.Done()
	s.log.Info("data collect received done signal")
}

func (s *ChainscanService) runWorker(gatherChan <-chan struct{}) {
	s.log.Infof("gather waiting started")
	defer s.log.Infof("gather waiting stopped")
	for range gatherChan {
		err := s.startScan(s.Done())
		if err != nil {
			s.log.Errorf("%v", err)
		}
	}
}

func (s *ChainscanService) startScan(dieChan <-chan struct{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("startScan panic: %v", r)
		}
	}()

	timeStart := time.Now()

	s.log.Infof("blockchain scan started, timeStart=%v", timeStart.UTC())
	defer func() {
		s.log.Infof("blockchain scan stopped, timeStart=%v, timeEnd=%v, duration=%v",
			timeStart.UTC(), time.Now().UTC(), time.Since(timeStart))
	}()

	cg := newScanWorker(s.log, s.reader)

	if err := cg.FetchBlockchain(dieChan); err != nil {
		return fmt.Errorf("failed to fetch blockchain: %w", err)
	}

	return nil
}
