package application

import (
	"context"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"SnipAndNeat/app/api"
	"SnipAndNeat/app/config"
	"SnipAndNeat/app/service/mailer"
	"SnipAndNeat/app/service/ozon"
	storage "SnipAndNeat/app/strorage"
)

var (
	ErrStartTimeout    = errors.New("start timeout")
	ErrShutdownTimeout = errors.New("shutdown timeout")
)

type Lifecycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type application struct {
	log     zerolog.Logger
	cfg     *config.Config
	syncCmp []Cmp
	cmp     map[string]Lifecycle
}

type Cmp struct {
	Name    string
	Service Lifecycle
}

func New(cfg *config.Config) Lifecycle {
	return &application{
		log:     zerolog.New(os.Stderr).With().Timestamp().Logger(),
		cfg:     cfg,
		syncCmp: []Cmp{},
		cmp:     map[string]Lifecycle{},
	}
}

func (a *application) Start(ctx context.Context) error {
	db := storage.New(a.cfg)
	a.addConsistentService(ctx, "db", db)

	mailer := mailer.NewMailer(a.cfg)

	ozon, err := ozon.New(ctx, a.cfg, db)
	if err != nil {
		a.log.Fatal().Err(err).Msg("cannot create ozon api")
	}

	httpAdmin, err := api.New(a.cfg, mailer, ozon)
	if err != nil {
		a.log.Fatal().Err(err).Msg("cannot create api listner")
	}
	a.addConcurrentService(ctx, "http_admin", httpAdmin)

	// sch := scheduler.New(a.cfg.Scheduler, itr)
	// a.addConcurrentService(ctx, "sch", sch)

	okCh, errCh := make(chan struct{}), make(chan error)
	go func() {
		if err = a.startConsistently(ctx); err != nil {
			a.log.Err(err).Msg("cannot start consistent services")
			errCh <- err
		}
		if err := a.startConcurrently(ctx); err != nil {
			a.log.Err(err).Msg("cannot start concurrent services")
			errCh <- err
		}
		okCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ErrStartTimeout
	case err := <-errCh:
		return err
	case <-okCh:
		return nil
	}

}

// Stop останавливает приложение и все его компоненты
func (a *application) Stop(ctx context.Context) error {
	a.log.Info().Msg("shutting down service...")

	errCh := make(chan error)
	go func() {
		gr, ctx := errgroup.WithContext(ctx)
		for name, service := range a.cmp {
			name, service := name, service
			a.log.Info().Msgf("stopping %q...", name)
			gr.Go(func() error {
				if err := service.Stop(ctx); err != nil {
					a.log.Error().Err(err).Msgf("cannot stop %q", name)
					return err
				}
				return nil
			})
		}
		errCh <- gr.Wait()
	}()

	select {
	case <-ctx.Done():
		return ErrShutdownTimeout
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	}
}
