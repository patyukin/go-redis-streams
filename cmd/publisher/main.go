package main

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/config"
	"github.com/patyukin/go-redis-streams/internal/streamer"
	"github.com/patyukin/go-redis-streams/pkg/logger"
	"go.uber.org/zap"
	"log"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=StreamWriter --with-expecter=true

type App struct {
	writer StreamWriter
}

func NewApp(writer StreamWriter) *App {
	return &App{
		writer: writer,
	}
}

func (a *App) RunStream(ctx context.Context, n int) error {
	ctx = logger.WithLoggerFields(ctx, zap.String("func", "RunConsume"))

	for i := 0; i < n; i++ {
		var err error
		if i%2 == 0 {
			logger.Debug(ctx, "i%2 == 0")
			err = a.writer.Publish(ctx, "books")
		} else {
			logger.Debug(ctx, "i%2 != 0")
			err = a.writer.Publish(ctx, "journals")
		}

		if err != nil {
			//l.Fatal("Failed to publish", zap.String("error", err.Error()))
		}
	}

	return nil
}

type StreamWriter interface {
	Publish(ctx context.Context, stream string) error
	Close() error
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadYamlConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	l, err := logger.InitLogger(cfg.Logger.Path, "dev")
	if err != nil {
		l.Fatal("Failed to init logger", zap.String("error", err.Error()))
	}

	wrapCtx := context.WithValue(ctx, "logger", l)

	writer, err := streamer.NewRedisStreamer(wrapCtx, cfg)
	if err != nil {
		l.Fatal("Failed to init redis", zap.String("error", err.Error()))
	}

	a := NewApp(writer)
	err = a.RunStream(wrapCtx, 100000000000)
	if err != nil {
		l.Fatal("Failed to init app", zap.String("error", err.Error()))
	}

	// TODO: gf shutdown
}
