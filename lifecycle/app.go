package lifecycle

import (
	"context"
	"sync"
)

type AppLifeCycle struct {
	//dieChan chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewAppLifeCycle() *AppLifeCycle {
	app := &AppLifeCycle{
		//dieChan: make(chan struct{}),
	}

	app.ctx, app.cancel = context.WithCancel(context.Background())

	return app
}

func (alc *AppLifeCycle) StopAppSignal() {
	alc.cancel()
}

func (alc *AppLifeCycle) WatchServiceLifeCycle(serviceLifeCycle *ServiceLifeCycle) {
	alc.wg.Add(1)
}

func (alc *AppLifeCycle) ServiceDone() <-chan struct{} {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		alc.wg.Wait()
	}()
	return ctx.Done()
}
