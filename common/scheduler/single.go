package scheduler

import (
	"context"
)

func (s *Scheduler) single(job job, ctx context.Context, ctrl contol) {
	go func() {
		if ctrl.isSingleStarted {
			<-ctrl.stopedCh
		}
	}()
	err := (*job.task)(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("single task error")
	}
	go func() {
		ctrl.stopedCh <- struct{}{}
	}()
}
