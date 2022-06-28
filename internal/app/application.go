package application

import (
	"context"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	ErrStartTimeout    = errors.New("start timeout")
	ErrShutdownTimeout = errors.New("shutdown timeout")
)

type Lifecycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func New(cfg Config) Lifecycle {
	return &application{
		log: zerolog.New(os.Stderr).With().Timestamp().Logger(),
		cfg: cfg,
		cmp: map[string]Lifecycle{},
	}
}

type application struct {
	log zerolog.Logger
	cfg Config
	cmp map[string]Lifecycle
}

func (a *application) Start(ctx context.Context) error {

	var wg sync.WaitGroup
	okCh, errCh := make(chan struct{}), make(chan error)
	for name, service := range a.cmp {
		wg.Add(1)
		go func(name string, service Lifecycle) {
			defer wg.Done()
			a.log.Debug().Msgf("starting %q...", name)
			if err := service.Start(ctx); err != nil {
				a.log.Err(err).Msgf("cannot start %q %v", name, err)
				errCh <- errors.Wrapf(err, "cannot start %q", name)
				return
			}
		}(name, service)
	}
	go func() {
		wg.Wait()
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
