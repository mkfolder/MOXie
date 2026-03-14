package workers

import (
	"context"
	"time"

	"github.com/Makefolder/cynero/internal/service"
	"go.uber.org/zap"
)

// CleanerWorker is supposed to
// - delete expired orders

type CleanerWorker struct {
	BaseWorker
	expirationTime time.Duration
}

func NewCleanerWorker(
	log *zap.SugaredLogger,
	interval, expirationTime time.Duration,
	s *service.Service,
) Worker {
	ctx, cancel := context.WithCancel(context.Background())
	cw := CleanerWorker{
		expirationTime: expirationTime,
		BaseWorker: BaseWorker{
			s:        s,
			log:      log,
			interval: interval,
			ctx:      ctx,
			cancel:   cancel,
		},
	}
	cw.BaseWorker.run = cw.run
	return &cw
}

func (cw *CleanerWorker) run() {
	cw.log.Info("cleaner worker started")
	defer cw.wg.Done()

	cw.work()

	ticker := time.NewTicker(cw.interval)
	defer ticker.Stop()

	for {
		select {
		case <-cw.ctx.Done():
			cw.log.Info("cleaner worker stopping")
			return
		case <-ticker.C:
			cw.work()
		}
	}
}

func (cw *CleanerWorker) work() {
	before := time.Now().UTC().Add(-cw.expirationTime)
	affected, err := cw.s.DeleteOrdersBefore(cw.ctx, before)
	if err != nil {
		cw.log.Errorw("failed to delete orders", "error", err)
		return
	}
	cw.log.Infof("deleted orders %d", affected)
}
