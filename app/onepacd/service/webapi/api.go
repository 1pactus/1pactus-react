package webapi

import (
	"context"
	"net/http"
	"time"

	"github.com/frimin/1pactus-react/app/onepacd/service/webapi/middleware"
	"github.com/frimin/1pactus-react/lifecycle"
	"github.com/frimin/1pactus-react/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type WebApiService struct {
	*lifecycle.ServiceLifeCycle
	log    log.ILogger
	config *Config
	mode   string
}

func NewWebApiService(appLifeCycle *lifecycle.AppLifeCycle, mode string, config *Config) *WebApiService {
	return &WebApiService{
		ServiceLifeCycle: lifecycle.NewServiceLifeCycle(appLifeCycle),
		log:              log.WithKv("service", "webapi"),
		config:           config,
	}
}

func (s *WebApiService) Run() {
	defer s.LifeCycleDead(true)
	defer s.log.Info("BYE")
	s.log.Info("HI")

	r := gin.New()

	_log := s.log

	r.Use(cors.Default())
	r.Use(middleware.CustomRecovery(s.log))
	r.Use(middleware.Log(_log.GetInternalLogger()))

	gin.DefaultWriter = _log
	gin.DefaultErrorWriter = _log

	gin.SetMode(s.mode)
	gin.DisableConsoleColor()

	if err := s.setupRoute(r); err != nil {
		log.Errorf("failed to setup route: %v", err)
		return
	}

	srv := &http.Server{
		Addr:           s.config.HttpListen,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.log.Infof("http api listening at %s", s.config.HttpListen)

	serverErr := make(chan error, 1)

	go func() {
		defer log.Infof("webapi server at %s stopped", s.config.HttpListen)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("webapi server error: %v", err)
			serverErr <- err
		}
	}()

	select {
	case <-s.Done():
		s.log.Info("shutting down webapi server...")
	case err := <-serverErr:
		s.log.Errorf("webapi server failed to start: %v", err)
		return
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer shutdownCancel()

	s.log.Info("shutting down webapi server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		s.log.Errorf("webapi server forced to shutdown: %v", err)
		if closeErr := srv.Close(); closeErr != nil {
			s.log.Errorf("failed to force close webapi server: %v", closeErr)
		}
		return
	}

	s.log.Info("webapi server gracefully stopped")
}
