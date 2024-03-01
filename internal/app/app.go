package app

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/config"
	"golang.org/x/sync/errgroup"
	"log"
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

func (a *App) RunStream(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, processor := range a.processors {
		processor := processor
		g.Go(func() error {
			err := processor.Run(ctx)
			if err != nil {
				log.Println("Failed to run processor:", err)
				// TODO: прервать все остальные стримы (контекст)
			}

			return err
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
