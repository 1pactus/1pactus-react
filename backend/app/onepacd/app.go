package onepacd

import (
	gather "github.com/frimin/1pactus-react/backend/app/onepacd/service/gather"
	"github.com/frimin/1pactus-react/backend/app/onepacd/service/webapi"
	"github.com/frimin/1pactus-react/backend/app/onepacd/store"
	"github.com/frimin/1pactus-react/backend/lifecycle"
	"github.com/frimin/1pactus-react/backend/log"
)

const (
	App     = "onepacd"
	Version = "1.0.0.0"
)

var (
	dataCollect *gather.DataGatherService
	webApi      *webapi.WebApiService
)

func InitServices(appLifeCycle *lifecycle.AppLifeCycle) error {
	dataCollect = gather.NewGatherService(appLifeCycle, launchConfig.Service.Gather)
	webApi = webapi.NewWebApiService(appLifeCycle)

	appLifeCycle.WatchServiceLifeCycle(dataCollect.ServiceLifeCycle)
	appLifeCycle.WatchServiceLifeCycle(webApi.ServiceLifeCycle)

	return nil
}

func RunServices() {
	go dataCollect.Run()
	go webApi.Run()
}

func Run() {
	log.Info("HI")

	defer log.Info("BYE")

	if err := store.Init(launchConfig.ConfigBase); err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}

	appLifeCycle := lifecycle.NewAppLifeCycle()

	if err := InitServices(appLifeCycle); err != nil {
		log.Fatalf("failed to initialize services: %v", err)
	}

	RunServices()

	appLifeCycle.WaitForAllServiceDone()
}
