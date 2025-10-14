package webapi

import (
	"github.com/frimin/1pactus-react/backend/lifecycle"
	"github.com/frimin/1pactus-react/backend/log"
)

type WebApiService struct {
	*lifecycle.ServiceLifeCycle
	log log.ILogger
}

func NewWebApiService(appLifeCycle *lifecycle.AppLifeCycle) *WebApiService {
	return &WebApiService{
		ServiceLifeCycle: lifecycle.NewServiceLifeCycle(appLifeCycle),
		log:              log.WithKv("service", "webapi"),
	}
}

func (s *WebApiService) Run() {
	defer s.LifeCycleDead(true)
	defer s.log.Info("BYE")
	s.log.Info("HI")

	select {
	case <-s.Done():
		s.log.Info("web api received done signal")
		return
	}
}
