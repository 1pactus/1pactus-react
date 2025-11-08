package lifecycle

import (
	"context"
)

type ServiceLifeCycle struct {
	appLifeCycle *AppLifeCycle
	ctx          context.Context
	dead         bool
}

func NewServiceLifeCycle(appLifeCycle *AppLifeCycle) *ServiceLifeCycle {
	ctx, _ := context.WithCancel(appLifeCycle.ctx)

	return &ServiceLifeCycle{
		appLifeCycle: appLifeCycle,
		ctx:          ctx,
		dead:         false,
	}
}

func (slc *ServiceLifeCycle) Done() <-chan struct{} {
	return slc.ctx.Done()
}

func (slc *ServiceLifeCycle) Context() context.Context {
	return slc.ctx
}

func (slc *ServiceLifeCycle) LifeCycleDead(stopApp bool) {
	if slc.dead {
		return
	}
	slc.dead = true
	slc.appLifeCycle.wg.Done()
	if stopApp {
		slc.appLifeCycle.StopAppSignal()
	}
}
