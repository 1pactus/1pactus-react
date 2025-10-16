package dummy

import (
	"github.com/frimin/1pactus-react/lifecycle"
	"github.com/frimin/1pactus-react/log"
)

type DummyService struct {
	*lifecycle.ServiceLifeCycle
	log log.ILogger
}

func NewDummyService(appLifeCycle *lifecycle.AppLifeCycle) *DummyService {
	return &DummyService{
		ServiceLifeCycle: lifecycle.NewServiceLifeCycle(appLifeCycle),
		log:              log.WithKv("service", "dummy"),
	}
}

func (s *DummyService) Run() {
	defer s.LifeCycleDead(true)
	defer s.log.Info("BYE")

	s.log.Info("HI")

	select {
	case <-s.Done():
		s.log.Info("data collect received done signal")
		return
	}
}
