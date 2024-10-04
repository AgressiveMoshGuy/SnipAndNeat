package scheded

import (
	"SnipAndNeat/app/config"
	"SnipAndNeat/app/service/mailer"
	"SnipAndNeat/common/scheduler"
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Scheduler struct {
	log   *zap.Logger
	cfg   *config.Config
	sch   *scheduler.Scheduler
	Tasks struct {
		ping          scheduler.Task
		sendOzonFiles scheduler.Task
	}

	mailer *mailer.Mailer
}

func New(cfg *config.Config, mailer *mailer.Mailer) *Scheduler {
	sh := &Scheduler{
		log:    zap.L().Named("scheduler"),
		cfg:    cfg,
		sch:    scheduler.New("scheduler", &scheduler.Config{}),
		mailer: mailer,
	}

	// Задание на периодические health-check'и в зависимые сервисы
	sh.Tasks.ping = func(ctx context.Context) error {
		var wg sync.WaitGroup

		for k, v := range map[*bool]func(context.Context) error{} {
			k, v := k, v
			wg.Add(1)
			go func(result *bool, f func(context.Context) error) {
				defer wg.Done()

				err := f(ctx)
				ok := err == nil
				result = &ok
			}(k, v)
		}
		wg.Wait()
		return nil
	}

	return sh
}

func (s *Scheduler) Start(ctx context.Context) error {
	if err := s.sch.AddPeriodic(&s.Tasks.ping, 10*time.Second, true); err != nil {
		s.log.Error("cannot start ping task", zap.Error(err))
		return err
	}

	err := s.sch.Start(ctx)
	if err != nil {
		s.log.Error("cannot start scheduler", zap.Error(err))
		return errors.Wrap(err, "cannot start scheduler")

	}
	return nil
}

func (s *Scheduler) Stop(ctx context.Context) error {
	return nil
}
