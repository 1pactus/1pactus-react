package onepacd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/1pactus/1pactus-react/app/onepacd/service/chainextract"
	"github.com/1pactus/1pactus-react/app/onepacd/service/chainscan"
	"github.com/1pactus/1pactus-react/app/onepacd/service/webapi"
	"github.com/1pactus/1pactus-react/app/onepacd/store"
	"github.com/1pactus/1pactus-react/lifecycle"
	"github.com/1pactus/1pactus-react/log"
)

const (
	App     = "onepacd"
	Version = "1.0.0.0"
)

var (
	chainExtractService *chainextract.ChainExtractService
	chainscanService    *chainscan.DataGatherService
	webApiService       *webapi.WebApiService
)

func InitServices(appLifeCycle *lifecycle.AppLifeCycle) error {
	chainExtractService = chainextract.NewChainExtractService(appLifeCycle, conf.Service.ChainExtract, conf.Kafka.Enable)
	chainscanService = chainscan.NewGatherService(appLifeCycle, conf.Service.Chainscan, chainExtractService)
	webApiService = webapi.NewWebApiService(appLifeCycle, conf.App.RunMode, conf.Service.WebApi)

	appLifeCycle.WatchServiceLifeCycle(chainExtractService.ServiceLifeCycle)
	appLifeCycle.WatchServiceLifeCycle(chainscanService.ServiceLifeCycle)
	appLifeCycle.WatchServiceLifeCycle(webApiService.ServiceLifeCycle)

	return nil
}

func RunServices() {
	go chainExtractService.Run()
	go chainscanService.Run()
	go webApiService.Run()
}

func Run() {
	log.Info("Starting application")

	defer log.Info("Stoped application")

	if err := store.Init(conf.ConfigBase); err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}

	appLifeCycle := lifecycle.NewAppLifeCycle()

	if err := InitServices(appLifeCycle); err != nil {
		log.Fatalf("failed to initialize services: %v", err)
	}

	RunServices()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	go func() {
		shutdownInitiated := false
		timeoutTimer := time.NewTimer(0)
		timeoutTimer.Stop()

		for {
			select {
			case sig := <-sigChan:
				if shutdownInitiated {
					log.Infof("received signal %v during shutdown, ignoring", sig)
					continue
				}
				shutdownInitiated = true
				log.Infof("received signal: %v, initiating graceful shutdown", sig)
				appLifeCycle.StopAppSignal()
				timeoutTimer.Reset(30 * time.Second)
			case <-appLifeCycle.ServiceDone():
				log.Info("all services shutdown completed")
				cancel()
				return
			case <-timeoutTimer.C:
				log.Warn("shutdown timeout (30s), forcing exit")
				cancel()
				return
			}
		}
	}()

	<-ctx.Done()
	log.Info("application shutdown completed")
}
