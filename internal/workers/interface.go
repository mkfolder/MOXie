package workers

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Makefolder/cynero/internal/service"
	"go.uber.org/zap"
)

type Worker interface {
	Start() error
	Stop()
}

type BaseWorker struct {
	s        *service.Service
	log      *zap.SugaredLogger
	interval time.Duration
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	run      func()
}

func (w *BaseWorker) Start() error {
	if w.run == nil {
		return errors.New("failed to start worker: run function is nil")
	}
	w.wg.Add(1)
	go w.run()
	return nil
}

func (w *BaseWorker) Stop() {
	w.cancel()
	w.wg.Wait()
}
