package app

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/closer"
	"github.com/patyukin/go-redis-streams/internal/config"
	"log/slog"
	"sync"
)

type App struct {
	serviceProvider *serviceProvider
}

func NewApp(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx, cfg, logger)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context, cfg *config.Config, logger *slog.Logger) error {
	init := []func(ctx context.Context, cfg *config.Config, logger *slog.Logger) error{
		a.initServiceProvider,
	}

	for _, f := range init {
		err := f(ctx, cfg, logger)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.serviceProvider.Streamer(ctx).LimitConsume(ctx, "journals", a.processMessageB)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.serviceProvider.Streamer(ctx).LimitConsume(ctx, "books", a.processMessage)
	}()

	wg.Wait()

	return nil
}

func (a *App) initServiceProvider(_ context.Context, cfg *config.Config, logger *slog.Logger) error {
	a.serviceProvider = newServiceProvider(cfg, logger)
	return nil
}
