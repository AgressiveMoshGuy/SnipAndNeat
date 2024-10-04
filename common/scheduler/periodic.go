package scheduler

import (
	"context"
)

func (s *Scheduler) periodic(job job, ctx context.Context, ctrl contol) {
	go func() {
		if ctrl.isPeriodicStarted {
			<-ctrl.stopedCh
		}
	}()
	err := (*job.task)(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("periodic task error")
	}
	go func() {
		ctrl.stopedCh <- struct{}{}
	}()
}
