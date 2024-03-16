package app

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/config"
	"github.com/patyukin/go-redis-streams/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Processor interface {
	Run(ctx context.Context) error
}

type App struct {
	cfg        *config.Config
	processors []Processor
}

func NewApp(cfg *config.Config, processors []Processor) (*App, error) {
	a := &App{
		cfg:        cfg,
		processors: processors,
	}

	return a, nil
}

func (a *App) RunConsume(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctx = logger.WithLoggerFields(ctx, zap.String("func", "RunConsume"))
	g, ctx := errgroup.WithContext(ctx)

	for _, processor := range a.processors {
		prcr := processor
		g.Go(func() error {
			logger.Debug(ctx, "Run processor")
			if err := ctx.Err(); err != nil {
				return err
			}

			err := prcr.Run(ctx)
			if err != nil {
				logger.Error(ctx, "Failed to run processor:", zap.String("error", err.Error()))
				cancel()
			}

			return err
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (a *App) Close() error {
	return nil
}
