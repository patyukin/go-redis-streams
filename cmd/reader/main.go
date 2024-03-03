package main

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/app"
	"github.com/patyukin/go-redis-streams/internal/config"
	"github.com/patyukin/go-redis-streams/internal/processor"
	"github.com/patyukin/go-redis-streams/internal/streamer/redis"
	"github.com/patyukin/go-redis-streams/pkg/logger"
	"go.uber.org/zap"
	"log"
)

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

	reader, err := redis.NewRedisStreamer(wrapCtx, cfg)
	processors := []app.Processor{
		processor.NewBookProcessor(reader),
		processor.NewJournalProcessor(reader),
		processor.NewBookProcessor(reader),
		processor.NewJournalProcessor(reader),
		processor.NewBookProcessor(reader),
		processor.NewJournalProcessor(reader),
	}
	if err != nil {
		l.Fatal("Failed to init redis", zap.String("error", err.Error()))
	}

	a, err := app.NewApp(cfg, processors)
	if err != nil {
		l.Fatal("Failed to init app", zap.String("error", err.Error()))
	}

	err = a.RunStream(wrapCtx)
	if err != nil {
		l.Fatal("Failed to run", zap.String("error", err.Error()))
	}
}
