package application

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func (a *application) addConsistentService(ctx context.Context, name string, service Lifecycle) {
	a.syncCmp = append(a.syncCmp, Cmp{
		Name:    name,
		Service: service,
	})
}

func (a *application) addConcurrentService(ctx context.Context, name string, service Lifecycle) {
	a.cmp[name] = service
}

func (a *application) startConsistently(ctx context.Context) error {
	for _, cmp := range a.syncCmp {
		a.log.Info().Msgf("starting %q...", cmp.Name)
		if err := cmp.Service.Start(ctx); err != nil {
			a.log.Error().Err(err).Msgf("cannot start %q", cmp.Name)
			return errors.Wrapf(err, "cannot start %q", cmp.Name)
		}
	}
	return nil
}

func (a *application) startConcurrently(ctx context.Context) error {
	gr, ctx := errgroup.WithContext(ctx)
	for name, service := range a.cmp {
		name, service := name, service
		gr.Go(func() error {
			a.log.Info().Msgf("starting %q...", name)
			if err := service.Start(ctx); err != nil {
				a.log.Error().Err(err).Msgf("cannot start %q", name)
				return errors.Wrapf(err, "cannot start %q", name)
			}
			return nil
		})
	}
	return gr.Wait()
}

func (a *application) stopConsistently(ctx context.Context) error {
	for i := len(a.syncCmp); i != 0; i-- {
		a.log.Info().Msgf("stopping %q...", a.syncCmp[i-1].Name)
		if err := a.syncCmp[i-1].Service.Stop(ctx); err != nil {
			a.log.Error().Err(err).Msgf("cannot stop %q", a.syncCmp[i-1].Name)
			return errors.Wrapf(err, "cannot stop %q", a.syncCmp[i-1].Name)
		}
	}
	return nil
}

func (a *application) stopConcurrently(ctx context.Context) error {
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
	return gr.Wait()
}
